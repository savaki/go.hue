package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Bridge is a representation of the Philips Hue bridge device.
type Bridge struct {
	IpAddr   string
	Username string
	debug    bool
}

// NewBridge instantiates a bridge object.  Use this method when you already
// know the ip address and username to use.  Saves the trouble of a lookup.
func NewBridge(ipAddr, username string) *Bridge {
	return &Bridge{IpAddr: ipAddr, Username: username}
}

func (self *Bridge) Debug() *Bridge {
	self.debug = true
	return self
}

func (self *Bridge) toUri(path string) string {
	return fmt.Sprintf("http://%s/api/%s%s", self.IpAddr, self.Username, path)
}

func (self *Bridge) get(path string) (*http.Response, error) {
	uri := self.toUri(path)
	if self.debug {
		log.Printf("GET %s\n", uri)
	}
	return client.Get(uri)
}

func (self *Bridge) post(path string, body io.Reader) (*http.Response, error) {
	uri := self.toUri(path)
	if self.debug {
		log.Printf("POST %s\n", uri)
	}
	return client.Post(uri, "application/json", body)
}

func (self *Bridge) put(path string, body io.Reader) (*http.Response, error) {
	uri := self.toUri(path)
	if self.debug {
		log.Printf("PUT %s\n", uri)
	}
	request, err := http.NewRequest("PUT", uri, body)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}

// GetNewLights - retrieves the list lights we've seen since
// the last scan.  returns the new lights, lastseen, and any error
// that may have occured as per:
// http://developers.meethue.com/1_lightsapi.html#12_get_new_lights
func (self *Bridge) GetNewLights() ([]*Light, string, error) {
	response, err := self.get("/lights/new")
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	results := make(map[string]interface{})
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, "", err
	}

	lastScan := results["lastscan"].(string)

	var lights []*Light
	for id, params := range results {
		if id != "lastscan" {
			value := params.(map[string]interface{})["name"]
			light := &Light{Id: id, Name: value.(string)}
			lights = append(lights, light)
		}
	}

	return lights, lastScan, nil
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

// Search - for new lights as per
// http://developers.meethue.com/1_lightsapi.html#13_search_for_new_lights
func (self *Bridge) Search() ([]Result, error) {
	response, err := self.post("/lights", nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results []Result
	err = json.NewDecoder(response.Body).Decode(&results)
	return results, err
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
	var results map[string]Light
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	// and convert them into lights
	var lights []*Light
	for id, params := range results {
		light := Light{Id: id, Name: params.Name, bridge: self}
		lights = append(lights, &light)
	}

	return lights, nil
}
