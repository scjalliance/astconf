package astconf

// fieldWriter is an io.Writer that wraps an Encoder. It will write the
// field name ahead of any field value written via fieldWriter.Write. The
// field name to be written will be pulled from the encoder's scratch data,
// which should have been prepared by a corresponding setting or object
// encoder.
//
// When a fieldEncoder is marshaling a particular field, it creates an
// ephemeral fieldWriter which is then passed to the field's type encoder,
// which may be a custom marshaler implementation.
//
// The first time fieldWriter.Write is called, it writes the field's
// indentation, name and field separator (typically = or =>). If Write is
// never called, the field name information is not written. This design
// allows custom type marshalers to omit a field simply by not writing a
// value for it.
//
// Every value written through fieldWriter is checked for characters that
// would break the line, and semicolons are escaped, so that all setting
// and object values are validated regardless of which encoder or custom
// marshaler produced them.
type fieldWriter struct {
	*Encoder
}

// Write will write field value p to the underlying writer.
func (fw fieldWriter) Write(p []byte) (n int, err error) {
	if err := errorIfAnyBytes("value", p, invalidValueChars); err != nil {
		return 0, err
	}
	escaped := escapeValue(p)

	e := fw.Encoder
	s := &e.scratch
	if !s.started {
		// fw == e.w at this point, so don't use e.w here or we'll loop forever
		printer := Printer{
			Writer:    e.base,
			Indent:    e.indent,
			Alignment: e.alignment,
			Started:   e.started,
		}

		if err := printer.start(s.name, s.sep); err != nil {
			return 0, err
		}

		s.started = true
		e.w = e.base
	}

	n, err = e.w.Write(escaped)
	if n > len(p) {
		// Escaping may have lengthened the output, but io.Writer requires
		// that n not exceed len(p).
		n = len(p)
	}
	return n, err
}

func (fw *fieldWriter) done() error {
	e := fw.Encoder
	if e.scratch.started {
		_, err := e.w.Write(newline)
		return err
	}
	return nil
}
