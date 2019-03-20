package voicemail

// Section is a mailbox section with a common context.
type Section struct {
	Context   string `astconf:"-"`
	Mailboxes []Box
}

// SectionName returns the name of the voicemail context section.
func (section *Section) SectionName() string {
	return section.Context
}
