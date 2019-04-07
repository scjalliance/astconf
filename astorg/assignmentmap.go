package astorg

// AssignmentMap is a map of phone assignments indexed by mac address.
// It tracks phone assignments and incorporates new addresses according to
// a hierarchy of preferred assignment types:
//
//  1. People
//  2. Roles
//  3. Unassigned
//
// When more than entity of the same assignment type attempts to claim a phone,
// the first entity added to the map will hold the claim.
type AssignmentMap map[string]Assignment

// Add attempts to add a phone assignment to the map for the given mac address.
// The assignment will succeed if an existing assignment of equal or
// greater priority does not already exist in the map.
//
// Add returns true if the assignment was added successfully.
func (m AssignmentMap) Add(mac string, assignment Assignment) bool {
	if mac == "" {
		return false
	}
	if existing, found := m[mac]; !found {
		m[mac] = assignment
		return true
	} else if existing.Type < assignment.Type {
		m[mac] = assignment
		return true
	}
	return false
}

// MergeAssignments merges any number of assignment maps according to
// the assignment priority rules.
func MergeAssignments(maps ...AssignmentMap) AssignmentMap {
	count := 0
	for _, m := range maps {
		count += len(m)
	}
	combined := make(AssignmentMap, count)
	for _, m := range maps {
		for mac, assignment := range m {
			combined.Add(mac, assignment)
		}
	}
	return combined
}
