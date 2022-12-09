package sip

import (
	"github.com/scjalliance/astconf/astmerge"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

// https://github.com/asterisk/asterisk/blob/master/configs/samples/sip.conf.sample
// https://www.voip-info.org/asterisk-config-sipconf/

// Entity is a sip entity.
//
// TODO: Add the rest of the possible fields.
type Entity struct {
	Username           string           `astconf:"-"`
	Templates          []string         `astconf:"-"`
	Type               Type             `astconf:"type"`
	AccountCode        string           `astconf:"accountcode,omitempty"`
	Disallow           []string         `astconf:"disallow,omitempty"`
	Allow              []string         `astconf:"allow,omitempty"`
	AllowGuest         astval.YesNoNone `astconf:"allowguest,omitempty"`
	AccountingFlags    string           `astconf:"amaflags,omitempty"`
	AsteriskDB         string           `astconf:"astdb,omitempty"`
	Auth               string           `astconf:"auth,omitempty"`
	BusyLevel          int              `astconf:"busylevel,omitempty"`
	CallLimit          int              `astconf:"call-limit,omitempty"`
	CallGroup          string           `astconf:"callgroup,omitempty"` // FIXME: Use a slice of something?
	CallerID           string           `astconf:"callerid,omitempty"`
	CallerPresentation string           `astconf:"callingpres,omitempty"`
	DirectMedia        string           `astconf:"directmedia,omitempty"`
	DirectMediaPermit  []string         `astconf:"directmediapermit,omitempty"`
	DirectMediaDeny    []string         `astconf:"directmediadeny,omitempty"`
	Context            string           `astconf:"context,omitempty"`
	Host               string           `astconf:"host,omitempty"`
	Transport          []string         `astconf:"transport,commaseparated,omitempty"`
	Mailbox            string           `astconf:"mailbox,omitempty"`
	Secret             string           `astconf:"secret,omitempty"`
	Variables          []astval.Var     `astconf:"setvar,omitempty"`
}

// SectionName returns the name of the section that the entity belongs to.
func (e *Entity) SectionName() string {
	return e.Username
}

// SectionTemplates returns the section templates used by an entity.
func (e *Entity) SectionTemplates() []string {
	return e.Templates
}

// OverlayEntities returns the overlayed configuration of all the given
// entities, in order of priority from least to greatest.
func OverlayEntities(entities ...Entity) (overlayed Entity) {
	for i := range entities {
		overlayEntityScalars(&entities[i], &overlayed)
		overlayEntityVectors(&entities[i], &overlayed)
	}
	return
}

// MergeEntities returns the merged configuration of all the given entities,
// in order of priority from least to greatest.
func MergeEntities(entities ...Entity) (merged Entity) {
	for i := range entities {
		overlayEntityScalars(&entities[i], &merged)
		mergeEntityVectors(&entities[i], &merged)
	}
	return
}

// overlayEntityScalars overlays all scalar values in from with values from to.
func overlayEntityScalars(from, to *Entity) {
	overlayType(&from.Type, &to.Type)
	astoverlay.String(&from.AccountCode, &to.AccountCode)
	astoverlay.AstYesNoNone(&from.AllowGuest, &to.AllowGuest)
	astoverlay.String(&from.AccountingFlags, &to.AccountingFlags)
	astoverlay.String(&from.AsteriskDB, &to.AsteriskDB)
	astoverlay.String(&from.Auth, &to.Auth)
	astoverlay.Int(&from.BusyLevel, &to.BusyLevel)
	astoverlay.Int(&from.CallLimit, &to.CallLimit)
	astoverlay.String(&from.CallGroup, &to.CallGroup)
	astoverlay.String(&from.CallerID, &to.CallerID)
	astoverlay.String(&from.CallerPresentation, &to.CallerPresentation)
	astoverlay.String(&from.DirectMedia, &to.DirectMedia)
	astoverlay.StringSlice(&from.DirectMediaPermit, &to.DirectMediaPermit)
	astoverlay.StringSlice(&from.DirectMediaDeny, &to.DirectMediaDeny)
	astoverlay.String(&from.Context, &to.Context)
	astoverlay.String(&from.Host, &to.Host)
	astoverlay.StringSlice(&from.Transport, &to.Transport)
	astoverlay.String(&from.Mailbox, &to.Mailbox)
	astoverlay.String(&from.Secret, &to.Secret)
	astoverlay.String(&from.Username, &to.Username)
}

// overlayEntityVectors overlays all vector values in from with values from to.
func overlayEntityVectors(from, to *Entity) {
	astoverlay.StringSlice(&from.Templates, &to.Templates)
	astoverlay.StringSlice(&from.Disallow, &to.Disallow)
	astoverlay.StringSlice(&from.Allow, &to.Allow)
	astoverlay.AstVarSlice(&from.Variables, &to.Variables)
}

// mergeEntityVectors merges all vector values in from with values from to.
func mergeEntityVectors(from, to *Entity) {
	astmerge.StringSlice(&from.Templates, &to.Templates)
	astmerge.StringSlice(&from.Disallow, &to.Disallow)
	astmerge.StringSlice(&from.Allow, &to.Allow)
	astmerge.AstVarSlice(&from.Variables, &to.Variables)
}
