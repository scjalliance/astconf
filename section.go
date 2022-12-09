package astconf

// SectionNamer provides the section name for a type.
type SectionNamer interface {
	SectionName() string
}

// SectionName is a simple implementation of SectionNamer. Its value won't
// be written as a setting.
type SectionName string

// SectionName returns s as the name for the section.
func (s SectionName) SectionName() string {
	return string(s)
}

// MarshalAsterisk prevents the section name from being marshaled.
func (s SectionName) MarshalAsterisk(e *Encoder) error {
	return nil
}

// SectionTemplater provides section templates for a type. It must be provided
// by the same type that provides the name for a section.
type SectionTemplater interface {
	SectionTemplates() []string
}

/*
// Section identifies a well-known section properties when embedded in a type.
// The section name to be used is derived from its tag by the encoder.
type Section struct{}
*/
