package dpma

import (
	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

// Line is a DPMA line definition.
type Line struct {
	Name                      astconf.SectionName
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

// MarshalAsteriskPreamble marshals the type.
func (line *Line) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "line")
}

// MergeLines returns the merged configuration of all the given lines,
// in order of priority from least to greatest.
func MergeLines(lines ...Line) (merged Line) {
	for i := range lines {
		line := &lines[i]
		astoverlay.SectionName(&line.Name, &merged.Name)
		astoverlay.String(&line.Extension, &merged.Extension)
		astoverlay.String(&line.DigitMap, &merged.DigitMap)
		astoverlay.String(&line.Label, &merged.Label)
		astoverlay.String(&line.Mailbox, &merged.Mailbox)
		astoverlay.String(&line.VoicemailURI, &merged.VoicemailURI)
		astoverlay.String(&line.OutboundProxyAddress, &merged.OutboundProxyAddress)
		astoverlay.String(&line.OutboundProxyPort, &merged.OutboundProxyPort)
		astoverlay.String(&line.Transport, &merged.Transport)
		astoverlay.String(&line.MediaEncryption, &merged.MediaEncryption)
		astoverlay.AstSeconds(&line.RegistrationTimeout, &merged.RegistrationTimeout)
		astoverlay.AstSeconds(&line.RegistrationRetryInterval, &merged.RegistrationRetryInterval)
		astoverlay.AstInt(&line.RegistrationMaxRetries, &merged.RegistrationMaxRetries)
		astoverlay.String(&line.Secret, &merged.Secret)
		astoverlay.String(&line.Context, &merged.Context)
		astoverlay.String(&line.CallerID, &merged.CallerID)
		astoverlay.String(&line.SubscribeContext, &merged.SubscribeContext)
		astoverlay.String(&line.PlarNumber, &merged.PlarNumber)
	}
	return
}
