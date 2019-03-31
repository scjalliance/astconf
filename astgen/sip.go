package astgen

import (
	"fmt"
	"strings"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/sip"
)

// SIP generates SIP entries for a dataset.
//
// The values in the base configuration will be merged into each entity
// that is returned.
func SIP(data *astorg.DataSet, base sip.Entity, context string) []sip.Entity {
	lookup := data.Lookup()

	// Progressively add entities to an ordered map
	var m sip.EntityMap

	// Step 1: Add all phones
	for _, phone := range data.Phones {
		username := phoneUsername(phone.MAC, lookup)
		entity := sip.Entity{
			Username: username,
			Secret:   phone.Secret,
		}
		if phone.Location != "" {
			entity.Variables = []string{
				fmt.Sprintf("USER_LOCATION=%s", phone.Location),
			}
		}
		m.Add(sip.MergeEntities(base, entity))
	}

	// Step 2: Merge person configuration
	finished := make(map[string]bool) // Keep track of which phones are fully configured

	for _, person := range data.People {
		var vars []string
		if person.Location != "" {
			//vars = append(vars, fmt.Sprintf("USER_LOCATION=%s", person.Location))
		}
		if person.VoicemailCode != "" {
			vars = append(vars, fmt.Sprintf("VMCODE=%s", person.VoicemailCode))
		}
		for _, number := range person.ContactNumbers {
			if strings.Contains(strings.ToLower(number.Label), "mobile") {
				vars = append(vars, fmt.Sprintf("MOBILE=%s", number.Dial))
				break
			}
		}
		for _, mac := range person.Phones {
			if !m.Contains(mac) || finished[mac] {
				continue
			}
			//vars := vars // New variable
			//vars = append(vars, fmt.Sprintf("USERNAME=%s", person.Username))
			username := phoneUsername(mac, lookup)
			entity := sip.Entity{
				Username:  username,
				CallerID:  fmt.Sprintf("\"%s\" <%s>", person.FullName, person.Extension),
				Mailbox:   fmt.Sprintf("%s@%s", person.Extension, context),
				Variables: vars,
			}
			m.Merge(entity)
			finished[mac] = true
		}
	}

	// Step 3: Merge phone role configuration
	for _, role := range data.PhoneRoles {
		for _, mac := range role.Phones {
			if !m.Contains(mac) || finished[mac] {
				continue
			}
			username := phoneUsername(mac, lookup)
			entity := sip.Entity{
				Username: username,
			}
			m.Merge(entity)
			finished[mac] = true
		}
	}

	// Return the compiled set of entity entries
	return m.Entities()
}
