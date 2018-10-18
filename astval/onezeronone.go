package astval

import (
	"bytes"
	"fmt"
	"io"
)

// OneZeroNone is a boolean value that will be marshaled as "1" or "0".
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type OneZeroNone ternary

// Boolean states for One/Zero.
var (
	One  = OneZeroNone(ternaryTrue)
	Zero = OneZeroNone(ternaryFalse)
)

// True returns true if t specifies a true value.
func (t OneZeroNone) True() bool {
	return ternary(t).True()
}

// False returns true if t specifies a false value.
func (t OneZeroNone) False() bool {
	return ternary(t).False()
}

// Specified returns true if t specifies a value.
func (t OneZeroNone) Specified() bool {
	return ternary(t).Specified()
}

// MarshalText marshals the boolean as "1", "0" or nil.
func (t OneZeroNone) MarshalText() ([]byte, error) {
	switch {
	case ternary(t).True():
		return []byte(one), nil
	case ternary(t).False():
		return []byte(zero), nil
	default:
		return nil, nil
	}
}

// UnmarshalText parses "1", "0" or "" as a ternary value.
func (t *OneZeroNone) UnmarshalText(text []byte) error {
	switch {
	case len(text) == 0:
		*t = OneZeroNone(ternaryUnspecified)
		return nil
	case bytes.Equal(text, []byte(zero)):
		*t = OneZeroNone(ternaryFalse)
		return nil
	case bytes.Equal(text, []byte(one)):
		*t = OneZeroNone(ternaryTrue)
		return nil
	default:
		return fmt.Errorf("cannot marshal \"%s\" as 1/0 value")
	}
}

// MarshalAsteriskSetting writes the value of t to w as "1" or "0".
//
// The value of t will not be written when t is unspecified.
func (t OneZeroNone) MarshalAsteriskSetting(w io.Writer) error {
	switch {
	case ternary(t).False():
		_, err := w.Write([]byte(zero))
		return err
	case ternary(t).True():
		_, err := w.Write([]byte(one))
		return err
	default:
		return nil
	}
}
