package astconf

import (
	"reflect"
)

var (
	valTrue  = []byte("yes")
	valFalse = []byte("false")
)

func boolEncoder(v reflect.Value, e *Encoder) error {
	var err error
	if v.Bool() {
		_, err = e.w.Write(valTrue)
	} else {
		_, err = e.w.Write(valFalse)
	}
	return err
}
