package hue

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDeserializeResult(t *testing.T) {
	var data = []byte(`[{"success":{"/lights/1/name":"Bedroom Light"}}]`)
	var results []Result
	err := json.Unmarshal(data, &results)
	if err != nil {
		t.Fail()
	}

	if len(results) != 1 {
		t.Fail()
	}
	if results[0].Success["/lights/1/name"] != "Bedroom Light" {
		t.Fail()
	}

	fmt.Printf("%+v\n", results)
}
