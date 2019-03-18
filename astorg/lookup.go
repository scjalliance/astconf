package astorg

// Lookup holds various lookup maps for a data set.
type Lookup struct {
	LocationByName map[string]Location  // Maps names to locations
	PersonByEmail  map[string]Person    // Maps email addresses to people
	PersonByNumber map[string]Person    // Maps phone numbers to people
	RoleByUsername map[string]PhoneRole // Maps usernames to phone roles
	RoleByNumber   map[string]PhoneRole // Maps phone numbers to phone roles
}
