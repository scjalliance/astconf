package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astmerge"
	"github.com/scjalliance/astconf/astoverlay"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip/#auth

// Auth is a pjsip auth section.
//
// TODO: Add the rest of the possible fields.
type Auth struct {
	Name      string   `astconf:"-"`
	Templates []string `astconf:"-"`
	AuthType  string   `astconf:"auth_type,omitempty"` // "userpass", "md5"
	Username  string   `astconf:"username,omitempty"`
	Password  string   `astconf:"password,omitempty"`
	MD5Cred   string   `astconf:"md5_cred,omitempty"`
	Realm     string   `astconf:"realm,omitempty"`
}

// SectionName returns the name of the auth section.
func (auth *Auth) SectionName() string {
	return auth.Name
}

// SectionTemplates returns the section templates used by the auth.
func (auth *Auth) SectionTemplates() []string {
	return auth.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (auth *Auth) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "auth")
}

// OverlayAuths returns the overlayed configuration of all the given
// auths, in order of priority from least to greatest.
func OverlayAuths(auths ...Auth) (overlayed Auth) {
	for i := range auths {
		overlayAuthScalars(&auths[i], &overlayed)
		overlayAuthVectors(&auths[i], &overlayed)
	}
	return
}

// MergeAuths returns the merged configuration of all the given auths,
// in order of priority from least to greatest.
func MergeAuths(auths ...Auth) (merged Auth) {
	for i := range auths {
		overlayAuthScalars(&auths[i], &merged)
		mergeAuthVectors(&auths[i], &merged)
	}
	return
}

// overlayAuthScalars overlays all scalar values in from with values from to.
func overlayAuthScalars(from, to *Auth) {
	astoverlay.String(&from.Name, &to.Name)
	astoverlay.String(&from.AuthType, &to.AuthType)
	astoverlay.String(&from.Username, &to.Username)
	astoverlay.String(&from.Password, &to.Password)
	astoverlay.String(&from.MD5Cred, &to.MD5Cred)
	astoverlay.String(&from.Realm, &to.Realm)
}

// overlayAuthVectors overlays all vector values in from with values from to.
func overlayAuthVectors(from, to *Auth) {
	astoverlay.StringSlice(&from.Templates, &to.Templates)
}

// mergeAuthVectors merges all vector values in from with values from to.
func mergeAuthVectors(from, to *Auth) {
	astmerge.StringSlice(&from.Templates, &to.Templates)
}
