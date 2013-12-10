package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AccessPoint interface {
	CreateUser(deviceType string) (*Bridge, error)

	Attach(username string) *Bridge
}

type localBridge struct {
	Id      string `json:"id"`
	IpAddr  string `json:"internalipaddress"`
	MacAddr string `json:"macaddress"`
}

func (self localBridge) CreateUser(deviceType string) (*Bridge, error) {
	// construct our json params
	params := map[string]string{"devicetype": deviceType}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// create a new user
	uri := fmt.Sprintf("http://%s/api", self.IpAddr)
	response, err := http.Post(uri, "text/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// extract username from the results
	var results []map[string]map[string]string
	json.NewDecoder(response.Body).Decode(&results)
	value := results[0]
	username := value["success"]["username"]

	// and create the new bridge object
	return &Bridge{IpAddr: self.IpAddr, Username: username}, nil
}

func (self localBridge) Attach(username string) *Bridge {
	return &Bridge{IpAddr: self.IpAddr, Username: username}
}

func Discover() ([]AccessPoint, error) {
	response, err := http.Get("https://www.meethue.com/api/nupnp")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var bridges []localBridge
	err = json.NewDecoder(response.Body).Decode(&bridges)
	if err != nil {
		return nil, err
	}

	// convert local bridges to access points
	var points []AccessPoint
	for _, bridge := range bridges {
		points = append(points, AccessPoint(bridge))
	}

	return points, err
}
