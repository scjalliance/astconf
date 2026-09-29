package pjsip

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
)

// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Module_Configuration/res_pjsip/#transport

// Transport is a pjsip transport section.
//
// Transports are written once per server and are not generated per
// device, so there are no overlay or merge functions for them.
//
// TODO: Add the rest of the possible fields.
type Transport struct {
	Name                     string           `astconf:"-"`
	Templates                []string         `astconf:"-"`
	Protocol                 string           `astconf:"protocol,omitempty"` // "udp", "tcp", "tls", "ws", "wss"
	Bind                     string           `astconf:"bind,omitempty"`
	ExternalMediaAddress     string           `astconf:"external_media_address,omitempty"`
	ExternalSignalingAddress string           `astconf:"external_signaling_address,omitempty"`
	ExternalSignalingPort    int              `astconf:"external_signaling_port,omitempty"`
	LocalNet                 []string         `astconf:"local_net,omitempty"`
	CertFile                 string           `astconf:"cert_file,omitempty"`
	PrivKeyFile              string           `astconf:"priv_key_file,omitempty"`
	CAListFile               string           `astconf:"ca_list_file,omitempty"`
	Method                   string           `astconf:"method,omitempty"` // "default", "unspecified", "tlsv1", "tlsv1_1", "tlsv1_2", "tlsv1_3", "sslv2", "sslv3", "sslv23"
	AllowReload              astval.YesNoNone `astconf:"allow_reload"`
}

// SectionName returns the name of the transport section.
func (transport *Transport) SectionName() string {
	return transport.Name
}

// SectionTemplates returns the section templates used by the transport.
func (transport *Transport) SectionTemplates() []string {
	return transport.Templates
}

// MarshalAsteriskPreamble marshals the type.
func (transport *Transport) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "transport")
}
