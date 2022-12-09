package astorg

// DataSet represents a complete dataset for an organization.
type DataSet struct {
	Locations    LocationList
	People       PersonList
	PhoneRoles   PhoneRoleList
	Phones       PhoneList
	Softphones   SoftphoneList
	PagingGroups PagingGroupList
	Alerts       AlertList
	Ringtones    RingtoneList
	Mailboxes    MailboxList
	Tags         TagList
}

// Size returns the total number of records in the data set.
func (d *DataSet) Size() int {
	length := 0
	length += len(d.Locations)
	length += len(d.People)
	length += len(d.PhoneRoles)
	length += len(d.Phones)
	length += len(d.Softphones)
	length += len(d.PagingGroups)
	length += len(d.Alerts)
	length += len(d.Ringtones)
	length += len(d.Mailboxes)
	length += len(d.Tags)
	return length
}

// Lookup returns a lookup constructed from the data set.
func (d *DataSet) Lookup() Lookup {
	return Lookup{
		LocationByName:    d.Locations.ByName(),
		PersonByEmail:     d.People.ByEmail(),
		PersonByNumber:    d.People.ByExtension(),
		RoleByUsername:    d.PhoneRoles.ByUsername(),
		RoleByNumber:      d.PhoneRoles.ByExtension(),
		PhoneAssignments:  MergeAssignments(d.Phones.Assignments(), d.PhoneRoles.Assignments(), d.People.Assignments()),
		PagingGroupsByExt: d.PagingGroups.ByExtension(),
		AlertsByName:      d.Alerts.ByName(),
		RingtonesByName:   d.Ringtones.ByName(),
		MailboxesByName:   d.Mailboxes.ByName(),
		TagsByName:        d.Tags.ByName(),
	}
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
	if len(d.Softphones) != len(e.Softphones) {
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
	if len(d.Mailboxes) != len(e.Mailboxes) {
		return false
	}
	if len(d.Tags) != len(e.Tags) {
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
		if !d.Phones[i].Equal(&e.Phones[i]) {
			return false
		}
	}
	for i := range d.Softphones {
		if d.Softphones[i] != e.Softphones[i] {
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
	for i := range d.Mailboxes {
		if d.Mailboxes[i] != e.Mailboxes[i] {
			return false
		}
	}
	for i := range d.Tags {
		if !d.Tags[i].Equal(&e.Tags[i]) {
			return false
		}
	}
	return true
}
