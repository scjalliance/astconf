package pjsip

import (
	"github.com/scjalliance/astconf"
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
		auth := &auths[i]
		astoverlay.String(&auth.Name, &overlayed.Name)
		astoverlay.StringSlice(&auth.Templates, &overlayed.Templates)
		astoverlay.String(&auth.AuthType, &overlayed.AuthType)
		astoverlay.String(&auth.Username, &overlayed.Username)
		astoverlay.String(&auth.Password, &overlayed.Password)
		astoverlay.String(&auth.MD5Cred, &overlayed.MD5Cred)
		astoverlay.String(&auth.Realm, &overlayed.Realm)
	}
	return
}
