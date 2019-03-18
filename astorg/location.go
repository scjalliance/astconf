package astorg

// Location represents the location of a person or phone.
type Location struct {
	Name                  string
	Server                string
	Network               string
	Timezone              string
	UnassignedPhonePrefix string
	CallerID              string
	AreaCode              string
	PagingGroups          []string
}

// Equal reports whether p and q are equal.
func (a *Location) Equal(b *Location) bool {
	// Compare nil-ness
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare simple fields
	if a.Name != b.Name {
		return false
	}
	if a.Server != b.Server {
		return false
	}
	if a.Network != b.Network {
		return false
	}
	if a.Timezone != b.Timezone {
		return false
	}
	if a.UnassignedPhonePrefix != b.UnassignedPhonePrefix {
		return false
	}
	if a.CallerID != b.CallerID {
		return false
	}
	if a.AreaCode != b.AreaCode {
		return false
	}

	// Compare slices
	if !equalStrings(a.PagingGroups, b.PagingGroups) {
		return false
	}

	return true
}
