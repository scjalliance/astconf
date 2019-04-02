package astgen

import "github.com/scjalliance/astconf/astorg"

func phoneUsername(mac string, lookup astorg.Lookup) string {
	return mac
}

func lineUsername(mac string, lookup astorg.Lookup) string {
	if assignment, ok := lookup.PhoneAssignments[mac]; ok && assignment.Username != "" {
		return assignment.Username
	}
	return mac
}
