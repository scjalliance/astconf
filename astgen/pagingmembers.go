package astgen

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// pagingMembers returns the device names in each paging group, keyed by
// the group's extension. A name appears in the order it was found and may
// appear more than once; callers sort and deduplicate.
//
// Membership is found in three steps. Every phone in a location joins the
// location's paging groups. Then each phone and softphone of a person or
// phone role joins the groups of its first owner, including the groups of
// the owner's tags.
func pagingMembers(data *astorg.DataSet) map[string][]string {
	lookup := data.Lookup()
	membership := make(map[string][]string)

	addDevice := func(name string, groups []string) {
		for _, group := range groups {
			membership[group] = append(membership[group], name)
		}
	}

	// Keep track of which devices have been assigned
	assigned := make(map[string]bool)

	assignOwner := func(groups, tags, phones, softphones []string) {
		for _, t := range tags {
			if tag, ok := lookup.TagsByName[t]; ok {
				groups = append(groups, tag.PagingGroups...)
			}
		}
		names := make([]string, 0, len(phones)+len(softphones))
		for _, mac := range phones {
			names = append(names, lineUsername(mac, lookup))
		}
		names = append(names, softphones...)
		for _, name := range names {
			if assigned[name] {
				continue
			}
			assigned[name] = true
			addDevice(name, groups)
		}
	}

	// Step 1: Add phones in locations with paging groups
	for _, phone := range data.Phones {
		if loc, ok := lookup.LocationByName[phone.Location]; ok {
			addDevice(lineUsername(phone.MAC, lookup), loc.PagingGroups)
		}
	}

	// Step 2: Add phones assigned to people with paging groups
	for _, person := range data.People {
		assignOwner(person.PagingGroups, person.Tags, person.Phones, person.Softphones)
	}

	// Step 3: Add phones assigned to roles with paging groups
	for _, role := range data.PhoneRoles {
		assignOwner(role.PagingGroups, role.Tags, role.Phones, role.Softphones)
	}

	return membership
}

// pageCallerIDActions returns the actions of a paging group extension that
// fill in a missing caller ID name and number with "Intercom".
func pageCallerIDActions() []dialplan.Action {
	name := dialplan.CallerID("name")
	number := dialplan.CallerID("num")
	intercom := dialplan.String("Intercom")
	zero := dialplan.Int(0)

	return []dialplan.Action{
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
	}
}
