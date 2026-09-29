package astconf

import "reflect"

type omitEmptyEncoder struct {
	elemEnc encoderFunc
}

func newOmitEmptyEncoder(elemEnc encoderFunc) encoderFunc {
	// TODO: Add support for types implementing the IsZero() interface?
	//       See: https://github.com/golang/go/issues/4357
	//            https://github.com/golang/go/issues/11939

	enc := omitEmptyEncoder{elemEnc: elemEnc}
	return enc.encode
}

func (oee omitEmptyEncoder) encode(v reflect.Value, e *Encoder) error {
	if isEmptyValue(v) {
		return nil
	}
	return oee.elemEnc(v, e)
}

// isEmptyValue reports whether v holds an empty value for the purposes of
// the omitempty tag option.
//
// A nil pointer is empty. A non-nil pointer is never empty, even when it
// points to a zero value, so that a pointer field can express an explicit
// zero. This matches encoding/json.
//
// Kinds not listed here are never considered empty.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}
