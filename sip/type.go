package sip

// Type is a sip entity type.
//
// The zero value of Type is a valid Default value that is equivalent to Peer.
// The only difference between Default and Peer is that Default values will
// not take precedence when types are overlayed.
type Type int

// Sip entity types.
const (
	Default Type = 0
	Peer    Type = 1
	User    Type = 2
	Friend  Type = 3
)

// Specified returns true if t holds a non-default value.
func (t Type) Specified() bool {
	return t != Default
}

// String returns a string representation of the entity type.
func (t Type) String() string {
	switch t {
	case Default, Peer:
		return "peer"
	case User:
		return "user"
	case Friend:
		return "friend"
	default:
		return ""
	}
}

// MarshalText returns a string representation of t marshaled as utf-8 bytes.
// It never returns an error.
func (t Type) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}
