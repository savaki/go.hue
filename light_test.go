package hue

import (
	"fmt"
	"testing"
)

func TestGetLightAttributes(t *testing.T) {
	bridge := NewBridge(host, username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		t.Fail()
	}

	light := lights[0]
	state, err := light.GetLightAttributes()
	if err != nil {
		t.Fail()
	}

	fmt.Printf("%+v\n", state)
}

func TestSetName(t *testing.T) {
	bridge := NewBridge(host, username)
	lights, err := bridge.GetAllLights()
	if err != nil {
		t.Fail()
	}

	light := lights[0]
	state, err := light.SetName("Bedroom Light")
	if err != nil {
		t.Fail()
	}

	fmt.Printf("%+v\n", state)
}

func TestColorLoop(t *testing.T) {
	bridge := NewBridge(host, username)
	light, err := bridge.FindLightById("3")
	if err != nil {
		t.Fail()
	}

	_, err = light.ColorLoop()
	if err != nil {
		fmt.Printf("unable to start color loop => %s\n", err)
		t.Fail()
	}
}

func TestOn(t *testing.T) {
	bridge := NewBridge(host, username)
	light, err := bridge.FindLightById("3")
	if err != nil {
		t.Fail()
	}

	_, err = light.On()
	if err != nil {
		fmt.Printf("unable to start color loop => %s\n", err)
		t.Fail()
	}
}
