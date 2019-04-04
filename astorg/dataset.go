package astorg

import "strings"

// DataSet represents a complete dataset for an organization.
type DataSet struct {
	Locations    []Location
	People       []Person
	PhoneRoles   []PhoneRole
	Phones       []Phone
	PagingGroups []PagingGroup
	Alerts       []Alert
	Ringtones    []Ringtone
}

// Size returns the total number of records in the data set.
func (d *DataSet) Size() int {
	length := 0
	length += len(d.Locations)
	length += len(d.People)
	length += len(d.PhoneRoles)
	length += len(d.Phones)
	length += len(d.PagingGroups)
	length += len(d.Alerts)
	length += len(d.Ringtones)
	return length
}

// Lookup returns a lookup constructed from the data set.
func (d *DataSet) Lookup() Lookup {
	lookup := Lookup{
		LocationByName:    make(map[string]Location),
		PersonByEmail:     make(map[string]Person),
		PersonByNumber:    make(map[string]Person),
		RoleByUsername:    make(map[string]PhoneRole),
		RoleByNumber:      make(map[string]PhoneRole),
		PhoneAssignments:  make(map[string]Assignment),
		PagingGroupsByExt: make(map[string]PagingGroup),
		AlertsByName:      make(map[string]Alert),
		RingtonesByNAme:   make(map[string]Ringtone),
	}
	for _, location := range d.Locations {
		if location.Name != "" {
			lookup.LocationByName[location.Name] = location
		}
	}
	for _, phone := range d.Phones {
		if _, found := lookup.PhoneAssignments[phone.MAC]; found {
			continue
		}
		lookup.PhoneAssignments[phone.MAC] = Assignment{Type: Unassigned, Username: phone.MAC}
	}
	for _, role := range d.PhoneRoles {
		if role.Username != "" {
			lookup.RoleByUsername[role.Username] = role
		}
		if role.Extension != "" {
			lookup.RoleByNumber[role.Extension] = role
		}
		for _, mac := range role.Phones {
			if assignment := lookup.PhoneAssignments[mac]; assignment.Type != Unassigned {
				continue
			}
			lookup.PhoneAssignments[mac] = Assignment{Type: RoleAssigned, Username: role.Username}
		}
	}
	for _, person := range d.People {
		for _, address := range person.EmailAddresses {
			lookup.PersonByEmail[strings.ToLower(address.Address)] = person
		}
		if person.Extension != "" {
			lookup.PersonByNumber[person.Extension] = person
		}
		for _, mac := range person.Phones {
			if assignment := lookup.PhoneAssignments[mac]; assignment.Type != Unassigned {
				continue
			}
			lookup.PhoneAssignments[mac] = Assignment{Type: PersonAssigned, Username: person.Username}
		}
	}
	for _, group := range d.PagingGroups {
		if group.Extension != "" {
			lookup.PagingGroupsByExt[group.Extension] = group
		}
	}
	for _, alert := range d.Alerts {
		if alert.Name != "" {
			lookup.AlertsByName[alert.Name] = alert
		}
	}
	for _, ringtone := range d.Ringtones {
		if ringtone.Name != "" {
			lookup.RingtonesByNAme[ringtone.Name] = ringtone
		}
	}
	return lookup
}

// Equal reports whether d and e are equal.
func (d *DataSet) Equal(e *DataSet) bool {
	// Compare slice lengths
	if len(d.Locations) != len(e.Locations) {
		return false
	}
	if len(d.People) != len(e.People) {
		return false
	}
	if len(d.PhoneRoles) != len(e.PhoneRoles) {
		return false
	}
	if len(d.Phones) != len(e.Phones) {
		return false
	}
	if len(d.PagingGroups) != len(e.PagingGroups) {
		return false
	}
	if len(d.Alerts) != len(e.Alerts) {
		return false
	}
	if len(d.Ringtones) != len(e.Ringtones) {
		return false
	}

	// Compare slice values
	for i := range d.Locations {
		if d.Locations[i].Equal(&e.Locations[i]) {
			return false
		}
	}
	for i := range d.People {
		if !d.People[i].Equal(&e.People[i]) {
			return false
		}
	}
	for i := range d.PhoneRoles {
		if !d.PhoneRoles[i].Equal(&e.PhoneRoles[i]) {
			return false
		}
	}
	for i := range d.Phones {
		if d.Phones[i] != e.Phones[i] {
			return false
		}
	}
	for i := range d.PagingGroups {
		if d.PagingGroups[i] != e.PagingGroups[i] {
			return false
		}
	}
	for i := range d.Alerts {
		if d.Alerts[i] != e.Alerts[i] {
			return false
		}
	}
	for i := range d.Ringtones {
		if d.Ringtones[i] != e.Ringtones[i] {
			return false
		}
	}
	return true
}
