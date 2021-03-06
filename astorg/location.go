package astorg

import "github.com/scjalliance/astconf/astoverlay"

// Location represents the location of a person or phone.
type Location struct {
	Name         string
	Abbreviation string
	Server       string
	Network      string
	Timezone     string
	CallerID     string
	AreaCode     string
	PagingGroups []string
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
	if a.Abbreviation != b.Abbreviation {
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

// OverlayLocations returns the overlayed configuration of all the given
// locations, in order of priority from least to greatest.
func OverlayLocations(locations ...Location) (overlayed Location) {
	for i := range locations {
		loc := &locations[i]
		astoverlay.String(&loc.Name, &overlayed.Name)
		astoverlay.String(&loc.Abbreviation, &overlayed.Abbreviation)
		astoverlay.String(&loc.Server, &overlayed.Server)
		astoverlay.String(&loc.Network, &overlayed.Network)
		astoverlay.String(&loc.Timezone, &overlayed.Timezone)
		astoverlay.String(&loc.CallerID, &overlayed.CallerID)
		astoverlay.String(&loc.AreaCode, &overlayed.AreaCode)
		astoverlay.StringSlice(&loc.PagingGroups, &overlayed.PagingGroups)
	}
	return
}
