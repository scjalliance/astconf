package astconf

import (
	"encoding"
	"errors"
	"reflect"
)

type marshalerEncoder int

const (
	mtBlock marshalerEncoder = iota
	mtSetting
	mtObject
	mtText
	mtPreamble
)

var (
	errMarshalerNotImplemented         = errors.New("type does not implement Marshaler as expected")
	errObjectMarshalerNotImplemented   = errors.New("type does not implement ObjectMarshaler as expected")
	errSettingMarshalerNotImplemented  = errors.New("type does not implement ValueMarshaler as expected")
	errTextMarshalerNotImplemented     = errors.New("type does not implement TextMarshaler as expected")
	errPreambleMarshalerNotImplemented = errors.New("type does not implement PreambleMarshaler as expected")
)

func (me marshalerEncoder) encode(v reflect.Value, e *Encoder) error {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	var err error
	switch me {
	case mtBlock:
		m, ok := v.Interface().(Marshaler)
		if ok {
			err = m.MarshalAsterisk(e)
		} else {
			err = errMarshalerNotImplemented
		}
	case mtObject:
		m, ok := v.Interface().(ObjectMarshaler)
		if ok {
			err = m.MarshalAsteriskObject(e.w)
		} else {
			err = errObjectMarshalerNotImplemented
		}
	case mtSetting:
		m, ok := v.Interface().(SettingMarshaler)
		if ok {
			err = m.MarshalAsteriskSetting(e.w)
		} else {
			err = errSettingMarshalerNotImplemented
		}
	case mtText:
		m, ok := v.Interface().(encoding.TextMarshaler)
		if ok {
			var b []byte
			b, err = m.MarshalText()
			if err == nil {
				_, err = e.w.Write(b)
			}
		} else {
			err = errTextMarshalerNotImplemented
		}
	case mtPreamble:
		m, ok := v.Interface().(PreambleMarshaler)
		if ok {
			err = m.MarshalAsteriskPreamble(e)
		} else {
			err = errPreambleMarshalerNotImplemented
		}
	}
	if err != nil {
		return &MarshalerError{v.Type(), err, me.Call()}
	}
	return nil
}

func (me marshalerEncoder) addrEncode(v reflect.Value, e *Encoder) error {
	va := v.Addr()
	if va.IsNil() {
		return nil
	}
	return me.encode(va, e)
}

func (me marshalerEncoder) Call() string {
	switch me {
	case mtBlock:
		return "MarshalAsterisk"
	case mtObject:
		return "MarshalAsteriskObject"
	case mtSetting:
		return "MarshalAsteriskSetting"
	case mtText:
		return "MarshalText"
	case mtPreamble:
		return "MarshalAsteriskPreamble"
	default:
		return "UnknownMarshalingFunction"
	}
}
