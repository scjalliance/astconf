package astgen

import (
	"fmt"
	"net/url"
	"path"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dpma"
)

// Phones generates DPMA phone entries for a dataset.
//
// The values in the base configuration will be merged into each entry
// that is returned.
//
// The contactsURL, if non-empty, will be used to prefix the contacts URL
// and BLF Items URL for each entry.
func Phones(data *astorg.DataSet, base dpma.Phone, contactsURL string) []dpma.Phone {
	lookup := data.Lookup()

	// Progressively add phone entries to an ordered map
	var m dpma.PhoneMap

	// Step 1: Add all phones
	for _, phone := range data.Phones {
		username := phoneUsername(phone.MAC, lookup)
		fullName := username
		tz := ""
		if loc, ok := lookup.LocationByName[phone.Location]; ok {
			if loc.Abbreviation != "" {
				fullName = loc.Abbreviation + "-" + fullName
			}
			tz = loc.Timezone
		}
		entry := dpma.Phone{
			Username: username,
			MAC:      phone.MAC,
			FullName: fullName,
			Timezone: tz,
		}
		if line := lineUsername(phone.MAC, lookup); line != "" {
			entry.Lines = []string{line}
		}
		m.Add(dpma.OverlayPhones(base, entry))
	}

	// Step 2: Merge person configuration
	finished := make(map[string]bool) // Keep track of which phones are fully configured

	for _, person := range data.People {
		for _, mac := range person.Phones {
			username := phoneUsername(mac, lookup)
			if !m.Contains(username) || finished[mac] {
				continue
			}
			entry := dpma.Phone{
				Username:  username,
				MAC:       mac,
				FullName:  person.FullName,
				Ringtones: person.Ringtones,
			}
			if line := lineUsername(mac, lookup); line != "" {
				entry.Lines = []string{line}
			}
			if person.Username != "" {
				entry.Contacts = []string{
					buildURL(contactsURL, "global.xml"),
					buildURL(contactsURL, fmt.Sprintf("%s.xml", person.Username)),
				}
				entry.BLFItems = buildURL(contactsURL, fmt.Sprintf("%s.blf.xml", person.Username))
			}
			if person.Firmware != "" {
				entry.Firmware = []string{person.Firmware}
			}
			m.Merge(entry)
			finished[mac] = true
		}
	}

	// Step 3: Merge phone role configuration
	for _, role := range data.PhoneRoles {
		for _, mac := range role.Phones {
			username := phoneUsername(mac, lookup)
			if !m.Contains(username) || finished[mac] {
				continue
			}
			entry := dpma.Phone{
				Username: username,
				MAC:      mac,
				FullName: role.DisplayName,
			}
			if line := lineUsername(mac, lookup); line != "" {
				entry.Lines = []string{line}
			}
			if role.Username != "" {
				entry.Contacts = []string{
					buildURL(contactsURL, "global.xml"),
					buildURL(contactsURL, fmt.Sprintf("%s.xml", role.Username)),
				}
				entry.BLFItems = buildURL(contactsURL, fmt.Sprintf("%s.blf.xml", role.Username))
			}
			m.Merge(entry)
			finished[mac] = true
		}
	}

	// Return the compiled set of phone entries
	return m.Phones()
}

func buildURL(prefix, file string) string {
	u, err := url.Parse(prefix)
	if err != nil {
		return file
	}
	u.Path = path.Join(u.Path, file)
	return u.String()
}
