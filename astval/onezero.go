package astval

import (
	"bytes"
	"fmt"
)

const (
	one  = "1"
	zero = "0"
)

// OneZero is a boolean value that will be marshaled as "1" or "0".
type OneZero bool

// String returns the booelan as "1" or "0".
func (b OneZero) String() string {
	if b {
		return one
	}
	return zero
}

// MarshalText marshals the boolean as "1" or "0".
func (b OneZero) MarshalText() ([]byte, error) {
	if b {
		return []byte(one), nil
	}
	return []byte(zero), nil
}

// UnmarshalText parses "1" or "0" as a boolean value.
func (b *OneZero) UnmarshalText(text []byte) error {
	switch {
	case bytes.Equal(text, []byte(one)):
		*b = OneZero(true)
		return nil
	case bytes.Equal(text, []byte(zero)):
		*b = OneZero(false)
		return nil
	default:
		return fmt.Errorf("cannot unmarshal \"%s\" as 1/0 value", string(text))
	}
}
