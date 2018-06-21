package astconf

import "reflect"

type settingEncoder struct {
	name    []byte // []byte(name)
	elemEnc encoderFunc
}

func newSettingEncoder(name string, elemEnc encoderFunc) encoderFunc {
	enc := &settingEncoder{name: []byte(name), elemEnc: elemEnc}
	return enc.encode
}

func (se settingEncoder) encode(v reflect.Value, e *Encoder) error {
	e = e.derive()
	e.scratch.name = se.name
	e.scratch.sep = settingSeparator
	e.scratch.started = false
	fw := fieldWriter{e}
	e.w = fw
	if err := se.elemEnc(v, e); err != nil {
		return err
	}
	return fw.done()
}
