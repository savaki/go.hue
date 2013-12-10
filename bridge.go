package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Bridge is a representation of the Philips Hue bridge device.
type Bridge struct {
	IpAddr   string
	Username string
}

// NewBridge instantiates a bridge object.  Use this method when you already
// know the ip address and username to use.  Saves the trouble of a lookup.
func NewBridge(ipAddr, username string) *Bridge {
	return &Bridge{IpAddr: ipAddr, Username: username}
}

func (self *Bridge) toUri(path string) string {
	return fmt.Sprintf("http://%s/api/%s%s", self.IpAddr, self.Username, path)
}

func (self *Bridge) get(path string) (*http.Response, error) {
	uri := self.toUri(path)
	log.Printf("GET %s\n", uri)
	return http.Get(uri)
}

func (self *Bridge) put(path string, body io.Reader) (*http.Response, error) {
	uri := self.toUri(path)
	log.Printf("PUT %s\n", uri)
	request, err := http.NewRequest("PUT", uri, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(request)
}

// FindLightById allows you to easily look up light if you know it's Id
func (self *Bridge) FindLightById(id string) (*Light, error) {
	lights, err := self.GetAllLights()
	if err != nil {
		return nil, err
	}

	for _, light := range lights {
		if light.Id == id {
			return light, nil
		}
	}

	return nil, errors.New("unable to find light with id, " + id)
}

// FindLightByName - similar to FindLightById, this is a convenience method
// for when you already know the name of the light
func (self *Bridge) FindLightByName(name string) (*Light, error) {
	lights, err := self.GetAllLights()
	if err != nil {
		return nil, err
	}

	for _, light := range lights {
		if light.Name == name {
			return light, nil
		}
	}

	return nil, errors.New("unable to find light with name, " + name)
}

// GetAllLights - retrieves all lights the Hue is aware of
func (self *Bridge) GetAllLights() ([]*Light, error) {
	// fetch all the lights
	response, err := self.get("/lights")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// deconstruct the json results
	var results map[string]map[string]string
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	// and convert them into lights
	var lights []*Light
	for id, params := range results {
		light := Light{Id: id, Name: params["name"], bridge: self}
		lights = append(lights, &light)
	}

	return lights, nil
}
