package dialplan

import (
	"strings"

	"github.com/scjalliance/astconf"
)

// Validator is the interface implemented by applications, functions and
// expressions that can check their literal content before it is rendered.
//
// Extension.MarshalAsterisk validates each action that implements it, so
// text that would change the structure of a dialplan line is rejected
// instead of being written.
type Validator interface {
	Validate() error
}

// Characters that would change the structure of a dialplan line if they
// appeared in literal text. Commas separate arguments, parentheses delimit
// calls, brackets delimit expressions, and $ begins a substitution.
const (
	invalidArgChars       = ",()$[]"
	invalidExtensionChars = ",()$"     // Extension patterns use brackets
	invalidResourceChars  = ",()$[]&/" // & separates devices and / separates resource parts
	invalidFileChars      = ",()$[]&"  // & separates files
)

func errorIfAny(component, value, badChars string) error {
	if pos := strings.IndexAny(value, badChars); pos >= 0 {
		return astconf.InvalidContentError{
			Component: component,
			Value:     value,
			Pos:       pos,
		}
	}
	return nil
}

// validate calls Validate on v if v implements Validator.
func validate(v interface{}) error {
	if validator, ok := v.(Validator); ok {
		return validator.Validate()
	}
	return nil
}
