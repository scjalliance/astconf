package astconf

import "reflect"

type omitEmptyEncoder struct {
	elemEnc encoderFunc
}

func newOmitEmptyEncoder(t reflect.Type, elemEnc encoderFunc) encoderFunc {
	// TODO: Add support for types implementing the IsZero() interface?
	//       See: https://github.com/golang/go/issues/4357
	//            https://github.com/golang/go/issues/11939

	enc := omitEmptyEncoder{elemEnc: elemEnc}
	switch t.Kind() {
	case reflect.Bool:
		return enc.encodeBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return enc.encodeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return enc.encodeUint
	case reflect.String:
		return enc.encodeString
	case reflect.Slice:
		return enc.encodeSlice
	case reflect.Ptr:
		return newIndirectEncoder(newOmitEmptyEncoder(t.Elem(), elemEnc))
	default:
		return elemEnc
	}
}

func (oee omitEmptyEncoder) encodeBool(v reflect.Value, e *Encoder) error {
	if v.Bool() == false {
		return nil
	}
	return oee.elemEnc(v, e)
}

func (oee omitEmptyEncoder) encodeInt(v reflect.Value, e *Encoder) error {
	if v.Int() == 0 {
		return nil
	}
	return oee.elemEnc(v, e)
}

func (oee omitEmptyEncoder) encodeUint(v reflect.Value, e *Encoder) error {
	if v.Uint() == 0 {
		return nil
	}
	return oee.elemEnc(v, e)
}

func (oee omitEmptyEncoder) encodeString(v reflect.Value, e *Encoder) error {
	if v.String() == "" {
		return nil
	}
	return oee.elemEnc(v, e)
}

func (oee omitEmptyEncoder) encodeSlice(v reflect.Value, e *Encoder) error {
	if v.Len() == 0 {
		return nil
	}
	return oee.elemEnc(v, e)
}
