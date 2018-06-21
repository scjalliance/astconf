package dpma

import "time"

// Line is a DPMA line definition.
type Line struct {
	Extension                 string
	Label                     string
	Mailbox                   string
	VoicemailURI              string
	OutboundProxyAddress      string
	Transport                 string
	MediaEncryption           string
	RegistrationTimeout       time.Duration // TODO: Convert to seconds
	RegistrationRetryInterval time.Duration // TODO: Convert to seconds
	RegistrationMaxRetries    int
	Secret                    string
	Context                   string
	CallerID                  string
	SubscribeContext          string
	PlarNumber                string
}
