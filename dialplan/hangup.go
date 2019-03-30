package dialplan

// Hangup returns a dialplan application that hangs up the calling channel.
func Hangup() HangupApp {
	return HangupApp{}
}

// HangupApp is a dialplan application that hangs up the calling channel.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Hangup
type HangupApp struct {
	Cause string
}

// App returns the assembled application call.
func (app HangupApp) App() AppCall {
	if app.Cause != "" {
		return AppCall{Name: "Hangup", Args: []string{app.Cause}}
	}
	return AppCall{Name: "Hangup"}
}
