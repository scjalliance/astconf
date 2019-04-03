package dpma

import "github.com/scjalliance/astconf"

// Alert is a DPMA alert definition.
type Alert struct {
	Name       string `astconf:"-"`
	InfoHeader string `astconf:"alert_info"`
	RingType   string `astconf:"ring_type"`
	Ringtone   string `astconf:"ringtone"`
}

// SectionName returns the name of the ringtone section.
func (a *Alert) SectionName() string {
	return a.Name
}

// MarshalAsteriskPreamble marshals the type.
func (*Alert) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "alert")
}
