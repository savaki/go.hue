package hue

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetLightAttributes(t *testing.T) {
	bridge := NewBridge("10.0.1.11", username)
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
	bridge := NewBridge("10.0.1.11", username)
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

func TestSetLightStateSerializes(t *testing.T) {
	state := SetLightState{}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fail()
	}
	if string(data) != "{}" {
		fmt.Printf("expected {} ; got => %+v\n", string(data))
		t.Fail()
	}
}

func TestSetLightStateSerializesStringToInt(t *testing.T) {
	state := SetLightState{Bri: "123"}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fail()
	}
	if string(data) != `{"bri":123}` {
		fmt.Printf(`expected {"bri":123} ; :123 => %+v\n`, string(data))
		t.Fail()
	}
}

func TestColorLoop(t *testing.T) {
	bridge := NewBridge("10.0.1.11", username)
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
	bridge := NewBridge("10.0.1.11", username)
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
