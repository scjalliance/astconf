package astorg

// SoftphoneList is a slice of software phones.
type SoftphoneList []Softphone

// Assignments returns a phone assignment map for software phones, indexed by
// softphone username.
func (phones SoftphoneList) Assignments() AssignmentMap {
	assignments := make(AssignmentMap, len(phones))
	for _, phone := range phones {
		assignments.Add(phone.Username, Assignment{Type: Unassigned, Username: phone.Username})
	}
	return assignments
}
