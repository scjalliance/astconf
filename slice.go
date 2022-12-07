package astconf

import (
	"reflect"
)

type sliceEncoder struct {
	elemEnc encoderFunc
}

func newSliceEncoder(elemEnc encoderFunc) encoderFunc {
	enc := sliceEncoder{elemEnc: elemEnc}
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

type commaSeparatedSliceEncoder struct {
	elemEnc encoderFunc
}

func newCommaSeparatedSliceEncoder(elemEnc encoderFunc) encoderFunc {
	enc := commaSeparatedSliceEncoder{elemEnc: elemEnc}
	return enc.encode
}

func (s commaSeparatedSliceEncoder) encode(v reflect.Value, e *Encoder) error {
	out := e.Printer()
	for i, n := 0, v.Len(); i < n; i++ {
		if i > 0 {
			out.Write([]byte(","))
		}
		if err := s.elemEnc(v.Index(i), e); err != nil {
			return err
		}
	}
	return nil
}
