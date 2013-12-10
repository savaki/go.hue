package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// BridgeLocator represents a factory method you can use to obtain a
// bridge reference.
type BridgeLocator interface {
	// If this is your first time connnecting to a bridge,
	// you'll want to call CreateUser(deviceType) where deviceType is any
	// arbitrary string you can use to name your device.  This will ask
	// the hub to create a new user with a random username.  The bridge
	// returned back to you will be populated with the username
	CreateUser(deviceType string) (*Bridge, error)

	// Attach should be used when you already know the username and want
	// to obtain a reference to the bridge
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

// DiscoverBridges utilizes the hue api (https://www.meethue.com/api/nupnp) to
// fetch a list of known bridges at the current location.  This method assumes
// that you've already set up the hue and can access it via your mobile device
func DiscoverBridges() ([]BridgeLocator, error) {
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
	var points []BridgeLocator
	for _, bridge := range bridges {
		points = append(points, BridgeLocator(bridge))
	}

	return points, err
}
