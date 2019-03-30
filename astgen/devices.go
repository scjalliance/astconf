package astgen

import (
	"sort"

	"github.com/scjalliance/astconf/dialplan"
)

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
