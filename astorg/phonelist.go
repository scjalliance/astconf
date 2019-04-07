package astorg

// PhoneList is a slice of phones.
type PhoneList []Phone

// Assignments returns a phone assignment map for phones, indexed by
// mac address.
func (phones PhoneList) Assignments() AssignmentMap {
	assignments := make(AssignmentMap, len(phones))
	for _, phone := range phones {
		assignments.Add(phone.MAC, Assignment{Type: Unassigned, Username: phone.MAC})
	}
	return assignments
}
