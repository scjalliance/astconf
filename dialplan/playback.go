package dialplan

import (
	"strings"
)

// Playback returns a dialplan application that plays sound files.
func Playback(files ...string) PlaybackApp {
	return PlaybackApp{
		Files: files,
	}
}

// PlaybackApp is a dialplan application that plays sound files.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Playback
type PlaybackApp struct {
	Files    []string
	Skip     bool
	NoAnswer bool
}

// App returns the assembled application call.
func (app PlaybackApp) App() AppCall {
	call := AppCall{
		Name: "Playback",
		Args: []string{strings.Join(app.Files, "&")},
	}
	if app.Skip {
		call.Args = append(call.Args, "skip")
	}
	if app.NoAnswer {
		call.Args = append(call.Args, "noanswer")
	}
	return call
}
