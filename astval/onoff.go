package astval

import (
	"bytes"
	"fmt"
)

const (
	on  = "on"
	off = "off"
)

// OnOff is a boolean value that will be marshaled as "on" or "off".
type OnOff bool

// String returns the booelan as the string "on" or "off".
func (b OnOff) String() string {
	if b {
		return on
	}
	return off
}

// MarshalText marshals the boolean as "on" or "off".
func (b OnOff) MarshalText() ([]byte, error) {
	if b {
		return []byte(on), nil
	}
	return []byte(off), nil
}

// UnmarshalText parses "on" or "off" as a boolean value.
func (b *OnOff) UnmarshalText(text []byte) error {
	switch {
	case bytes.Equal(text, []byte(on)):
		*b = OnOff(true)
		return nil
	case bytes.Equal(text, []byte(off)):
		*b = OnOff(false)
		return nil
	default:
		return fmt.Errorf("cannot marshal \"%s\" as on/off value", string(text))
	}
}
