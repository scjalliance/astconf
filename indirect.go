package astconf

import "reflect"

type indirectEncoder struct {
	elemEnc encoderFunc
}

// newIndirectEncoder returns an encoder that performs elemEnc on the type
// pointed to by value.
func newIndirectEncoder(elemEnc encoderFunc) encoderFunc {
	enc := &indirectEncoder{elemEnc: elemEnc}
	return enc.encode
}

func (ie *indirectEncoder) encode(v reflect.Value, e *Encoder) error {
	if v.IsNil() {
		return nil
	}
	return ie.elemEnc(v.Elem(), e)
}

/*
func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
}
*/
