package dpma

import (
	"fmt"
	"strings"

	"github.com/scjalliance/astconf"
)

// LogoFile holds the path to a logo for a particular model of phone.
type LogoFile struct {
	Model string
	Path  string
}

// MarshalAsterisk marshals the logo file as an asterisk setting.
func (logo LogoFile) MarshalAsterisk(e *astconf.Encoder) error {
	name := fmt.Sprintf("%s_logo_file", strings.ToLower(logo.Model))
	return e.Printer().Setting(name, logo.Path)
}
