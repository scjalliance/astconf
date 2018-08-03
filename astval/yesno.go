package astval

import "io"

// YesNo is a boolean value that will be marshaled as "yes" or "no".
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type YesNo struct {
	state     bool
	specified bool
}

// Boolean states
var (
	Yes = YesNo{state: true, specified: true}
	No  = YesNo{state: false, specified: true}
)

// Specified returns true if the value of b has been specified.
func (b YesNo) Specified() bool {
	return b.specified
}

// MarshalAsteriskSetting marshals the boolean as "yes" or "no" to w.
func (b YesNo) MarshalAsteriskSetting(w io.Writer) error {
	if !b.specified {
		return nil
	}
	var err error
	if b.state {
		_, err = w.Write([]byte("yes"))
	} else {
		_, err = w.Write([]byte("no"))
	}
	return err
}
