package astorg

// ServerList is a slice of servers.
type ServerList []Server

// ByID returns a map of servers indexed by ID. When more than one server
// has the same ID, the first one in the list is kept.
func (servers ServerList) ByID() ServerMap {
	lookup := make(ServerMap, len(servers))
	for _, server := range servers {
		lookup.Add(server)
	}
	return lookup
}
