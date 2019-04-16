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
		if username == "" {
			continue
		}
		entry := dpma.Phone{
			Username: username,
			MAC:      phone.MAC,
			FullName: username,
		}
		if line := lineUsername(phone.MAC, lookup); line != "" {
			entry.Lines = []string{line}
		}
		if loc, ok := lookup.LocationByName[phone.Location]; ok {
			if loc.Abbreviation != "" {
				entry.FullName = loc.Abbreviation + "-" + entry.FullName
			}
			entry.Timezone = loc.Timezone
			entry.Alerts = dedupStringSlice(append(base.Alerts, alertsForPagingGroups(lookup, loc.PagingGroups...)...))
			entry.Ringtones = dedupStringSlice(append(base.Ringtones, ringtonesForAlerts(lookup, entry.Alerts...)...))
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
				Username:     username,
				MAC:          mac,
				FullName:     person.FullName,
				Alerts:       person.Alerts,
				Ringtones:    person.Ringtones,
				Applications: person.Apps,
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
			if person.DefaultRingtone != "" {
				entry.ActiveRingtone = person.DefaultRingtone
			}
			applyPagingGroupsToPhone(lookup, &entry, person.PagingGroups...)
			applyTagsToPhone(lookup, &entry, person.Tags...)
			dedupPhoneData(&entry)
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
				Username:     username,
				MAC:          mac,
				FullName:     role.DisplayName,
				Applications: role.Apps,
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
			applyPagingGroupsToPhone(lookup, &entry, role.PagingGroups...)
			applyTagsToPhone(lookup, &entry, role.Tags...)
			dedupPhoneData(&entry)
			m.Merge(entry)
			finished[mac] = true
		}
	}

	// Return the compiled set of phone entries
	return m.Phones()
}

func applyPagingGroupsToPhone(lookup astorg.Lookup, entry *dpma.Phone, pagingGroups ...string) {
	entry.Alerts = append(entry.Alerts, alertsForPagingGroups(lookup, pagingGroups...)...)
	entry.Ringtones = append(entry.Ringtones, ringtonesForAlerts(lookup, entry.Alerts...)...)
}

func applyTagsToPhone(lookup astorg.Lookup, entry *dpma.Phone, tags ...string) {
	for _, t := range tags {
		if tag, ok := lookup.TagsByName[t]; ok {
			entry.Alerts = append(entry.Alerts, tag.Alerts...)
			entry.Ringtones = append(entry.Ringtones, tag.Ringtones...)
			entry.Applications = append(entry.Applications, tag.Apps...)
			applyPagingGroupsToPhone(lookup, entry, tag.PagingGroups...)
		}
	}
}
func dedupPhoneData(entry *dpma.Phone) {
	entry.Alerts = dedupStringSlice(entry.Alerts)
	entry.Ringtones = dedupStringSlice(entry.Ringtones)
	entry.Applications = dedupStringSlice(entry.Applications)
}

func buildURL(prefix, file string) string {
	u, err := url.Parse(prefix)
	if err != nil {
		return file
	}
	u.Path = path.Join(u.Path, file)
	return u.String()
}

func alertsForPagingGroups(lookup astorg.Lookup, extensions ...string) (alerts []string) {
	for _, ext := range extensions {
		if group, ok := lookup.PagingGroupsByExt[ext]; ok {
			if group.Alert == "" {
				continue
			}
			alerts = append(alerts, group.Alert)
		}
	}
	return alerts
}

func ringtonesForAlerts(lookup astorg.Lookup, alerts ...string) (ringtones []string) {
	for _, alert := range alerts {
		if alert, ok := lookup.AlertsByName[alert]; ok {
			if alert.Ringtone == "" {
				continue
			}
			if _, ok := lookup.RingtonesByName[alert.Ringtone]; !ok {
				// Not having a ringtone entry implies that it's a built-in ringtone
				continue
			}
			ringtones = append(ringtones, alert.Ringtone)
		}
	}
	return ringtones
}
