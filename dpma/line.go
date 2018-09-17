package dpma

import (
	"github.com/scjalliance/astconf"
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
		mergeSectionName(&line.Name, &merged.Name)
		mergeString(&line.Extension, &merged.Extension)
		mergeString(&line.DigitMap, &merged.DigitMap)
		mergeString(&line.Label, &merged.Label)
		mergeString(&line.Mailbox, &merged.Mailbox)
		mergeString(&line.VoicemailURI, &merged.VoicemailURI)
		mergeString(&line.OutboundProxyAddress, &merged.OutboundProxyAddress)
		mergeString(&line.OutboundProxyPort, &merged.OutboundProxyPort)
		mergeString(&line.Transport, &merged.Transport)
		mergeString(&line.MediaEncryption, &merged.MediaEncryption)
		mergeAstSeconds(&line.RegistrationTimeout, &merged.RegistrationTimeout)
		mergeAstSeconds(&line.RegistrationRetryInterval, &merged.RegistrationRetryInterval)
		mergeAstInt(&line.RegistrationMaxRetries, &merged.RegistrationMaxRetries)
		mergeString(&line.Secret, &merged.Secret)
		mergeString(&line.Context, &merged.Context)
		mergeString(&line.CallerID, &merged.CallerID)
		mergeString(&line.SubscribeContext, &merged.SubscribeContext)
		mergeString(&line.PlarNumber, &merged.PlarNumber)
	}
	return
}
