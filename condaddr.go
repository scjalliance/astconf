package astconf

import "reflect"

type condAddrEncoder struct {
	canAddrEnc, elseEnc encoderFunc
}

// newCondAddrEncoder returns an encoder that checks whether its value
// CanAddr and delegates to canAddrEnc if so, else to elseEnc.
func newCondAddrEncoder(canAddrEnc, elseEnc encoderFunc) encoderFunc {
	enc := &condAddrEncoder{canAddrEnc: canAddrEnc, elseEnc: elseEnc}
	return enc.encode
}

func (ce *condAddrEncoder) encode(v reflect.Value, e *Encoder) error {
	if v.CanAddr() {
		return ce.canAddrEnc(v, e)
	}
	return ce.elseEnc(v, e)
}
