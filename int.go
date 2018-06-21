package astconf

import (
	"reflect"
	"strconv"
)

func intEncoder(v reflect.Value, e *Encoder) error {
	var scratch [64]byte
	_, err := e.w.Write(strconv.AppendInt(scratch[:0], v.Int(), 10))
	return err
}

func uintEncoder(v reflect.Value, e *Encoder) error {
	var scratch [64]byte
	_, err := e.w.Write(strconv.AppendUint(scratch[:0], v.Uint(), 10))
	return err
}
