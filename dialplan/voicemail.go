package dialplan

// VoiceMail returns a dialplan application that leaves a message in a
// mailbox. Options are passed through, such as "u" for the unavailable
// greeting.
func VoiceMail(mailbox, context, options string) VoiceMailApp {
	return VoiceMailApp{
		Mailbox: mailbox,
		Context: context,
		Options: options,
	}
}

// VoiceMailApp is a dialplan application that leaves a message in a
// mailbox.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/VoiceMail/
type VoiceMailApp struct {
	Mailbox string
	Context string
	Options string
}

// App returns the assembled application call.
func (app VoiceMailApp) App() AppCall {
	box := app.Mailbox
	if app.Context != "" {
		box += "@" + app.Context
	}
	call := AppCall{Name: "VoiceMail", Args: []string{box}}
	if app.Options != "" {
		call.Args = append(call.Args, app.Options)
	}
	return call
}

// Validate returns an error if the mailbox, context or options contain
// invalid characters.
func (app VoiceMailApp) Validate() error {
	if err := errorIfAny("voicemail mailbox", app.Mailbox, invalidResourceChars); err != nil {
		return err
	}
	if err := errorIfAny("voicemail context", app.Context, invalidResourceChars); err != nil {
		return err
	}
	return errorIfAny("voicemail options", app.Options, invalidArgChars)
}
