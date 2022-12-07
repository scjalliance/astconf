package astconf

import (
	"reflect"
)

func newStructEncoder(t reflect.Type, canAddr bool) encoderFunc {
	// Note that if t provides its own custom marshaling that will be
	// handled by newTypeEncoder, not newStructEncoder.

	// We assume throughout this function that even if we can't take the
	// address of t, we can take the address of t's fields.

	var encoders encoderSet

	features := typeFeatures(t)
	fields := typeFields(t, canAddr)

	// Step 1a: Calculate the indent of t
	switch {
	case features.Indenter():
		encoders = append(encoders, indenterEncoder)
	case canAddr && features.IndenterAddr():
		encoders = append(encoders, newAddrEncoder(indenterEncoder))
	}

	// Step 1b: Calculate the indent of t's fields
	var indent int
	for i := range fields {
		f := &fields[i]

		// Variable indentation
		switch {
		case f.Indenter():
			encoders = append(encoders, newFieldEncoder(f.index, indenterEncoder))
		case f.IndenterAddr():
			encoders = append(encoders, newFieldEncoder(f.index, newAddrEncoder(indenterEncoder)))
		}
		if f.embedded || (f.Block() && !f.MultiValue()) || (f.BlockAddr() && !f.MultiValueAddr()) {
			// It doesn't have a printed field name, so don't count it
			continue
		}

		// Fixed indentation
		nlen := len(fields[i].name)
		if nlen > indent {
			indent = nlen
		}
	}
	encoders = append(encoders, newIndentEncoder(indent))

	// Step 2: Print the section name
	switch {
	case features.Sectioner():
		encoders = append(encoders, sectionEncoder)
	case canAddr && features.SectionerAddr():
		encoders = append(encoders, newAddrEncoder(sectionEncoder))
	default:
		var sectioner *field
		for i := range fields {
			f := &fields[i]
			if f.Sectioner() || f.SectionerAddr() {
				if sectioner == nil || len(sectioner.index) >= len(f.index) {
					sectioner = f
				}
			}
		}
		if sectioner != nil {
			switch {
			case sectioner.Sectioner():
				encoders = append(encoders, newFieldEncoder(sectioner.index, sectionEncoder))
			case sectioner.SectionerAddr():
				encoders = append(encoders, newFieldEncoder(sectioner.index, newAddrEncoder(sectionEncoder)))
			}
		}
	}

	// Step 3: Execute all preamble marshalers
	switch {
	case features.PreambleMarshaler():
		encoders = append(encoders, mtPreamble.encode)
	case canAddr && features.PreambleMarshalerAddr():
		encoders = append(encoders, mtPreamble.addrEncode)
	}

	for i := range fields {
		f := &fields[i]
		switch {
		case f.PreambleMarshaler():
			encoders = append(encoders, newFieldEncoder(f.index, mtPreamble.encode))
		case f.PreambleMarshalerAddr():
			encoders = append(encoders, newFieldEncoder(f.index, mtPreamble.addrEncode))
		}
	}

	// Step 4: Marshal all field values
	for i := range fields {
		f := &fields[i]
		if f.Omit() || f.embedded {
			continue
		}
		t := typeByIndex(t, f.index)
		var elemEnc encoderFunc
		switch {
		case f.MultiValue(), f.MultiValueAddr():
			if f.CommaSeparated() && !f.Object() && !f.ObjectAddr() {
				// Special handling for multi-valued fields with the
				// commaseparated tag
				elemEnc = newCommaSeparatedSliceEncoder(typeEncoder(t.Elem()))
				elemEnc = newSettingEncoder(f.name, elemEnc)
			} else {
				// Standard handling for multi-valued fields
				switch {
				case f.Object(), f.ObjectAddr():
					elemEnc = newObjectEncoder(f.name, typeEncoder(t.Elem()))
				default:
					elemEnc = newSettingEncoder(f.name, typeEncoder(t.Elem()))
				}
				elemEnc = newSliceEncoder(elemEnc)
			}
		case f.Block(), f.BlockAddr():
			elemEnc = newBlockEncoder(typeEncoder(t))
		case f.Object(), f.ObjectAddr():
			elemEnc = newObjectEncoder(f.name, typeEncoder(t))
		default:
			elemEnc = newSettingEncoder(f.name, typeEncoder(t))
		}
		if f.OmitEmtpy() {
			elemEnc = newOmitEmptyEncoder(t, elemEnc)
		}
		encoders = append(encoders, newFieldEncoder(f.index, elemEnc))
	}

	return encoders.finalize()
}
