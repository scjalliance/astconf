package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip/#identify

// Identify is a pjsip identify section. It matches inbound requests to an
// endpoint by source address or by a header, which is how a trunk with no
// registration is recognized.
//
// Identify sections are written once per peer and are not generated per
// device, so there are no overlay or merge functions for them.
//
// Only the fields that multi-server rendering needs are modeled.
type Identify struct {
	Name        string           `astconf:"-"`
	Templates   []string         `astconf:"-"`
	Endpoint    string           `astconf:"endpoint,omitempty"`
	Match       []string         `astconf:"match,omitempty"` // IP addresses, networks or hostnames
	SRVLookups  astval.YesNoNone `astconf:"srv_lookups"`
	MatchHeader string           `astconf:"match_header,omitempty"` // "Header-Name: value"
}

// SectionName returns the name of the identify section.
func (identify *Identify) SectionName() string {
	return identify.Name
}

// SectionTemplates returns the section templates used by the identify.
func (identify *Identify) SectionTemplates() []string {
	return identify.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (identify *Identify) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "identify")
}
