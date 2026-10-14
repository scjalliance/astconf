package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_publish_asterisk/

// AsteriskPublication is a pjsip asterisk-publication section. It ties
// device and mailbox state sharing with one other Asterisk server together.
//
// DeviceState and MailboxState permit incoming state from that server, and
// the two filters limit which incoming devices and mailboxes are accepted.
// What is sent is limited by the filters on the outbound-publish sections
// that DeviceStatePublish and MailboxStatePublish name.
//
// Asterisk-publication sections are written once per peer and are not
// generated per device, so there are no overlay or merge functions for
// them.
type AsteriskPublication struct {
	Name                string           `astconf:"-"`
	Templates           []string         `astconf:"-"`
	DeviceStatePublish  string           `astconf:"devicestate_publish,omitempty"`  // Name of an outbound-publish section
	MailboxStatePublish string           `astconf:"mailboxstate_publish,omitempty"` // Name of an outbound-publish section
	DeviceState         astval.YesNoNone `astconf:"device_state"`
	DeviceStateFilter   string           `astconf:"device_state_filter,omitempty"` // Regular expression
	MailboxState        astval.YesNoNone `astconf:"mailbox_state"`
	MailboxStateFilter  string           `astconf:"mailbox_state_filter,omitempty"` // Regular expression
}

// SectionName returns the name of the asterisk-publication section.
func (publication *AsteriskPublication) SectionName() string {
	return publication.Name
}

// SectionTemplates returns the section templates used by the
// asterisk-publication.
func (publication *AsteriskPublication) SectionTemplates() []string {
	return publication.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (publication *AsteriskPublication) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "asterisk-publication")
}
