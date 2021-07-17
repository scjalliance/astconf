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
	var deviceStates []dialplan.Expression
	for _, device := range devices {
		deviceNotInUse := dialplan.Or(
			dialplan.Equal(dialplan.DeviceState(device), dialplan.String("NOT_INUSE")),
			dialplan.Equal(dialplan.DeviceState(device), dialplan.String("UNAVAILABLE")),
			dialplan.Equal(dialplan.DeviceState(device), dialplan.String("INVALID")),
		)
		deviceStates = append(deviceStates, deviceNotInUse)
	}
	switch len(deviceStates) {
	case 0:
		return nil
	case 1:
		return deviceStates[0]
	case 2:
		return dialplan.And(deviceStates[0], deviceStates[1])
	default:
		return dialplan.And(deviceStates[0], deviceStates[1], deviceStates[2:]...)
	}
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
