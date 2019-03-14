package astval

import (
	"bytes"
	"fmt"
	"io"
)

// YesNoNone is a ternary value that will be marshaled as "yes" or "no".
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type YesNoNone ternary

// Boolean states for Yes/No.
var (
	Yes = YesNoNone(ternaryTrue)
	No  = YesNoNone(ternaryFalse)
)

// True returns true if t specifies a true value.
func (t YesNoNone) True() bool {
	return ternary(t).True()
}

// False returns true if t specifies a false value.
func (t YesNoNone) False() bool {
	return ternary(t).False()
}

// Specified returns true if t specifies a value.
func (t YesNoNone) Specified() bool {
	return ternary(t).Specified()
}

// MarshalText marshals the ternary value as "on", "off" or nil.
func (t YesNoNone) MarshalText() ([]byte, error) {
	switch {
	case ternary(t).True():
		return []byte(yes), nil
	case ternary(t).False():
		return []byte(no), nil
	default:
		return nil, nil
	}
}

// UnmarshalText parses "on", "off" or "" as a ternary value.
func (t *YesNoNone) UnmarshalText(text []byte) error {
	switch {
	case len(text) == 0:
		*t = YesNoNone(ternaryUnspecified)
		return nil
	case bytes.Equal(text, []byte(no)):
		*t = YesNoNone(ternaryFalse)
		return nil
	case bytes.Equal(text, []byte(yes)):
		*t = YesNoNone(ternaryTrue)
		return nil
	default:
		return fmt.Errorf("cannot marshal \"%s\" as on/off value", string(text))
	}
}

// MarshalAsteriskSetting writes the value of b to w as "on" or "off".
//
// The value of b will not be written when b is equal to the zero value of
// OnOff.
func (t YesNoNone) MarshalAsteriskSetting(w io.Writer) error {
	switch {
	case ternary(t).False():
		_, err := w.Write([]byte(no))
		return err
	case ternary(t).True():
		_, err := w.Write([]byte(yes))
		return err
	default:
		return nil
	}
}
