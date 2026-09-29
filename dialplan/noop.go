package dialplan

// Noop returns a dialplan application that does nothing.
func Noop(text string) NoopApp {
	return NoopApp{
		Text: text,
	}
}

// NoopApp is a dialplan application that does nothing.
type NoopApp struct {
	Text string
}

// App returns the assembled application call.
func (noop NoopApp) App() AppCall {
	if noop.Text == "" {
		return AppCall{Name: "Noop"}
	}
	return AppCall{
		Name: "Noop",
		Args: []string{noop.Text},
	}
}

// Validate returns an error if the text contains invalid characters.
func (noop NoopApp) Validate() error {
	return errorIfAny("noop text", noop.Text, invalidArgChars)
}
