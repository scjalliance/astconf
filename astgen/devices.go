package astgen

import (
	"sort"

	"github.com/scjalliance/astconf/dialplan"
)

func makeDevices(names ...string) []dialplan.Device {
	devices := make([]dialplan.Device, 0, len(names))
	for _, name := range names {
		devices = append(devices, dialplan.SIP(name))
	}
	return devices
}

func devicesNotInUse(devices []dialplan.Device) dialplan.Expression {
	var expr dialplan.Expression
	for _, device := range devices {
		deviceNotInUse := dialplan.Equal(dialplan.DeviceState(device), dialplan.String("NOT_INUSE"))
		if expr == nil {
			expr = deviceNotInUse
		} else {
			expr = dialplan.And(expr, deviceNotInUse)
		}
	}
	return expr
}

func dedupAndSortDevices(devices []dialplan.Device) []dialplan.Device {
	if len(devices) <= 1 {
		return devices
	}

	seen := make(map[string]dialplan.Device, len(devices))
	keys := make([]string, 0, len(devices))
	for _, device := range devices {
		key := device.String()
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = device
		keys = append(keys, key)
	}

	sort.Strings(keys)

	output := make([]dialplan.Device, 0, len(keys))
	for _, key := range keys {
		output = append(output, seen[key])
	}
	return output
}
