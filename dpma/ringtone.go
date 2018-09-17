package dpma

import (
	"strings"

	"github.com/scjalliance/astconf"
)

// Ringtone is a DPMA ringtone definition.
type Ringtone struct {
	Alias    string `astconf:"alias"`
	Filename string `astconf:"filename"`
}

// SectionName returns the name of the ringtone section.
func (r *Ringtone) SectionName() string {
	return strings.ToLower(strings.Replace(string(r.Alias), " ", "-", -1))
}

// MarshalAsteriskPreamble marshals the type.
func (r *Ringtone) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "ringtone")
}
