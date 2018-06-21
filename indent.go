package astconf

import "reflect"

type indentEncoder int

func (indent indentEncoder) encode(v reflect.Value, e *Encoder) error {
	if int(indent) > e.indent {
		e.indent = int(indent)
	}
	return nil
}

func newIndentEncoder(indent int) encoderFunc {
	return indentEncoder(indent).encode
}
