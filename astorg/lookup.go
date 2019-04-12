package astorg

// Lookup holds various lookup maps for a data set.
type Lookup struct {
	LocationByName    map[string]Location    // Maps names to locations
	PersonByEmail     map[string]Person      // Maps email addresses to people
	PersonByNumber    map[string]Person      // Maps phone numbers to people
	RoleByUsername    map[string]PhoneRole   // Maps usernames to phone roles
	RoleByNumber      map[string]PhoneRole   // Maps phone numbers to phone roles
	PhoneAssignments  AssignmentMap          // Maps phone mac addresses to their highest priority assignments
	PagingGroupsByExt map[string]PagingGroup // Maps extensions to paging groups
	AlertsByName      map[string]Alert       // Maps alert names to alerts
	RingtonesByName   map[string]Ringtone    // Maps ringtone names to ringtones
	MailboxesByName   map[string]Mailbox     // Maps mailbox names to mailboxes
}
