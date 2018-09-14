package astorg

// Location represents the location of a person or phone.
type Location struct {
	Name                  string
	Network               string
	Timezone              string
	UnassignedPhonePrefix string
	Server                string
}
