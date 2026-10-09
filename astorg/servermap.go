package astorg

// ServerMap is a map of servers indexed by ID.
//
// When more than one server claims an ID, the first server added to the map
// holds it.
type ServerMap map[string]Server

// Add attempts to add a server to the map. The server is not added if its ID
// is empty or the map already holds a server with that ID.
//
// Add returns true if the server was added successfully.
func (m ServerMap) Add(server Server) bool {
	if server.ID == "" {
		return false
	}
	if _, found := m[server.ID]; found {
		return false
	}
	m[server.ID] = server
	return true
}
