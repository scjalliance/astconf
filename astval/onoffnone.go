package astval

import (
	"bytes"
	"fmt"
	"io"
)

// OnOffNone is a ternary value that will be marshaled as "on" or "off".
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type OnOffNone ternary

// Boolean states for On/Off.
var (
	On  = OnOffNone(ternaryTrue)
	Off = OnOffNone(ternaryFalse)
)

// True returns true if t specifies a true value.
func (t OnOffNone) True() bool {
	return ternary(t).True()
}

// False returns true if t specifies a false value.
func (t OnOffNone) False() bool {
	return ternary(t).False()
}

// Specified returns true if t specifies a value.
func (t OnOffNone) Specified() bool {
	return ternary(t).Specified()
}

// MarshalText marshals the ternary value as "on", "off" or nil.
func (t OnOffNone) MarshalText() ([]byte, error) {
	switch {
	case ternary(t).True():
		return []byte(on), nil
	case ternary(t).False():
		return []byte(off), nil
	default:
		return nil, nil
	}
}

// UnmarshalText parses "on", "off" or "" as a ternary value.
func (t *OnOffNone) UnmarshalText(text []byte) error {
	switch {
	case len(text) == 0:
		*t = OnOffNone(ternaryUnspecified)
		return nil
	case bytes.Equal(text, []byte(on)):
		*t = OnOffNone(ternaryTrue)
		return nil
	case bytes.Equal(text, []byte(off)):
		*t = OnOffNone(ternaryFalse)
		return nil
	default:
		return fmt.Errorf("cannot marshal \"%s\" as on/off value")
	}
}

// MarshalAsteriskSetting writes the value of b to w as "on" or "off".
//
// The value of b will not be written when b is equal to the zero value of
// OnOff.
func (t OnOffNone) MarshalAsteriskSetting(w io.Writer) error {
	switch {
	case ternary(t).False():
		_, err := w.Write([]byte(off))
		return err
	case ternary(t).True():
		_, err := w.Write([]byte(on))
		return err
	default:
		return nil
	}
}
