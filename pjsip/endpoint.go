package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astmerge"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

// https://github.com/asterisk/asterisk/blob/master/configs/samples/pjsip.conf.sample
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip/#endpoint

// Endpoint is a pjsip endpoint section.
//
// TODO: Add the rest of the possible fields.
type Endpoint struct {
	Name                       string           `astconf:"-"`
	Templates                  []string         `astconf:"-"`
	Transport                  string           `astconf:"transport,omitempty"`
	Context                    string           `astconf:"context,omitempty"`
	Disallow                   []string         `astconf:"disallow,omitempty"`
	Allow                      []string         `astconf:"allow,omitempty"`
	Auth                       string           `astconf:"auth,omitempty"`
	OutboundAuth               string           `astconf:"outbound_auth,omitempty"`
	AORs                       []string         `astconf:"aors,commaseparated,omitempty"`
	IdentifyBy                 string           `astconf:"identify_by,omitempty"` // "username", "auth_username", "ip", "header"
	CallerID                   string           `astconf:"callerid,omitempty"`
	FromUser                   string           `astconf:"from_user,omitempty"`
	FromDomain                 string           `astconf:"from_domain,omitempty"`
	Mailboxes                  []string         `astconf:"mailboxes,commaseparated,omitempty"`
	DirectMedia                astval.YesNoNone `astconf:"direct_media"`
	DirectMediaGlareMitigation string           `astconf:"direct_media_glare_mitigation,omitempty"` // "none", "outgoing", "incoming"
	RTPSymmetric               astval.YesNoNone `astconf:"rtp_symmetric"`
	ForceRport                 astval.YesNoNone `astconf:"force_rport"`
	RewriteContact             astval.YesNoNone `astconf:"rewrite_contact"`
	DTMFMode                   string           `astconf:"dtmf_mode,omitempty"` // "rfc4733", "inband", "info", "auto"
	DeviceStateBusyAt          astval.Int       `astconf:"device_state_busy_at"`
	Variables                  []astval.Var     `astconf:"set_var,omitempty"`
}

// SectionName returns the name of the endpoint section.
func (endpoint *Endpoint) SectionName() string {
	return endpoint.Name
}

// SectionTemplates returns the section templates used by the endpoint.
func (endpoint *Endpoint) SectionTemplates() []string {
	return endpoint.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (endpoint *Endpoint) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "endpoint")
}

// OverlayEndpoints returns the overlayed configuration of all the given
// endpoints, in order of priority from least to greatest.
func OverlayEndpoints(endpoints ...Endpoint) (overlayed Endpoint) {
	for i := range endpoints {
		overlayEndpointScalars(&endpoints[i], &overlayed)
		overlayEndpointVectors(&endpoints[i], &overlayed)
	}
	return
}

// MergeEndpoints returns the merged configuration of all the given
// endpoints, in order of priority from least to greatest.
func MergeEndpoints(endpoints ...Endpoint) (merged Endpoint) {
	for i := range endpoints {
		overlayEndpointScalars(&endpoints[i], &merged)
		mergeEndpointVectors(&endpoints[i], &merged)
	}
	return
}

// overlayEndpointScalars overlays all scalar values in from with values from to.
func overlayEndpointScalars(from, to *Endpoint) {
	astoverlay.String(&from.Name, &to.Name)
	astoverlay.String(&from.Transport, &to.Transport)
	astoverlay.String(&from.Context, &to.Context)
	astoverlay.String(&from.Auth, &to.Auth)
	astoverlay.String(&from.OutboundAuth, &to.OutboundAuth)
	astoverlay.String(&from.IdentifyBy, &to.IdentifyBy)
	astoverlay.String(&from.CallerID, &to.CallerID)
	astoverlay.String(&from.FromUser, &to.FromUser)
	astoverlay.String(&from.FromDomain, &to.FromDomain)
	astoverlay.AstYesNoNone(&from.DirectMedia, &to.DirectMedia)
	astoverlay.String(&from.DirectMediaGlareMitigation, &to.DirectMediaGlareMitigation)
	astoverlay.AstYesNoNone(&from.RTPSymmetric, &to.RTPSymmetric)
	astoverlay.AstYesNoNone(&from.ForceRport, &to.ForceRport)
	astoverlay.AstYesNoNone(&from.RewriteContact, &to.RewriteContact)
	astoverlay.String(&from.DTMFMode, &to.DTMFMode)
	astoverlay.AstInt(&from.DeviceStateBusyAt, &to.DeviceStateBusyAt)
}

// overlayEndpointVectors overlays all vector values in from with values from to.
func overlayEndpointVectors(from, to *Endpoint) {
	astoverlay.StringSlice(&from.Templates, &to.Templates)
	astoverlay.StringSlice(&from.Disallow, &to.Disallow)
	astoverlay.StringSlice(&from.Allow, &to.Allow)
	astoverlay.StringSlice(&from.AORs, &to.AORs)
	astoverlay.StringSlice(&from.Mailboxes, &to.Mailboxes)
	astoverlay.AstVarSlice(&from.Variables, &to.Variables)
}

// mergeEndpointVectors merges all vector values in from with values from to.
func mergeEndpointVectors(from, to *Endpoint) {
	astmerge.StringSlice(&from.Templates, &to.Templates)
	astmerge.StringSlice(&from.Disallow, &to.Disallow)
	astmerge.StringSlice(&from.Allow, &to.Allow)
	astmerge.StringSlice(&from.AORs, &to.AORs)
	astmerge.StringSlice(&from.Mailboxes, &to.Mailboxes)
	astmerge.AstVarSlice(&from.Variables, &to.Variables)
}
