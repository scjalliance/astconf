package astgen

import (
	"fmt"
	"net/url"
	"path"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dpma"
)

// Phones generates DPMA phone entries for a dataset.
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
			Lines:    []string{username},
			Timezone: tz,
		}
		m.Add(dpma.OverlayPhones(base, entry))
	}

	// Step 2: Add person configuration
	finished := make(map[string]bool) // Keep track of which phones are fully configured

	for _, person := range data.People {
		for _, mac := range person.Phones {
			username := phoneUsername(mac, lookup)
			if finished[mac] {
				continue
			}
			entry := dpma.Phone{
				Username:  username,
				MAC:       mac,
				FullName:  person.FullName,
				Ringtones: person.Ringtones,
				Lines:     []string{username},
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
			if !m.Contains(mac) {
				m.Add(dpma.OverlayPhones(base, entry))
			} else {
				m.Merge(entry)
			}
			finished[mac] = true
		}
	}

	// Step 3: Add phone role configuration
	for _, role := range data.PhoneRoles {
		for _, mac := range role.Phones {
			if finished[mac] {
				continue
			}
			username := phoneUsername(mac, lookup)
			entry := dpma.Phone{
				Username: username,
				MAC:      mac,
				FullName: role.DisplayName,
				Lines:    []string{username},
			}
			if role.Username != "" {
				entry.Contacts = []string{
					buildURL(contactsURL, "global.xml"),
					buildURL(contactsURL, fmt.Sprintf("%s.xml", role.Username)),
				}
				entry.BLFItems = buildURL(contactsURL, fmt.Sprintf("%s.blf.xml", role.Username))
			}
			if !m.Contains(mac) {
				m.Add(dpma.OverlayPhones(base, entry))
			} else {
				m.Merge(entry)
			}
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
