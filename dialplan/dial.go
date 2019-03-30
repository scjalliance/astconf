package dialplan

import (
	"strconv"
	"strings"
)

// Dial returns a dial application that dials a single device.
func Dial(device Device, seconds int) DialApp {
	return DialApp{
		Devices: []Device{device},
		Timeout: seconds,
	}
}

// DialMany returns a dial application that dials one or more devices.
func DialMany(devices []Device, seconds int) DialApp {
	return DialApp{
		Devices: devices,
		Timeout: seconds,
	}
}

// DialApp is a dialplan application that dials other devices.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Dial
type DialApp struct {
	Devices []Device
	Timeout int // Seconds
}

// App returnes the assembled application call.
func (dial DialApp) App() AppCall {
	devices := make([]string, 0, len(dial.Devices))
	for _, device := range dial.Devices {
		devices = append(devices, device.String())
	}
	app := AppCall{
		Name: "Dial",
		Args: []string{strings.Join(devices, "&")},
	}
	if dial.Timeout > 0 {
		app.Args = append(app.Args, strconv.Itoa(dial.Timeout))
	}
	return app
}
