package dialplan

// Device is an addressable resource within a channel.
type Device struct {
	Technology string // Channel Driver
	Resource   string
}

// String returns a string representation of the device address.
func (d Device) String() string {
	return d.Technology + "/" + d.Resource
}

// PJSIP returns a PJSIP device address.
func PJSIP(resource string) Device {
	return Device{Technology: "PJSIP", Resource: resource}
}

// SIP returns a SIP device address.
func SIP(resource string) Device {
	return Device{Technology: "SIP", Resource: resource}
}
