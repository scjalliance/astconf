package astorg

// Lookup holds various lookup maps for a data set.
type Lookup struct {
	LocationByName   map[string]Location   // Maps names to locations
	PersonByEmail    map[string]Person     // Maps email addresses to people
	PersonByNumber   map[string]Person     // Maps phone numbers to people
	RoleByUsername   map[string]PhoneRole  // Maps usernames to phone roles
	RoleByNumber     map[string]PhoneRole  // Maps phone numbers to phone roles
	PhoneAssignments map[string]Assignment // Maps phone mac addresses to their highest priority assignments
	AlertsByName     map[string]Alert      // Maps alert names to alerts
	RingtonesByNAme  map[string]Ringtone   // Maps ringtone names to ringtones
}
