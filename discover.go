package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// BridgeLocator represents a factory method you can use to obtain a
// bridge reference.  Depending on whether this is your first time
// connecting to a bridge and you need a new username or you're reconnecting
// to an existing bridge, you call #CreateUser and #Attach respectively
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
	Serial  string `json:"id"`
	IpAddr  string `json:"internalipaddress"`
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
	response, err := client.Post(uri, "text/json", bytes.NewReader(data))
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

// DiscoverBridges is a two-step approach trying to find your hue bridges.
// First it will try to discover bridges in your network using UPnP. Should
// this fail it will utilize the hue api (https://www.meethue.com/api/nupnp) to
// fetch a list of known bridges at the current location. This fallback method
// assumes that you've already set up the hue and can access it via your mobile
// device. If findAllBridges is true the UPnP discovery will wait for all
// bridges to respond. When set to false, this method will return as soon as it
// found the first bridge in your network.
func DiscoverBridges(findAllBridges bool) ([]BridgeLocator, error) {
	var bridges []localBridge

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 1900})
	if err == nil {
		conn.SetDeadline(time.Now().Add(3 * time.Second))
		b := "M-SEARCH * HTTP/1.1\r\n" +
				"HOST: 239.255.255.250:1900\r\n" +
				"MAN: \"ssdp:discover\"\r\n" +
				"MX: 3\r\n" +
				"ST: go.hue:idl\r\n"

		_, err := conn.WriteToUDP([]byte(b), &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1900})
		if err == nil {
			var buf []byte = make([]byte, 8192)
			for {
				_, addr, err := conn.ReadFromUDP(buf)
				if err != nil {
					break
				}

				// Since the hue responds 6 times via UDP, we need to filter out dupes
				dupe := false
				for _, v := range bridges {
					if v.IpAddr == addr.IP.String() {
						dupe = true
					}
				}
				if dupe {
					continue
				}

				// Response sanity check
				if !strings.Contains(string(buf), "LOCATION: ") {
					continue
				}
				var descUrl string
				{
					s := strings.SplitAfter(string(buf), "LOCATION: ")
					descUrl = strings.Split(s[1], "\n")[0]
				}

				// Fetch description.xml from hue
				response, err := http.Get(descUrl)
				if err != nil {
					continue
				}
				defer response.Body.Close()
				body, err := ioutil.ReadAll(response.Body)

				// Make sure we really found a hue
				if !strings.Contains(string(body), "Philips hue") {
					continue
				}

				// Extract serial number, avoid nasty xml unmarshalling
				var serial string
				{
					s := strings.SplitAfter(string(body), "serialNumber>")
					serial = strings.Split(s[1], "</serialNumber>")[0]
				}
				bridges = append(bridges, localBridge{
					Serial: serial,
					IpAddr: addr.IP.String(),
				})
				if !findAllBridges {
					break
				}
			}
		}
	}

	if len(bridges) == 0 {
		// fallback method
		response, err := http.Get("https://www.meethue.com/api/nupnp")
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		err = json.NewDecoder(response.Body).Decode(&bridges)
		if err != nil {
			return nil, err
		}
	}

	// convert local bridges to access points
	var points []BridgeLocator
	for _, bridge := range bridges {
		points = append(points, BridgeLocator(bridge))
	}

	return points, err
}
