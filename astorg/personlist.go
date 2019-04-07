package astorg

// PersonList is a slice of people.
type PersonList []Person

// ByEmail returns a map of people indexed by email address.
func (people PersonList) ByEmail() map[string]Person {
	lookup := make(map[string]Person, len(people))
	for _, person := range people {
		for _, email := range person.EmailAddresses {
			if email.Address == "" {
				continue
			}
			lookup[email.Address] = person
		}
	}
	return lookup
}

// ByExtension returns a map of people indexed by extension.
func (people PersonList) ByExtension() map[string]Person {
	lookup := make(map[string]Person, len(people))
	for _, person := range people {
		if person.Extension == "" {
			continue
		}
		lookup[person.Extension] = person
	}
	return lookup
}

// Assignments returns a phone assignment map for people, indexed by
// mac address.
func (people PersonList) Assignments() AssignmentMap {
	assignments := make(AssignmentMap, len(people))
	for _, person := range people {
		for _, mac := range person.Phones {
			assignments.Add(mac, Assignment{Type: PersonAssigned, Username: person.Username})
		}
	}
	return assignments
}
