package astval

import (
	"io"
	"strconv"
	"time"
)

// Seconds is a time duration that will be marshaled as an integer number of
// seconds.
//
// Its zero value indicates an unspecified condition that will not be
// marshaled.
type Seconds struct {
	value     time.Duration
	specified bool
}

// NewSeconds returns a time duration that will be marshaled as an integer
// number of seconds.
func NewSeconds(value time.Duration) Seconds {
	return Seconds{value: value, specified: true}
}

// Set applies the given time duration to seconds.
func (seconds *Seconds) Set(value time.Duration) {
	seconds.value = value
	seconds.specified = true
}

// Value returns the value of seconds as a duration.
func (seconds Seconds) Value() time.Duration {
	return seconds.value
}

// Specified returns true if the value of seconds has been specified.
func (seconds Seconds) Specified() bool {
	return seconds.specified
}

// MarshalAsteriskSetting writes the value of seconds to w as an integral
// number of seconds.
//
// The value of seconds will not be written when seconds is equal to the
// zero value of Seconds.
func (seconds Seconds) MarshalAsteriskSetting(w io.Writer) error {
	if !seconds.specified {
		return nil
	}
	_, err := w.Write([]byte(strconv.FormatInt(int64(seconds.value/time.Second), 10)))
	return err
}
