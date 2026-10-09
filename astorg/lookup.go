package astorg

// Lookup holds various lookup maps for a data set.
type Lookup struct {
	ServerByID        ServerMap              // Maps IDs to servers
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
	TagsByName        map[string]Tag         // Maps tag names to tags
}

// LocationServer returns the server that serves the named location.
//
// The returned id is the location's Server field. It is empty when the
// location is unknown or has no server, which is normal for an office with
// no desk phones. When id is not empty, ok reports whether the data set
// contains a server with that ID; a false ok then means the location names
// a server that is missing.
func (l Lookup) LocationServer(location string) (id string, server Server, ok bool) {
	loc, found := l.LocationByName[location]
	if !found || loc.Server == "" {
		return "", Server{}, false
	}
	server, ok = l.ServerByID[loc.Server]
	return loc.Server, server, ok
}

// Hub returns the hub server.
//
// It reports false unless exactly one server in the lookup is the hub, so
// a data set with no hub or with several hubs has no hub.
func (l Lookup) Hub() (Server, bool) {
	var (
		hub   Server
		count int
	)
	for _, server := range l.ServerByID {
		if server.Hub {
			hub = server
			count++
		}
	}
	if count != 1 {
		return Server{}, false
	}
	return hub, true
}
