package astgen

import (
	"fmt"
	"sort"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// PhoneExtensions generates a dialplan section that includes all phone
// extensions in a dataset.
func PhoneExtensions(data *astorg.DataSet, context string) dialplan.Section {
	section := dialplan.Section{Context: context}

	for _, person := range data.People {
		if person.Extension == "" || person.Username == "" {
			continue
		}

		device := dialplan.SIP(person.Username)
		presence := fmt.Sprintf("CustomPresence:%s", person.Username)

		ext := dialplan.Extension{
			Comment: person.FullName,
			Number:  person.Extension,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Call %s", person.Username)),
				dialplan.Macro("pre-answer", dialplan.Var("EXTEN")),
				dialplan.ExecIf(dialplan.Equal(dialplan.DeviceState(device), dialplan.String("NOT_INUSE")), dialplan.Dial(device, 20)),
				dialplan.ExecIf(dialplan.Equal(dialplan.PresenceState(presence, "subtype"), dialplan.String("with Call Waiting Enabled")), dialplan.Dial(device, 20)),
				dialplan.Macro("no-answer", dialplan.Var("EXTEN")),
			},
		}

		section.Extensions = append(section.Extensions, ext)
	}

	for _, role := range data.PhoneRoles {
		if role.Extension == "" || role.Username == "" {
			continue
		}

		device := dialplan.SIP(role.Username)
		presence := fmt.Sprintf("CustomPresence:%s", role.Username)

		ext := dialplan.Extension{
			Comment: role.DisplayName,
			Number:  role.Extension,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Call %s", role.Username)),
				dialplan.Macro("pre-answer", dialplan.Var("EXTEN")),
				dialplan.ExecIf(dialplan.Equal(dialplan.DeviceState(device), dialplan.String("NOT_INUSE")), dialplan.Dial(device, 20)),
				dialplan.ExecIf(dialplan.Equal(dialplan.PresenceState(presence, "subtype"), dialplan.String("with Call Waiting Enabled")), dialplan.Dial(device, 20)),
				dialplan.Congestion(),
			},
		}

		section.Extensions = append(section.Extensions, ext)
	}

	return section
}

// IntercomExtensions generates a dialplan section that includes all intercom
// extensions in a dataset.
func IntercomExtensions(data *astorg.DataSet, context string) dialplan.Section {
	lookup := data.Lookup()

	// Build up the set of devices contained in each paging group
	membership := make(map[string][]dialplan.Device)

	// Step 1: Add phones in locations with paging groups
	for _, phone := range data.Phones {
		username := phoneUsername(phone.MAC, lookup)
		device := dialplan.SIP(username)
		if loc, ok := lookup.LocationByName[phone.Location]; ok {
			for _, group := range loc.PagingGroups {
				if members, ok := membership[group]; ok {
					membership[group] = append(members, device)
				} else {
					membership[group] = []dialplan.Device{device}
				}
			}
		}
	}

	// Step 2: Add phones assigned to people with paging groups
	finished := make(map[string]bool) // Keep track of which phones are fully configured

	for _, person := range data.People {
		for _, mac := range person.Phones {
			if finished[mac] {
				continue
			}
			username := phoneUsername(mac, lookup)
			device := dialplan.SIP(username)
			for _, group := range person.PagingGroups {
				if members, ok := membership[group]; ok {
					membership[group] = append(members, device)
				} else {
					membership[group] = []dialplan.Device{device}
				}
			}
			finished[mac] = true
		}
	}

	// Step 3: Add phones assigned to roles with paging groups
	for _, role := range data.PhoneRoles {
		for _, mac := range role.Phones {
			if finished[mac] {
				continue
			}
			username := phoneUsername(mac, lookup)
			device := dialplan.SIP(username)
			for _, group := range role.PagingGroups {
				if members, ok := membership[group]; ok {
					membership[group] = append(members, device)
				} else {
					membership[group] = []dialplan.Device{device}
				}
			}
			finished[mac] = true
		}
	}

	// Create a sorted set of paging groups
	groups := make([]string, 0, len(membership))
	for group := range membership {
		groups = append(groups, group)
	}
	sort.Strings(groups)

	// Build an extension for each paging group
	section := dialplan.Section{Context: context}

	for _, group := range groups {
		devices := dedupAndSortDevices(membership[group])

		name := dialplan.CallerID("name")
		number := dialplan.CallerID("num")
		intercom := dialplan.String("Intercom")
		zero := dialplan.Int(0)

		ext1 := dialplan.Extension{
			Number: group,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Intercom %s", group)),
				dialplan.Playback("beep"),
				dialplan.SIPAddHeader("Alert-Info", "beep-answer"),
				dialplan.Set(name, dialplan.IfElse(
					dialplan.GreaterThan(dialplan.Len(number), zero),
					name,
					intercom,
				)),
				dialplan.Set(number, dialplan.IfElse(
					dialplan.GreaterThan(dialplan.Len(number), zero),
					number,
					intercom,
				)),
				dialplan.PageApp{
					Recipients:   devices,
					Announcement: "beep",
				},
				dialplan.Hangup(),
			},
		}

		ext2 := dialplan.Extension{
			Number: "*" + group,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Intercom %s", group)),
				dialplan.Playback("beep"),
				dialplan.SIPAddHeader("Alert-Info", "intercom"),
				dialplan.Set(name, dialplan.IfElse(
					dialplan.GreaterThan(dialplan.Len(number), zero),
					name,
					intercom,
				)),
				dialplan.Set(number, dialplan.IfElse(
					dialplan.GreaterThan(dialplan.Len(number), zero),
					number,
					intercom,
				)),
				dialplan.Page(devices...),
				dialplan.Hangup(),
			},
		}

		section.Extensions = append(section.Extensions, ext1, ext2)
	}

	return section
}
