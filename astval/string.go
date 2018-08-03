package astval

import (
	"io"
)

// String is a string value.
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type String struct {
	value     string
	specified bool
}

// NewString returns an string with the given value.
func NewString(value int) Int {
	return Int{value: value, specified: true}
}

// Set applies value to i.
func (s *String) Set(value string) {
	s.value = value
	s.specified = true
}

// Value returns the string value of s.
func (s String) Value() string {
	return s.value
}

// Specified returns true if the value of s has been specified.
func (s String) Specified() bool {
	return s.specified
}

// MarshalAsteriskSetting writes the value of s to w.
//
// The value of s will not be written when s is equal to the zero value of
// String.
func (s String) MarshalAsteriskSetting(w io.Writer) error {
	if !s.specified {
		return nil
	}
	_, err := w.Write([]byte(s.value))
	return err
}
