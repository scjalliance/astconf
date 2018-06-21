package astconf

import (
	"fmt"
	"reflect"
)

type addrEncoder struct {
	elemEnc encoderFunc
}

// newAddrEncoder returns an encoder that performs elemEnc on the address
// of the value to be encoded.
func newAddrEncoder(elemEnc encoderFunc) encoderFunc {
	enc := addrEncoder{elemEnc: elemEnc}
	return enc.encode
}

func (ae *addrEncoder) encode(v reflect.Value, e *Encoder) error {
	if v.CanAddr() {
		va := v.Addr()
		if va.IsNil() {
			return nil
		}
		return ae.elemEnc(va, e)
	}
	return fmt.Errorf("unable to take address of %s", v)
}
