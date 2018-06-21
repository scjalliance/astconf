package astconf

import "reflect"

type objectEncoder struct {
	name    []byte // []byte(name)
	elemEnc encoderFunc
}

func newObjectEncoder(name string, elemEnc encoderFunc) encoderFunc {
	enc := &objectEncoder{name: []byte(name), elemEnc: elemEnc}
	return enc.encode
}

func (oe *objectEncoder) encode(v reflect.Value, e *Encoder) error {
	e = e.derive()
	e.scratch.name = oe.name
	e.scratch.sep = objectSeparator
	e.scratch.started = false
	fw := fieldWriter{e}
	e.w = fw
	if err := oe.elemEnc(v, e); err != nil {
		return err
	}
	return fw.done()
}
