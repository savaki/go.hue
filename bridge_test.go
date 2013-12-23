package hue

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetNewLights(t *testing.T) {
	bridge := NewBridge(host, username)
	_, _, err := bridge.GetNewLights()
	if err != nil {
		t.Fail()
	}
}

func TestSearchAndGetNewLights(t *testing.T) {
	bridge := NewBridge(host, username)
	_, err := bridge.Search()
	if err != nil {
		fmt.Printf("unable to search => %s\n", err)
		t.Fail()
	}

	lights, lastScan, err := bridge.GetNewLights()
	if err != nil {
		fmt.Printf("GetNewLights => %s\n", err)
		t.Fail()
	}

	fmt.Printf("lights   => %+v\n", lights)
	fmt.Printf("lastScan => %+v\n", lastScan)
}

func TestParseGetNewLights(t *testing.T) {
	data := []byte(`{
		"7": {"name": "Hue Lamp 7"},
		"8": {"name": "Hue Lamp 8"},
		"lastscan": "2012-10-29T12:00:00"
	}`)

	results := make(map[string]interface{})
	err := json.Unmarshal(data, &results)
	if err != nil {
		fmt.Printf("%+v\n", err)
		t.Fail()
	}

	value := results["7"].(map[string]interface{})
	fmt.Printf("%#v\n", value["name"])
}
