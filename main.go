package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var baseURL = os.Getenv("HUE_URL")
var username = os.Getenv("HUE_USERNAME")

func main() {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	lightsV1, err := getLightsV1()
	if err != nil {
		log.Fatalf("error getting lights: %s", err)
	}

	scenes, err := getSceneTemplates()
	if err != nil {
		log.Fatalf("error getting scenes: %s", err)
	}

	smartScenes, err := getSmartSceneTemplates()
	if err != nil {
		log.Fatalf("error getting smartscenes: %s", err)
	}

	// Create zones for each light
	for lightIDV1, light := range lightsV1 {
		// Check if a zone already exists for the light
		zoneName := fmt.Sprintf("DEV - %s", light.Name)
		zoneExists, err := zoneExists(zoneName)
		if err != nil {
			log.Fatalf("error checking if zone exists: %s", err)
		}
		if zoneExists {
			zoneID, err := getZoneIDByIDName(zoneName)
			if err != nil {
				log.Fatalf("error getting zone ID: %s", err)
			}
			err = deleteZone(zoneID)
			if err != nil {
				log.Fatalf("error deleting zone: %s", err)
			}
		}

		zoneID, err := createZone(light.Name, lightIDV1)
		if err != nil {
			log.Fatalf("error creating zone: %s", err)
		}

		// Create scenes for each entry in scenes.json
		for _, scene := range scenes {
			// Check if a scene already exists for the light and scene combination
			sceneExists, err := sceneExists(fmt.Sprintf("%s - %s", scene.Metadata.Name, light.Name))
			if err != nil {
				log.Fatalf("error checking if scene exists: %s", err)
			}
			if !sceneExists {
				err = createScene(zoneID, lightIDV1, light, scene)
				if err != nil {
					log.Fatalf("error creating scene: %s", err)
				}
			}
		}

		// Create smartscenes for each entry in smart_scenes.json
		for _, smartScene := range smartScenes {
			// Check if a scene already exists for the light and scene combination
			smartSceneExists, err := smartSceneExists(fmt.Sprintf("%s - %s", smartScene.Metadata.Name, light.Name))
			if err != nil {
				log.Fatalf("error checking if scene exists: %s", err)
			}
			if !smartSceneExists {
				err = createSmartScene(zoneID, lightIDV1, light.Name, smartScene)
				if err != nil {
					log.Fatalf("error creating scene: %s", err)
				}
			}
		}
	}
}

func getLightsV1() (map[string]LightV1, error) {
	resp, err := get("lights", false)
	if err != nil {
		return nil, err
	}

	lights := make(map[string]LightV1)
	err = json.Unmarshal(resp, &lights)
	if err != nil {
		return nil, err
	}

	return lights, nil
}

func getSceneTemplates() ([]Scene, error) {
	// Load scenes.json
	file, err := os.ReadFile("scenes.json")
	if err != nil {
		return nil, err
	}

	// Unmarshal scenes.json into a slice of Scene structs
	var scenes []Scene
	err = json.Unmarshal(file, &scenes)
	if err != nil {
		return nil, err
	}

	return scenes, nil
}

func getSmartSceneTemplates() ([]SmartScene, error) {
	// Load scenes.json
	file, err := os.ReadFile("smart_scenes.json")
	if err != nil {
		return nil, err
	}

	// Unmarshal scenes.json into a slice of Scene structs
	var smartScenes []SmartScene
	err = json.Unmarshal(file, &smartScenes)
	if err != nil {
		return nil, err
	}

	return smartScenes, nil
}

func zoneExists(name string) (bool, error) {
	resBody, err := get("zone", true)
	if err != nil {
		return false, err
	}

	zones := ZoneResult{}
	err = json.Unmarshal(resBody, &zones)
	if err != nil {
		return false, err
	}

	for _, zone := range zones.Data {
		if zone.Metadata.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func createZone(lightName, lightIDV1 string) (string, error) {
	fmt.Printf("Creating zone for %s\n", lightName)

	lightID, err := getLightIDByIDV1(lightIDV1)
	if err != nil {
		return "", err
	}

	child := Child{
		ID:   lightID,
		Type: "light",
	}

	zone := ZoneCreate{
		Metadata: Metadata{
			Name:      fmt.Sprintf("DEV - %s", lightName),
			Archetype: "other",
		},
		Children: []Child{child},
	}

	fmt.Println(zone)

	response, statusCode, err := post("zone", true, zone)
	if err != nil {
		return "", err
	}

	zoneResult := ZoneCreateResult{}
	err = json.Unmarshal(response, &zoneResult)
	if err != nil {
		return "", err
	}

	fmt.Println(string(response))
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return "", fmt.Errorf("error creating zone: %s", string(response))
	}

	return zoneResult.Data[0].ID, nil
}

func deleteZone(id string) error {
	fmt.Printf("Deleting zone with id %s\n", id)

	response, statusCode, err := delete(fmt.Sprintf("zone/%s", id), true)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return fmt.Errorf("error deleting zone: %s", string(response))
	}

	return nil
}

