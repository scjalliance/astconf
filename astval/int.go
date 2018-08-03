package astval

import (
	"io"
	"strconv"
)

// Int is an integer value.
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type Int struct {
	value     int
	specified bool
}

// NewInt returns an integer with the given value.
func NewInt(value int) Int {
	return Int{value: value, specified: true}
}

// Set applies the given value to i.
func (i *Int) Set(value int) {
	i.value = value
	i.specified = true
}

// Value returns the integer value of i.
func (i Int) Value() int {
	return i.value
}

// Specified returns true if the value of i has been specified.
func (i Int) Specified() bool {
	return i.specified
}

// MarshalAsteriskSetting writes the value of i to w.
//
// The value of i will not be written when i is equal to the zero value of
// Int.
func (i Int) MarshalAsteriskSetting(w io.Writer) error {
	if !i.specified {
		return nil
	}
	var scratch [64]byte
	_, err := w.Write(strconv.AppendInt(scratch[:0], int64(i.value), 10))
	return err
}
