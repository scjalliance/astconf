package astconf

import (
	"fmt"
	"io"
	"reflect"
)

// Marshaler is the interface implemented by types that can marshal
// themselves into valid asterisk configuration blocks.
//
// The provided encoder is only valid for the duration of the call. Marshalers
// must not retain a reference to e.
type Marshaler interface {
	MarshalAsterisk(e *Encoder) error
}

// ObjectMarshaler is the interface implemented by types that can marshal
// themselves into valid asterisk configuration objects.
type ObjectMarshaler interface {
	MarshalAsteriskObject(w io.Writer) error
}

// SettingMarshaler is the interface implemented by types that can marshal
// themselves into valid asterisk configuration settings.
type SettingMarshaler interface {
	MarshalAsteriskSetting(w io.Writer) error
}

// PreambleMarshaler is the interface implemented by types that marshal a
// an asterisk configuration preamble in addition to their regular value.
//
// The provided encoder is only valid for the duration of the call. Marshalers
// must not retain a reference to e.
type PreambleMarshaler interface {
	MarshalAsteriskPreamble(e *Encoder) error
}

// MarshalerError is returned when a custom marshaling function encounters an error.
type MarshalerError struct {
	Type reflect.Type
	Err  error
	call string // The type of marshaler
}

func (e *MarshalerError) Error() string {
	return fmt.Sprintf("astconf: error calling %s for type %s: %s", e.call, e.Type, e.Err)
}
