package astorg

// Server represents one PBX node: an office server or the hub.
//
// A location names the server that serves it through Location.Server, which
// holds the server's ID. The server does not list its locations, so the
// mapping is stored in one place.
type Server struct {
	ID      string // Stable identifier, referenced by Location.Server
	Name    string
	Hub     bool   // The node that carries the PSTN trunks and central voicemail
	Address string // The address other servers use to reach this one
}

// Equal reports whether a and b are equal.
func (a *Server) Equal(b *Server) bool {
	// Compare nil-ness
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare simple fields
	if a.ID != b.ID {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Hub != b.Hub {
		return false
	}
	if a.Address != b.Address {
		return false
	}

	return true
}
