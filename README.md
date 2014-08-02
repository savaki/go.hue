go.hue
======

[![GoDoc](http://godoc.org/github.com/savaki/go.hue?status.png)](http://godoc.org/github.com/savaki/go.hue)

An easy to use api to manage your phillips hue.  For documentation, check out the link to godoc above.

### Example - Register a New Device

To start using the hue api, you first need to register your device.

```
package main

import (
	"fmt"
	"github.com/savaki/go.hue"
)

func main() {
	locators, _ := hue.DiscoverBridges(false)
	locator := locators[0] // find the first locator
	deviceType := "my nifty app"

	// remember to push the button on your hue first
	bridge, _ := locator.CreateUser(deviceType)
	fmt.Printf("registered new device => %+v\n", bridge)
}
```

### Example - Turn on all the lights

```
package main

import (
	"github.com/savaki/go.hue"
)

func main() {
	bridge := hue.NewBridge("your-ip-address", "your-username")
	lights, _ := bridge.GetAllLights()

	for _, light := range lights {
		light.On()
	}
}

```

### Example - Disco Time!  Turn all lights on with colorloop

```
package main

import (
	"github.com/savaki/go.hue"
)

func main() {
	bridge := hue.NewBridge("your-ip-address", "your-username")
	lights, _ := bridge.GetAllLights()

	for _, light := range lights {
		light.ColorLoop()
	}
}

```

### Example - Easy Access to Lights

```
package main

import (
	"github.com/savaki/go.hue"
)

func main() {
	bridge := hue.NewBridge("your-ip-address", "your-username")
	light, _ := bridge.FindLightByName("Bathroom Light")
	light.On()
}

```

