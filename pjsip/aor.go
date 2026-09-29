package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astmerge"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip/#aor

// AOR is a pjsip address of record section.
//
// TODO: Add the rest of the possible fields.
type AOR struct {
	Name              string           `astconf:"-"`
	Templates         []string         `astconf:"-"`
	Contacts          []string         `astconf:"contact,omitempty"`
	MaxContacts       astval.Int       `astconf:"max_contacts"`
	RemoveExisting    astval.YesNoNone `astconf:"remove_existing"`
	QualifyFrequency  astval.Seconds   `astconf:"qualify_frequency"`
	DefaultExpiration astval.Seconds   `astconf:"default_expiration"`
	MinimumExpiration astval.Seconds   `astconf:"minimum_expiration"`
	MaximumExpiration astval.Seconds   `astconf:"maximum_expiration"`
	Mailboxes         []string         `astconf:"mailboxes,commaseparated,omitempty"`
}

// SectionName returns the name of the aor section.
func (aor *AOR) SectionName() string {
	return aor.Name
}

// SectionTemplates returns the section templates used by the aor.
func (aor *AOR) SectionTemplates() []string {
	return aor.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (aor *AOR) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "aor")
}

// OverlayAORs returns the overlayed configuration of all the given
// aors, in order of priority from least to greatest.
func OverlayAORs(aors ...AOR) (overlayed AOR) {
	for i := range aors {
		overlayAORScalars(&aors[i], &overlayed)
		overlayAORVectors(&aors[i], &overlayed)
	}
	return
}

// MergeAORs returns the merged configuration of all the given aors,
// in order of priority from least to greatest.
func MergeAORs(aors ...AOR) (merged AOR) {
	for i := range aors {
		overlayAORScalars(&aors[i], &merged)
		mergeAORVectors(&aors[i], &merged)
	}
	return
}

// overlayAORScalars overlays all scalar values in from with values from to.
func overlayAORScalars(from, to *AOR) {
	astoverlay.String(&from.Name, &to.Name)
	astoverlay.AstInt(&from.MaxContacts, &to.MaxContacts)
	astoverlay.AstYesNoNone(&from.RemoveExisting, &to.RemoveExisting)
	astoverlay.AstSeconds(&from.QualifyFrequency, &to.QualifyFrequency)
	astoverlay.AstSeconds(&from.DefaultExpiration, &to.DefaultExpiration)
	astoverlay.AstSeconds(&from.MinimumExpiration, &to.MinimumExpiration)
	astoverlay.AstSeconds(&from.MaximumExpiration, &to.MaximumExpiration)
}

// overlayAORVectors overlays all vector values in from with values from to.
func overlayAORVectors(from, to *AOR) {
	astoverlay.StringSlice(&from.Templates, &to.Templates)
	astoverlay.StringSlice(&from.Contacts, &to.Contacts)
	astoverlay.StringSlice(&from.Mailboxes, &to.Mailboxes)
}

// mergeAORVectors merges all vector values in from with values from to.
func mergeAORVectors(from, to *AOR) {
	astmerge.StringSlice(&from.Templates, &to.Templates)
	astmerge.StringSlice(&from.Contacts, &to.Contacts)
	astmerge.StringSlice(&from.Mailboxes, &to.Mailboxes)
}
