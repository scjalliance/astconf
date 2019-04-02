package dpma

import (
	"github.com/scjalliance/astconf"
)

// Ringtone is a DPMA ringtone definition.
type Ringtone struct {
	Name     string `astconf:"-"`
	Alias    string `astconf:"alias"`
	Filename string `astconf:"filename"`
}

// SectionName returns the name of the ringtone section.
func (r *Ringtone) SectionName() string {
	return r.Name
}

// MarshalAsteriskPreamble marshals the type.
func (*Ringtone) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "ringtone")
}
