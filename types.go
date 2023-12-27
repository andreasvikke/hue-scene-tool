package main

type ZoneResult struct {
	Errors []string `json:"errors"`
	Data   []Zone   `json:"data"`
}

type LightResult struct {
	Errors []string `json:"errors"`
	Data   []Light  `json:"data"`
}

type SceneResult struct {
	Errors []string `json:"errors"`
	Data   []Scene  `json:"data"`
}

type SceneIDResult struct {
	Errors []string  `json:"errors"`
	Data   []SceneID `json:"data"`
}

type SmartSceneResult struct {
	Errors []string `json:"errors"`
	Data   []Scene  `json:"data"`
}

type ZoneCreateResult struct {
	Errors []string `json:"errors"`
	Data   []Child  `json:"data"`
}

type Zone struct {
	ID       string   `json:"id"`
	Metadata Metadata `json:"metadata"`
	Children []Child  `json:"children"`
}

type ZoneCreate struct {
	Metadata Metadata `json:"metadata"`
	Children []Child  `json:"children"`
}

type Child struct {
	ID   string `json:"rid"`
	Type string `json:"rtype"`
}

type Light struct {
	ID    string `json:"id"`
	ID_v1 string `json:"id_v1"`
}

type LightV1 struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SceneID struct {
	ID       string   `json:"id"`
	Metadata Metadata `json:"metadata"`
}

type Scene struct {
	Actions  []Actions `json:"actions"`
	Metadata Metadata  `json:"metadata"`
	Group    Group     `json:"group"`
	Palette  Palette   `json:"palette"`
}

type SmartScene struct {
	Metadata      Metadata        `json:"metadata"`
	Group         Group           `json:"group"`
	WeekTimeslots []WeekTimeslots `json:"week_timeslots"`
}

type Timeslots struct {
	StartTime StartTime `json:"start_time"`
	Target    Target    `json:"target"`
}

type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

type StartTime struct {
	Kind string `json:"kind"`
	Time Time   `json:"time"`
}

type WeekTimeslots struct {
	Timeslots  []Timeslots `json:"timeslots"`
	Recurrence []string    `json:"recurrence"`
}

type Target struct {
	Rid   string `json:"rid"`
	RType string `json:"rtype"`
}
type On struct {
	On bool `json:"on"`
}
type Action struct {
	On               On                `json:"on"`
	Dimming          Dimming           `json:"dimming"`
	ColorTemperature *ColorTemperature `json:"color_temperature,omitempty"`
}
type Actions struct {
	Target Target `json:"target"`
	Action Action `json:"action"`
}
type Metadata struct {
	Name      string `json:"name"`
	Archetype string `json:"archetype"`
}
type Group struct {
	Rid   string `json:"rid"`
	RType string `json:"rtype"`
}
type ColorTemperature struct {
	Mirek int `json:"mirek,omitempty"`
}
type Dimming struct {
	Brightness float64 `json:"brightness"`
}
type ColorTemperaturePallet struct {
	ColorTemperature ColorTemperature `json:"color_temperature"`
	Dimming          Dimming          `json:"dimming"`
}
type Palette struct {
	Color            []any                    `json:"color"`
	Dimming          []Dimming                `json:"dimming"`
	ColorTemperature []ColorTemperaturePallet `json:"color_temperature"`
	Effects          []any                    `json:"effects"`
}
