package astgen

import (
	"fmt"
	"strings"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
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
		username := lineUsername(phone.MAC, lookup)
		var vars []astval.Var
		var callerID string
		if phone.Location != "" {
			vars = append(vars, astval.NewVar("USER_LOCATION", phone.Location))
			if loc, ok := lookup.LocationByName[phone.Location]; ok {
				vars = append(vars, astval.NewVar("OUTBOUND_CALLERID", loc.CallerID))
				vars = append(vars, astval.NewVar("AREACODE", loc.AreaCode))
				if loc.Abbreviation != "" {
					callerID = fmt.Sprintf("\"%s-%s\" <UNAVAILABLE>", loc.Abbreviation, strings.ToUpper(phone.MAC))
				}
			}
		}
		entity := sip.Entity{
			Username:  username,
			Secret:    phone.Secret,
			CallerID:  callerID,
			Variables: vars,
		}
		m.Add(sip.MergeEntities(base, entity))
	}

	// Step 2: Merge person configuration
	finished := make(map[string]bool) // Keep track of which phones are fully configured

	for _, person := range data.People {
		var vars []astval.Var
		if person.VoicemailCode != "" {
			vars = append(vars, astval.NewVar("VMCODE", person.VoicemailCode))
		}
		for _, number := range person.ContactNumbers {
			if strings.Contains(strings.ToLower(number.Label), "mobile") {
				vars = append(vars, astval.NewVar("MOBILE", number.Dial))
				break
			}
		}
		vars = append(vars, astval.NewVar("USERNAME", person.Username))
		for _, mac := range person.Phones {
			username := lineUsername(mac, lookup)
			if !m.Contains(username) || finished[mac] {
				continue
			}
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
			username := lineUsername(mac, lookup)
			if !m.Contains(username) || finished[mac] {
				continue
			}
			entity := sip.Entity{
				Username: username,
				CallerID: fmt.Sprintf("\"%s\" <%s>", role.DisplayName, role.Extension),
			}
			m.Merge(entity)
			finished[mac] = true
		}
	}

	// Return the compiled set of entity entries
	return m.Entities()
}
