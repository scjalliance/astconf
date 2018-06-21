package astconf

import (
	"reflect"
)

type sliceEncoder struct {
	elemEnc encoderFunc
}

func newSliceEncoder(t reflect.Type) encoderFunc {
	// TODO: Handle byte slices in some special way?
	enc := sliceEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

func (s sliceEncoder) encode(v reflect.Value, e *Encoder) error {
	for i, n := 0, v.Len(); i < n; i++ {
		e := e.derive()
		if err := s.elemEnc(v.Index(i), e); err != nil {
			return err
		}
	}
	return nil
}
