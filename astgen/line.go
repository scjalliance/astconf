package astgen

import (
	"fmt"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dpma"
)

// Lines generates DPMA line entries for a dataset.
func Lines(data *astorg.DataSet) []dpma.Line {
	lookup := data.Lookup()

	// Progressively add line entries to an ordered map
	var m dpma.LineMap

	// Step 1: Add all phones
	for _, phone := range data.Phones {
		username := phoneUsername(phone.MAC, lookup)
		label := phone.MAC
		if loc, ok := lookup.LocationByName[phone.Location]; ok {
			if loc.Abbreviation != "" {
				label = loc.Abbreviation + "-" + label
			}
		}
		line := dpma.Line{
			Name:  username,
			Label: label,
		}
		m.Add(line)
	}

	// Step 2: Overlay person configuration
	finished := make(map[string]bool) // Keep track of which phones are fully configured

	for _, person := range data.People {
		for _, mac := range person.Phones {
			if !m.Contains(mac) || finished[mac] {
				continue
			}
			username := phoneUsername(mac, lookup)
			line := dpma.Line{
				Name:  username,
				Label: fmt.Sprintf("%s %s", person.Extension, person.FullName),
			}
			m.Overlay(line)
			finished[mac] = true
		}
	}

	// Step 3: Overlay phone role configuration
	for _, role := range data.PhoneRoles {
		for _, mac := range role.Phones {
			if !m.Contains(mac) || finished[mac] {
				continue
			}
			username := phoneUsername(mac, lookup)
			line := dpma.Line{
				Name:  username,
				Label: fmt.Sprintf("%s %s", role.Extension, role.DisplayName),
			}
			m.Overlay(line)
			finished[mac] = true
		}
	}

	// Return the compiled set of line entries
	return m.Lines()
}
