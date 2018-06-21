package astconf

import "reflect"

type encoderFunc func(v reflect.Value, e *Encoder) error

type encoderSet []encoderFunc

func (encoders encoderSet) encode(v reflect.Value, e *Encoder) error {
	for i := range encoders {
		if err := encoders[i](v, e); err != nil {
			return err
		}
	}
	return nil
}

func (encoders encoderSet) finalize() encoderFunc {
	if len(encoders) == 1 {
		return encoders[0]
	}
	return encoders.encode
}
