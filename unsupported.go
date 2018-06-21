package astconf

import (
	"reflect"
)

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "astconf: unsupported type: " + e.Type.String()
}

func unsupportedTypeEncoder(v reflect.Value, e *Encoder) error {
	return &UnsupportedTypeError{v.Type()}
}
