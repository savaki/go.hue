package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

var host = os.Getenv("HUE_HOST")
var username = os.Getenv("HUE_USERNAME")

func TestDiscoverBridges(t *testing.T) {
	locators, err := DiscoverBridges(false)
	if err != nil {
		log.Printf("%+v\n", err)
		t.Fail()
	}

	if locators == nil {
		t.Error("locators was null! this should never happen")
	}

	log.Printf("%+v\n", locators)
}

func TestDeserializeComplexMap(t *testing.T) {
	data := []byte(`[{"success":{"username": "1234567890"}}]`)

	var results []map[string]map[string]string
	json.NewDecoder(bytes.NewReader(data)).Decode(&results)

	if results == nil {
		t.Fail()
	}
	fmt.Printf("%+v\n", results)
}

func TestGetAllLights(t *testing.T) {
	bridge := NewBridge(host, username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		log.Printf("%+v\n", err)
		t.Fail()
	}

	if len(lights) != 8 {
		t.Error(log.Printf("expected 8 lights"))
	}

	fmt.Printf("%+v\n", lights)
}
