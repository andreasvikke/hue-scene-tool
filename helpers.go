package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func get(path string, v2 bool) ([]byte, error) {
	url := baseURL
	if v2 {
		url += "/clip/v2/resource/" + path
	} else {
		url += fmt.Sprintf("/api/%s/", username) + path
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("hue-application-key", username)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func post(path string, v2 bool, data interface{}) ([]byte, int, error) {
	url := baseURL
	if v2 {
		url += "/clip/v2/resource/" + path
	} else {
		url += fmt.Sprintf("/api/%s/", username) + path
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("hue-application-key", username)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

func delete(path string, v2 bool) ([]byte, int, error) {
	url := baseURL
	if v2 {
		url += "/clip/v2/resource/" + path
	} else {
		url += fmt.Sprintf("/api/%s/", username) + path
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("hue-application-key", username)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}
