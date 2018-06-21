package astconf

import "reflect"

// An InvalidValueError is returned by Marshal when attempting
// to encode an invalid value.
type InvalidValueError struct {
	Value reflect.Value
}

func (e *InvalidValueError) Error() string {
	return "astconf: invalid value for type: " + e.Value.Type().String()
}

func invalidValueEncoder(v reflect.Value, e *Encoder) error {
	return &InvalidValueError{v}
}
