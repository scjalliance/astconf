package astconf

import "reflect"

type blockEncoder struct {
	elemEnc encoderFunc
}

func newBlockEncoder(elemEnc encoderFunc) encoderFunc {
	enc := blockEncoder{elemEnc: elemEnc}
	return enc.encode
}

func (be *blockEncoder) encode(v reflect.Value, e *Encoder) error {
	e = e.derive()
	return be.elemEnc(v, e)
}
