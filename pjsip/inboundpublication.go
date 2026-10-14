package pjsip

import (
	"github.com/scjalliance/astconf"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip_pubsub/#inbound-publication

// InboundPublication is a pjsip inbound-publication section. It accepts
// PUBLISH requests for one resource and hands each event to a handler.
//
// The section name must match the user part of the request URI that the
// publishing server sends to, which is the user part of its outbound-publish
// server_uri.
//
// Inbound-publication sections are written once per peer and are not
// generated per device, so there are no overlay or merge functions for
// them.
type InboundPublication struct {
	Name      string   `astconf:"-"`
	Templates []string `astconf:"-"`
	Endpoint  string   `astconf:"endpoint,omitempty"` // The only endpoint allowed to publish
	Events    PublicationEvents
}

// SectionName returns the name of the inbound-publication section.
func (publication *InboundPublication) SectionName() string {
	return publication.Name
}

// SectionTemplates returns the section templates used by the
// inbound-publication.
func (publication *InboundPublication) SectionTemplates() []string {
	return publication.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (publication *InboundPublication) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "inbound-publication")
}

// PublicationEvent maps one event of an inbound publication to the
// configuration that handles it. For the asterisk-devicestate and
// asterisk-mwi events, Handler is the name of an asterisk-publication
// section.
type PublicationEvent struct {
	Event   string // "asterisk-devicestate", "asterisk-mwi"
	Handler string
}

// PublicationEvents is an ordered list of publication events. Each one is
// written as its own "event_<event> = <handler>" setting, because Asterisk
// reads every setting whose name starts with "event_" as an event mapping.
type PublicationEvents []PublicationEvent

// MarshalAsterisk writes one setting for each event, in order.
func (events PublicationEvents) MarshalAsterisk(e *astconf.Encoder) error {
	for _, event := range events {
		if err := e.Printer().Setting("event_"+event.Event, event.Handler); err != nil {
			return err
		}
	}
	return nil
}
