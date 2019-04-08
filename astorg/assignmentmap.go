package astorg

// AssignmentMap is a map of phone assignments indexed a key, which is
// typically the phone's mac address. The map is used to progressively
// tracks phone assignments; it incorporates new values according to
// a prioritized order of assignment types:
//
//  1. People
//  2. Roles
//  3. Unassigned
//
// When more than one entity of the same assignment type claims a phone,
// the first entity added to the map will hold the claim.
type AssignmentMap map[string]Assignment

// Add attempts to add a phone assignment to the map for the given key.
// The assignment will succeed if an existing assignment of equal or
// greater priority does not already exist in the map with that key.
//
// Add returns true if the assignment was added successfully.
func (m AssignmentMap) Add(key string, assignment Assignment) bool {
	if key == "" {
		return false
	}
	if existing, found := m[key]; !found {
		m[key] = assignment
		return true
	} else if existing.Type < assignment.Type {
		m[key] = assignment
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
		for key, assignment := range m {
			combined.Add(key, assignment)
		}
	}
	return combined
}
