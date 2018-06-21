package astconf

import (
	"reflect"
)

func typeEncoder(t reflect.Type) encoderFunc {
	// TODO: Add caching
	return newTypeEncoder(t, true)
}

// newTypeEncoder constructs an encoderFunc for a type.
// The returned encoder only checks CanAddr when allowAddr is true.
func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
	features := typeFeatures(t)

	switch {
	case features.BlockMarshaler():
		return mtBlock.encode
	case allowAddr && features.BlockMarshalerAddr():
		return newCondAddrEncoder(mtBlock.addrEncode, newTypeEncoder(t, false))
	case features.ObjectMarshaler():
		return mtObject.encode
	case allowAddr && features.ObjectMarshalerAddr():
		return newCondAddrEncoder(mtObject.addrEncode, newTypeEncoder(t, false))
	case features.SettingMarshaler():
		return mtSetting.encode
	case allowAddr && features.SettingMarshalerAddr():
		return newCondAddrEncoder(mtSetting.addrEncode, newTypeEncoder(t, false))
	case features.TextMarshaler():
		return mtText.encode
	case allowAddr && features.TextMarshalerAddr():
		return newCondAddrEncoder(mtText.addrEncode, newTypeEncoder(t, false))
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.String:
		return stringEncoder
	case reflect.Struct:
		if allowAddr {
			return newCondAddrEncoder(newStructEncoder(t, true), newStructEncoder(t, false))
		}
		return newStructEncoder(t, false)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Ptr:
		return newIndirectEncoder(newTypeEncoder(t.Elem(), allowAddr))
	default:
		return unsupportedTypeEncoder
	}
}
