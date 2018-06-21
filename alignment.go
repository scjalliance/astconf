package astconf

// Alignment indicates whether names should be aligned
// in some fashion.
type Alignment int

// Name alignments.
const (
	AlignmentNone Alignment = iota
	AlignmentLeft
	AlignmentRight
)

// Unaligned configures an encoder not to align names.
func Unaligned(e *Encoder) {
	e.alignment = AlignmentNone
}

// AlignLeft configures an encoder to align names to the left.
func AlignLeft(e *Encoder) {
	e.alignment = AlignmentLeft
}

// AlignRight configures an encoder to align names to the right.
func AlignRight(e *Encoder) {
	e.alignment = AlignmentRight
}
