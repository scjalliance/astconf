package astconf

import "reflect"

func sectionEncoder(v reflect.Value, e *Encoder) error {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil // FIXME: Return some sort of invalid section error?
	}
	namer, ok := v.Interface().(SectionNamer)
	if !ok {
		return nil // FIXME: Return some sort of invalid section error? panic?
	}
	name := namer.SectionName()

	if name == "" {
		// TODO: Decide whether we should return an error, panic, or
		//       silently skip the section header.
		return nil
	}

	return e.Printer().Section(name)
}
