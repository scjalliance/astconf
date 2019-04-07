package astorg

// LocationList is a slice of locations.
type LocationList []Location

// ByName returns a map of locations indexed by name.
func (locations LocationList) ByName() map[string]Location {
	lookup := make(map[string]Location, len(locations))
	for _, location := range locations {
		if location.Name == "" {
			continue
		}
		lookup[location.Name] = location
	}
	return lookup
}
