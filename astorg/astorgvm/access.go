package astorgvm

// Access is a voicemail access mode for a person.
//
// The zero value of Access is a valid Default value that is equivalent to
// PhoneAndEmail. The only difference between Default and PhoneAndEmail is
// that Default values will not take precedence when access modes are
// overlayed.
type Access int

// Voicemail access modes.
const (
	Default       Access = 0
	Phone         Access = 1
	Email         Access = 2
	PhoneAndEmail Access = 3
	Disabled      Access = 4
)

// Specified returns true if a holds a non-default value.
func (a Access) Specified() bool {
	return a != Default
}

// String returns a string representation of the voicemail access mode.
func (a Access) String() string {
	switch a {
	case Default, PhoneAndEmail:
		return "phone+email"
	case Phone:
		return "phone"
	case Email:
		return "email"
	case Disabled:
		return "disabled"
	default:
		return ""
	}
}
