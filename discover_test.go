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

func ExampleDiscoverBridges() {
	locators, _ := DiscoverBridges(false)
	locator := locators[0] // find the first locator
	deviceType := "my nifty app"

	// remember to push the button on your hue first
	bridge, _ := locator.CreateUser(deviceType)
	fmt.Printf("registered new device => %+v\n", bridge)
}

func ExampleNewBridge() {
	bridge := NewBridge("your-ip-address", "your-username")
	lights, _ := bridge.GetAllLights()

	for _, light := range lights {
		light.On()
	}
}

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
		t.Error("expected 8 lights")
	}

	fmt.Printf("%+v\n", lights)
}
