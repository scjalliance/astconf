package dpma

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

// Line is a DPMA line definition.
type Line struct {
	Name                      string         `astconf:"-"`
	Extension                 string         `astconf:"exten,omitempty"`
	DigitMap                  string         `astconf:"digit_map,omitempty"`
	Label                     string         `astconf:"line_label,omitempty"`
	Mailbox                   string         `astconf:"mailbox,omitempty"`
	VoicemailURI              string         `astconf:"voicemail_uri,omitempty"`
	OutboundProxyAddress      string         `astconf:"outboundproxy_address,omitempty"`
	OutboundProxyPort         string         `astconf:"outboundproxy_port,omitempty"`
	Transport                 string         `astconf:"transport,omitempty"`
	MediaEncryption           string         `astconf:"media_encryption,omitempty"`
	RegistrationTimeout       astval.Seconds `astconf:"reregistration_timeout"`
	RegistrationRetryInterval astval.Seconds `astconf:"registration_retry_interval"`
	RegistrationMaxRetries    astval.Int     `astconf:"registration_max_retries"`
	Secret                    string         `astconf:"secret,omitempty"`
	Context                   string         `astconf:"context,omitempty"`
	CallerID                  string         `astconf:"callerid,omitempty"`
	SubscribeContext          string         `astconf:"subscribecontext,omitempty"`
	PlarNumber                string         `astconf:"plar_number,omitempty"`
}

// SectionName returns the name of the line section.
func (line *Line) SectionName() string {
	return line.Name
}

// MarshalAsteriskPreamble marshals the type.
func (line *Line) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "line")
}

// OverlayLines returns the overlayed configuration of all the given lines,
// in order of priority from least to greatest.
func OverlayLines(lines ...Line) (overlayed Line) {
	for i := range lines {
		line := &lines[i]
		astoverlay.String(&line.Name, &overlayed.Name)
		astoverlay.String(&line.Extension, &overlayed.Extension)
		astoverlay.String(&line.DigitMap, &overlayed.DigitMap)
		astoverlay.String(&line.Label, &overlayed.Label)
		astoverlay.String(&line.Mailbox, &overlayed.Mailbox)
		astoverlay.String(&line.VoicemailURI, &overlayed.VoicemailURI)
		astoverlay.String(&line.OutboundProxyAddress, &overlayed.OutboundProxyAddress)
		astoverlay.String(&line.OutboundProxyPort, &overlayed.OutboundProxyPort)
		astoverlay.String(&line.Transport, &overlayed.Transport)
		astoverlay.String(&line.MediaEncryption, &overlayed.MediaEncryption)
		astoverlay.AstSeconds(&line.RegistrationTimeout, &overlayed.RegistrationTimeout)
		astoverlay.AstSeconds(&line.RegistrationRetryInterval, &overlayed.RegistrationRetryInterval)
		astoverlay.AstInt(&line.RegistrationMaxRetries, &overlayed.RegistrationMaxRetries)
		astoverlay.String(&line.Secret, &overlayed.Secret)
		astoverlay.String(&line.Context, &overlayed.Context)
		astoverlay.String(&line.CallerID, &overlayed.CallerID)
		astoverlay.String(&line.SubscribeContext, &overlayed.SubscribeContext)
		astoverlay.String(&line.PlarNumber, &overlayed.PlarNumber)
	}
	return
}
