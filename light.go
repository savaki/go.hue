package hue

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Light - encapsulates the controls for a specific philips hue light
type Light struct {
	Id     string
	Name   string
	bridge *Bridge
}

type LightState struct {
	Hue       int       `json:"hue"`
	On        bool      `json:"on"`
	Effect    string    `json:"effect"`
	Alert     string    `json:"effect"`
	Bri       int       `json:"bri"`
	Sat       int       `json:"sat"`
	Ct        int       `json:"ct"`
	Xy        []float32 `json:"xy"`
	Reachable bool      `json:"reachable"`
	ColorMode string    `json:"colormode"`
}

type SetLightState struct {
	On             string
	Bri            string
	Hue            string
	Sat            string
	Xy             []float32
	Ct             string
	Alert          string
	Effect         string
	TransitionTime string
}

type LightAttributes struct {
	State           LightState        `json:"state"`
	Type            string            `json:"type"`
	Name            string            `json:"name"`
	ModelId         string            `json:"modelid"`
	SoftwareVersion string            `json:"swversion"`
	PointSymbol     map[string]string `json:"pointsymbol"`
}

// GetLightAttributes - retrieves light attributes and state as per
// http://developers.meethue.com/1_lightsapi.html#14_get_light_attributes_and_state
func (self *Light) GetLightAttributes() (*LightAttributes, error) {
	response, err := self.bridge.get("/lights/" + self.Id)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result := new(LightAttributes)
	err = json.NewDecoder(response.Body).Decode(&result)
	return result, err
}

// SetName - sets the name of a light as per
// http://developers.meethue.com/1_lightsapi.html#15_set_light_attributes_rename
func (self *Light) SetName(newName string) ([]Result, error) {
	params := map[string]string{"name": newName}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	response, err := self.bridge.put("/lights/"+self.Id, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results []Result
	err = json.NewDecoder(response.Body).Decode(&results)
	return results, err
}

// On - a convenience method to turn on a light and set its effect to "none"
func (self *Light) On() ([]Result, error) {
	state := SetLightState{
		On:     "true",
		Effect: "none",
	}
	return self.SetState(state)
}

// Off - a convenience method to turn off a light
func (self *Light) Off() ([]Result, error) {
	state := SetLightState{On: "false"}
	return self.SetState(state)
}

// ColorLoop - a convenience method to turn on a light and have it begin
// a colorloop effect
func (self *Light) ColorLoop() ([]Result, error) {
	state := SetLightState{
		On:     "true",
		Effect: "colorloop",
	}
	return self.SetState(state)
}

// SetState sets the state of a light as per
// http://developers.meethue.com/1_lightsapi.html#16_set_light_state
func (self *Light) SetState(state SetLightState) ([]Result, error) {
	params := make(map[string]interface{})

	if state.On != "" {
		value, _ := strconv.ParseBool(state.On)
		params["on"] = value
	}
	if state.Bri != "" {
		params["bri"], _ = strconv.Atoi(state.Bri)
	}
	if state.Hue != "" {
		params["hue"], _ = strconv.Atoi(state.Hue)
	}
	if state.Sat != "" {
		params["sat"], _ = strconv.Atoi(state.Sat)
	}
	if state.Xy != nil {
		params["xy"] = state.Xy
	}
	if state.Ct != "" {
		params["ct"], _ = strconv.Atoi(state.Ct)
	}
	if state.Alert != "" {
		params["alert"] = state.Alert
	}
	if state.Effect != "" {
		params["effect"] = state.Effect
	}
	if state.TransitionTime != "" {
		params["transitiontime"], _ = strconv.Atoi(state.TransitionTime)
	}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	response, err := self.bridge.put("/lights/"+self.Id+"/state", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results []Result
	err = json.NewDecoder(response.Body).Decode(&results)
	return results, err
}
