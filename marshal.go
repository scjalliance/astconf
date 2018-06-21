package astconf

import (
	"bytes"
	"io"
)

// MarshalTo marshals the asterisk configuration of v to w.
func MarshalTo(v interface{}, w io.Writer) error {
	return NewEncoder(w).Encode(v)
}

// Marshal marshals the asterisk configuration of v and returns it as a slice
// of bytes.
//
// Marshal will return an error if asked to marshal a channel, function, or
// map.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
