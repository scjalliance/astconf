package astconf

import (
	"reflect"
)

func stringEncoder(v reflect.Value, e *Encoder) error {
	_, err := e.w.Write([]byte(v.String()))
	return err
}
