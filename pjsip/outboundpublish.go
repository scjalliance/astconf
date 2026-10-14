package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_outbound_publish/#outbound-publish
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_publish_asterisk/

// OutboundPublish is a pjsip outbound-publish section. It sends PUBLISH
// requests for one event to one server, such as device or mailbox state
// for another Asterisk server.
//
// The user part of ServerURI names the inbound-publication section that
// receives the requests on the other server.
//
// DeviceStateFilter and MailboxStateFilter are the extended options that
// res_pjsip_publish_asterisk reads from an outbound-publish. They limit
// which local devices or mailboxes are published, and are written with a
// leading "@" as Asterisk requires for extended options.
//
// Outbound-publish sections are written once per peer and are not
// generated per device, so there are no overlay or merge functions for
// them.
//
// Only the fields that multi-server rendering needs are modeled.
type OutboundPublish struct {
	Name               string         `astconf:"-"`
	Templates          []string       `astconf:"-"`
	ServerURI          string         `astconf:"server_uri,omitempty"`
	FromURI            string         `astconf:"from_uri,omitempty"`
	ToURI              string         `astconf:"to_uri,omitempty"`
	Event              string         `astconf:"event,omitempty"` // "asterisk-devicestate", "asterisk-mwi"
	OutboundAuth       string         `astconf:"outbound_auth,omitempty"`
	Expiration         astval.Seconds `astconf:"expiration"`
	Transport          string         `astconf:"transport,omitempty"`
	OutboundProxy      string         `astconf:"outbound_proxy,omitempty"`
	MaxAuthAttempts    astval.Int     `astconf:"max_auth_attempts"`
	DeviceStateFilter  string         `astconf:"@device_state_filter,omitempty"`  // Regular expression
	MailboxStateFilter string         `astconf:"@mailbox_state_filter,omitempty"` // Regular expression
}

// SectionName returns the name of the outbound-publish section.
func (publish *OutboundPublish) SectionName() string {
	return publish.Name
}

// SectionTemplates returns the section templates used by the
// outbound-publish.
func (publish *OutboundPublish) SectionTemplates() []string {
	return publish.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (publish *OutboundPublish) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "outbound-publish")
}
