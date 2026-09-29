package dialplan

import "fmt"

// SIPAddHeader returns a dialplan application that adds a SIP header to
// outbound calls.
func SIPAddHeader(header, content string) SIPAddHeaderApp {
	return SIPAddHeaderApp{
		Header:  header,
		Content: content,
	}
}

// SIPAddHeaderApp is a dialplan application that adds a SIP header to
// outbound calls.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_SIPAddHeader
type SIPAddHeaderApp struct {
	Header  string
	Content string
}

// App returns the assembled application call.
func (app SIPAddHeaderApp) App() AppCall {
	return AppCall{
		Name: "SIPAddHeader",
		Args: []string{fmt.Sprintf("%s:%s", app.Header, app.Content)},
	}
}

// Validate returns an error if the header or its content contains invalid
// characters.
func (app SIPAddHeaderApp) Validate() error {
	if err := errorIfAny("header name", app.Header, invalidArgChars); err != nil {
		return err
	}
	return errorIfAny("header content", app.Content, invalidArgChars)
}
