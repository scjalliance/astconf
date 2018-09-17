package astorg

// DataSet represents a complete dataset for an organization.
type DataSet struct {
	Locations  []Location
	People     []Person
	PhoneRoles []PhoneRole
	Phones     []Phone
}

// Size returns the total number of records in the data set.
func (d *DataSet) Size() int {
	length := 0
	length += len(d.Locations)
	length += len(d.People)
	length += len(d.PhoneRoles)
	length += len(d.Phones)
	return length
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

	// Compare slice values
	for i := range d.Locations {
		if d.Locations[i] != e.Locations[i] {
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
	return true
}
