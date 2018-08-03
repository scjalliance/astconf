package astval

import "io"

// OnOff is a boolean value that will be marshaled as "on" or "off".
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type OnOff struct {
	state     bool
	specified bool
}

// Boolean states
var (
	On  = OnOff{state: true, specified: true}
	Off = OnOff{state: false, specified: true}
)

// MarshalAsteriskSetting writes the value of b to w as "on" or "off".
//
// The value of b will not be written when b is equal to the zero value of
// OnOff.
func (b OnOff) MarshalAsteriskSetting(w io.Writer) error {
	if !b.specified {
		return nil
	}
	var err error
	if b.state {
		_, err = w.Write([]byte("on"))
	} else {
		_, err = w.Write([]byte("off"))
	}
	return err
}
