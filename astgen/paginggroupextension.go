package astgen

import (
	"fmt"
	"sort"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// PagingGroupExtensions generates a dialplan section that includes all
// paging group extensions in a dataset.
func PagingGroupExtensions(data *astorg.DataSet, context string) dialplan.Section {
	lookup := data.Lookup()

	// Build up the set of devices contained in each paging group
	membership := make(map[string][]dialplan.Device)

	addDevice := func(device dialplan.Device, groups []string) {
		for _, group := range groups {
			if members, ok := membership[group]; ok {
				membership[group] = append(members, device)
			} else {
				membership[group] = []dialplan.Device{device}
			}
		}
	}

	// Keep track of which devices have been assigned
	assigned := make(map[dialplan.Device]bool)

	assignDevice := func(device dialplan.Device, groups []string) {
		if assigned[device] {
			return
		}
		assigned[device] = true
		addDevice(device, groups)
	}

	// Step 1: Add phones in locations with paging groups
	for _, phone := range data.Phones {
		if loc, ok := lookup.LocationByName[phone.Location]; ok {
			addDevice(dialplan.SIP(lineUsername(phone.MAC, lookup)), loc.PagingGroups)
		}
	}

	// Step 2: Add phones assigned to people with paging groups
	for _, person := range data.People {
		groups := person.PagingGroups
		for _, t := range person.Tags {
			if tag, ok := lookup.TagsByName[t]; ok {
				groups = append(groups, tag.PagingGroups...)
			}
		}
		for _, mac := range person.Phones {
			assignDevice(dialplan.SIP(lineUsername(mac, lookup)), groups)
		}
		for _, entity := range person.Softphones {
			assignDevice(dialplan.SIP(entity), groups)
		}
	}

	// Step 3: Add phones assigned to roles with paging groups
	for _, role := range data.PhoneRoles {
		groups := role.PagingGroups
		for _, t := range role.Tags {
			if tag, ok := lookup.TagsByName[t]; ok {
				groups = append(groups, tag.PagingGroups...)
			}
		}
		for _, mac := range role.Phones {
			assignDevice(dialplan.SIP(lineUsername(mac, lookup)), groups)
		}
		for _, entity := range role.Softphones {
			assignDevice(dialplan.SIP(entity), groups)
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

	for _, group := range data.PagingGroups {
		devices := dedupAndSortDevices(membership[group.Extension])

		name := dialplan.CallerID("name")
		number := dialplan.CallerID("num")
		intercom := dialplan.String("Intercom")
		zero := dialplan.Int(0)

		actions := []dialplan.Action{
			dialplan.Noop(fmt.Sprintf("Paging %s", group.Name)),
			dialplan.Playback("beep"),
		}
		if group.Alert != "" {
			actions = append(actions, dialplan.SIPAddHeader("Alert-Info", group.Alert))
		}
		actions = append(actions, []dialplan.Action{
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
		}...)
		if len(devices) > 0 {
			actions = append(actions, dialplan.Page(devices...))
		}
		actions = append(actions, dialplan.Hangup())

		ext := dialplan.Extension{
			Comment: group.Name,
			Number:  group.Extension,
			Actions: actions,
		}

		section.Extensions = append(section.Extensions, ext)
	}

	return section
}
