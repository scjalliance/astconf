package dialplan

// Section is a dialplan section with a common context.
type Section struct {
	Context    string `astconf:"-"`
	Extensions []Extension
}

// SectionName returns the name of the voicemail context section.
func (section *Section) SectionName() string {
	return section.Context
}
