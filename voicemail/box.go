package voicemail

import (
	"fmt"
	"strings"

	"github.com/scjalliance/astconf"
)

// https://github.com/asterisk/asterisk/blob/master/configs/samples/voicemail.conf.sample
// https://www.voip-info.org/asterisk-config-voicemailconf/

// Box holds settings for a voicemail mailbox.
type Box struct {
	// Extension is the number of the mailbox.
	Extension string

	// Name is the name of the mailbox.
	Name string

	// Password is the password used to retrieve messages.
	Password string

	// PasswordIsChangeable enables self-service voicemail password changes.
	PasswordIsChangeable bool

	// EmailAddresses are the set of recipients that will receive email
	// messages.
	EmailAddresses []string

	// PagerEmailAddress is the email address that will receive pager
	// notifications.
	PagerEmailAddress string

	// Timezone assigns a particular timezone message to the mailbox.
	Timezone string

	// Locale declares a particular locale to use for date and time strings.
	Locale string

	// SendToPager causes the voicemail to be sent to the pager address.
	SendToPager bool

	// Format overrides the default audio format that will be attached
	// to an email.
	Format string

	// SayCallerID causes the caller to be identified before playing back a
	// message.
	SayCallerID bool

	// SkipEnvelope prevents the envelope (date/time) from being presented
	//  to the user before playing back a message.
	SkipEnvelope bool

	// EmailOnly causes messages to be deleted from the server after they
	// have been sent via email.
	EmailOnly bool
}

// Options returns an options string representing the mailbox options.
func (box *Box) Options() (options []string) {
	if box.Timezone != "" {
		options = append(options, fmt.Sprintf("tz=%s", sanitizeOption(box.Timezone)))
	}
	if box.Locale != "" {
		options = append(options, fmt.Sprintf("locale=%s", sanitizeOption(box.Locale)))
	}
	if !box.SendToPager {
		options = append(options, "attach=yes")
	}
	if box.Format != "" {
		options = append(options, fmt.Sprintf("attachfmt=%s", sanitizeOption(box.Format)))
	}
	if box.SayCallerID {
		options = append(options, "saycid=yes")
	}
	if box.SkipEnvelope {
		options = append(options, "envelope=no")
	}
	if box.EmailOnly {
		options = append(options, "delete=yes")
	} else {
		options = append(options, "delete=no")
	}
	return
}

// MarshalAsterisk marshals the mailbox to an asterisk encoder.
func (box *Box) MarshalAsterisk(e *astconf.Encoder) error {
	// Mailbox values are a comma-separated set of five components
	components := make([]string, 5)

	// Construct the component values (password,name,email,pager_email,options)
	if box.PasswordIsChangeable {
		components[0] = box.Password
	} else {
		components[0] = "-" + box.Password
	}
	components[1] = box.Name
	components[2] = strings.Join(box.EmailAddresses, "|")
	components[3] = box.PagerEmailAddress
	components[4] = strings.Join(box.Options(), "|")

	// Determine the number of components we'll need to include
	count := len(components)
	for count > 0 {
		if components[count-1] != "" {
			break
		}
		count--
	}

	// Construct the value
	value := strings.Join(components[:count], ",")

	// Print the mailbox as an object
	return e.Printer().Object(box.Extension, value)
}
