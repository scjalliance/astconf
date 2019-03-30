package dialplan

import (
	"fmt"

	"github.com/scjalliance/astconf"
)

// Extension is a dialplan extension.
type Extension struct {
	Comment string
	Number  string
	Actions []Action
}

// MarshalAsterisk marshals the extension to an asterisk encoder.
func (exten Extension) MarshalAsterisk(e *astconf.Encoder) error {
	if len(exten.Actions) == 0 {
		return nil
	}

	p := e.Printer()

	// Extensions typically are multi-line, so we give them some space
	p.Break()

	// If the extension has a comment we'll print it
	if exten.Comment != "" {
		p.Comment(exten.Comment)
	}

	// Print each of the actions in the extension
	for i, action := range exten.Actions {
		var obj, value string
		if i == 0 {
			obj = "exten"
			value = fmt.Sprintf("%s,%d,%s", exten.Number, i+1, action.App().String())
		} else {
			obj = "same"
			value = fmt.Sprintf("n,%s", action.App().String())
		}

		if err := e.Printer().Object(obj, value); err != nil {
			return err
		}
	}

	return nil
}
