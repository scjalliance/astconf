package astconf

import "reflect"

// Indenter is the interface implemented by types that can request a minimum
// level of indentation when marshaling.
type Indenter interface {
	AsteriskIndent() int
}

func indenterEncoder(v reflect.Value, e *Encoder) error {
	indenter, ok := v.Interface().(Indenter)
	if !ok {
		return nil // FIXME: Return an error? Panic?
	}
	indent := indenter.AsteriskIndent()
	if indent > e.indent {
		e.indent = indent
	}
	return nil
}