func sceneExists(sceneName string) (bool, error) {
	response, err := get("scene", true)
	if err != nil {
		return false, err
	}

	var scenes SceneResult
	if err := json.Unmarshal(response, &scenes); err != nil {
		return false, err
	}

	for _, scene := range scenes.Data {
		if scene.Metadata.Name == sceneName {
			return true, nil
		}
	}

	return false, nil
}

func smartSceneExists(smartSceneName string) (bool, error) {
	response, err := get("smart_scene", true)
	if err != nil {
		return false, err
	}

	var smartScenes SmartSceneResult
	if err := json.Unmarshal(response, &smartScenes); err != nil {
		return false, err
	}

	for _, smartScene := range smartScenes.Data {
		if smartScene.Metadata.Name == smartSceneName {
			return true, nil
		}
	}

	return false, nil
}

func createScene(zoneID, lightIDV1 string, light LightV1, sceneTemplate Scene) error {
	fmt.Printf("Creating scene %s for %s\n", sceneTemplate.Metadata.Name, light.Name)

	realID, err := getLightIDByIDV1(lightIDV1)
	if err != nil {
		return err
	}

	newScene := sceneTemplate
	newScene.Metadata.Name = fmt.Sprintf("%s - %s", sceneTemplate.Metadata.Name, light.Name)
	newScene.Group = Group{
		Rid:   zoneID,
		RType: "zone",
	}

	if light.Type == "Dimmable light" {
		newScene.Actions = append(newScene.Actions, Actions{
			Target: Target{
				Rid:   realID,
				RType: "light",
			},
			Action: Action{
				On:               On{On: true},
				Dimming:          sceneTemplate.Palette.ColorTemperature[0].Dimming,
				ColorTemperature: nil,
			},
		})
	} else {
		newScene.Actions = append(newScene.Actions, Actions{
			Target: Target{
				Rid:   realID,
				RType: "light",
			},
			Action: Action{
				On:               On{On: true},
				Dimming:          sceneTemplate.Palette.ColorTemperature[0].Dimming,
				ColorTemperature: &sceneTemplate.Palette.ColorTemperature[0].ColorTemperature,
			},
		})
	}

	respBody, statusCode, err := post("scene", true, newScene)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		ns, err := json.Marshal(newScene)
		if err != nil {
			return fmt.Errorf("failed to marshall scene")
		}
		return fmt.Errorf("failed to create scene: %d, %s, %s", statusCode, string(respBody), ns)
	}

	return nil
}

func createSmartScene(zoneID, lightIDV1, lightName string, smartSceneTemplate SmartScene) error {
	fmt.Printf("Creating smart scene %s for %s\n", smartSceneTemplate.Metadata.Name, lightName)

	newScene := smartSceneTemplate
	newScene.Metadata.Name = fmt.Sprintf("%s - %s", smartSceneTemplate.Metadata.Name, lightName)
	newScene.Group = Group{
		Rid:   zoneID,
		RType: "zone",
	}
	newScene.WeekTimeslots = []WeekTimeslots{
		{
			Timeslots:  []Timeslots{},
			Recurrence: smartSceneTemplate.WeekTimeslots[0].Recurrence,
		},
	}

	for _, timeslot := range smartSceneTemplate.WeekTimeslots[0].Timeslots {
		sceneName := fmt.Sprintf("%s - %s", timeslot.Target.Rid, lightName)
		sceneID, err := getSceneIDByName(sceneName)
		if err != nil {
			return err
		}

		newScene.WeekTimeslots[0].Timeslots = append(newScene.WeekTimeslots[0].Timeslots, Timeslots{
			StartTime: timeslot.StartTime,
			Target: Target{
				Rid:   sceneID,
				RType: "scene",
			},
		})
	}

	respBody, statusCode, err := post("smart_scene", true, newScene)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return fmt.Errorf("failed to create smartscene: %s", string(respBody))
	}

	return nil
}

func getLightIDByIDV1(lightIDV1 string) (string, error) {
	response, err := get("light", true)
	if err != nil {
		return "", err
	}

	var lights LightResult
	if err := json.Unmarshal(response, &lights); err != nil {
		return "", err
	}

	for _, light := range lights.Data {
		if light.ID_v1 == "/lights/"+lightIDV1 {
			return light.ID, nil
		}
	}

	return "", fmt.Errorf("no lights found for v1 id: %s", lightIDV1)
}

func getZoneIDByIDName(name string) (string, error) {
	resBody, err := get("zone", true)
	if err != nil {
		return "", err
	}

	zones := ZoneResult{}
	err = json.Unmarshal(resBody, &zones)
	if err != nil {
		return "", err
	}

	for _, zone := range zones.Data {
		if zone.Metadata.Name == name {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("no zone found for name: %s", name)
}

func getSceneIDByName(sceneName string) (string, error) {
	response, err := get("scene", true)
	if err != nil {
		return "", err
	}

	var scenes SceneIDResult
	if err := json.Unmarshal(response, &scenes); err != nil {
		return "", err
	}

	for _, scene := range scenes.Data {
		if scene.Metadata.Name == sceneName {
			return scene.ID, nil
		}
	}

	return "", fmt.Errorf("no scene found with name: %s", sceneName)
}
