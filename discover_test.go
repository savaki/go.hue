package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

var username = os.Getenv("HUE_USERNAME")

func TestDiscoverBridges(t *testing.T) {
	locators, err := DiscoverBridges()
	if err != nil {
		log.Printf("%+v\n", err)
		t.Fail()
	}

	if locators == nil {
		log.Printf("locators was null!  this should never happen\n")
		t.Fail()
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
	bridge := NewBridge("10.0.1.11", username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		log.Printf("%+v\n", err)
		t.Fail()
	}

	if len(lights) != 8 {
		log.Printf("expected 6 lights\n")
		t.Fail()
	}

	fmt.Printf("%+v\n", lights)
}
