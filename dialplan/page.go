package dialplan

import (
	"strconv"
	"strings"
)

// http://www.asteriskdocs.org/en/3rd_Edition/asterisk-book-html-chunk/AdditionalConfig_id257169.html

// Page returns a dialplan application that pages a set of devices.
func Page(recipients ...Device) PageApp {
	return PageApp{
		Recipients: recipients,
	}
}

// PageApp is a dialplan application that pages a set of devices.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Page
type PageApp struct {
	Recipients       []Device
	FullDuplex       bool   //  d: Enable two-way audio for all recipients
	Forward          bool   // !i: Forward the page to forwarded recipients
	Quiet            bool   //  q: Don't present a paging beep to recipients
	Record           bool   //  r: Record the page to a file
	InUse            bool   // !s: Include recipeints with a phone state other than NOT_INUSE
	Announcement     string //  A: Play an announcement to all recipients
	AnnounceToCaller bool   //  x: Play an announcement to the caller as well
	Timeout          int    // Number of seconds to attempt a call before giving up
}

// App returns the assembled application call.
func (app PageApp) App() AppCall {
	devices := make([]string, 0, len(app.Recipients))
	for _, device := range app.Recipients {
		devices = append(devices, device.String())
	}
	call := AppCall{
		Name: "Page",
		Args: []string{strings.Join(devices, "&")},
	}
	options := app.Options()
	if len(options) == 0 && app.Timeout <= 0 {
		return call
	}
	call.Args = append(call.Args, options)
	if app.Timeout > 0 {
		call.Args = append(call.Args, strconv.Itoa(app.Timeout))
	}
	return call
}

// Options returns the options string for the page.
func (app PageApp) Options() string {
	var opts string
	if app.FullDuplex {
		opts += "d"
	}
	if !app.Forward {
		opts += "i"
	}
	if app.Quiet {
		opts += "q"
	}
	if app.Record {
		opts += "r"
	}
	if !app.InUse {
		opts += "s"
	}
	if app.Announcement != "" {
		opts += "A(" + app.Announcement + ")" // FIXME: Sanitize?
		if app.AnnounceToCaller {
			opts += "n"
		}
	}
	return opts
}
