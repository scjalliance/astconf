package astval

// Boolean constants.
const (
	ternaryUnspecified ternary = 0
	ternaryTrue        ternary = 1 // For convenience, but not the only value for true
	ternaryFalse       ternary = 'F'
	ternaryTrueMask    ternary = ^ternaryFalse // Any 8-bit value besides 0 and 'F' is considered true
)

// ternary holds a value in one of three states: true, false or unspecified.
type ternary byte

// True returns true if t specifies a true value.
func (t ternary) True() bool {
	return t&ternaryTrueMask != 0
}

// False returns false if t specifies a true value.
func (t ternary) False() bool {
	return t == ternaryFalse
}

// Specified returns true if t specifies a value.
func (t ternary) Specified() bool {
	return t != ternaryUnspecified
}
