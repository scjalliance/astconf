package astval

import (
	"bytes"
	"fmt"
)

const (
	yes = "yes"
	no  = "no"
)

// YesNo is a boolean value that will be marshaled as "yes" or "no".
type YesNo bool

// String returns the boolean as the string "yes" or "no".
func (b YesNo) String() string {
	if b {
		return yes
	}
	return no
}

// MarshalText marshals the boolean the boolean as "yes" or "no".
func (b YesNo) MarshalText() ([]byte, error) {
	if b {
		return []byte(yes), nil
	}
	return []byte(no), nil
}

// UnmarshalText parses "yes" or "no" as a boolean value.
func (b *YesNo) UnmarshalText(text []byte) error {
	switch {
	case bytes.Equal(text, []byte(yes)):
		*b = YesNo(true)
		return nil
	case bytes.Equal(text, []byte(no)):
		*b = YesNo(false)
		return nil
	default:
		return fmt.Errorf("cannot marshal \"%s\" as yes/no value", string(text))
	}
}
