package astconf

import (
	"encoding"
	"reflect"
)

var (
	blockMarshalerType    = reflect.TypeOf(new(Marshaler)).Elem()
	objectMarshalerType   = reflect.TypeOf(new(ObjectMarshaler)).Elem()
	settingMarshalerType  = reflect.TypeOf(new(SettingMarshaler)).Elem()
	textMarshalerType     = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
	preambleMarshalerType = reflect.TypeOf(new(PreambleMarshaler)).Elem()
	sectionerType         = reflect.TypeOf(new(SectionNamer)).Elem()
	indenterType          = reflect.TypeOf(new(Indenter)).Elem()
)

type typeFeature int

const (
	tfBlockMarshaler typeFeature = 1 << iota
	tfBlockMarshalerAddr
	tfObjectMarshaler
	tfObjectMarshalerAddr
	tfSettingMarshaler
	tfSettingMarshalerAddr
	tfTextMarshaler
	tfTextMarshalerAddr
	tfPreambleMarshaler
	tfPreambleMarshalerAddr
	tfSectioner
	tfSectionerAddr
	tfIndenter
	tfIndenterAddr
	tfMultiValue
	tfMultiValueAddr
	tfBlock
	tfBlockAddr
	tfObject
	tfObjectAddr
	tfOmit
	tfOmitEmpty

	tfMarshalerMask     = tfBlockMarshaler | tfObjectMarshaler | tfSettingMarshaler | tfTextMarshaler
	tfMarshalerAddrMask = tfBlockMarshalerAddr | tfObjectMarshalerAddr | tfSettingMarshalerAddr | tfTextMarshalerAddr
)

func tagFeatures(tag *fieldTag) (features typeFeature) {
	if tag.omit {
		features |= tfOmit
	}
	if tag.Contains("omitempty") {
		features |= tfOmitEmpty
	}
	if tag.Contains("object") {
		features |= tfObject
	}
	return features
}

func typeFeatures(t reflect.Type) (features typeFeature) {
	// Marshaler
	switch {
	case typeImplements(t, blockMarshalerType):
		features |= tfBlockMarshaler
		features |= tfBlock
	case typeAddrImplements(t, blockMarshalerType):
		features |= tfBlockMarshalerAddr
		features |= tfBlockAddr
	}

	// ObjectMarshaler
	switch {
	case typeImplements(t, objectMarshalerType):
		features |= tfObjectMarshaler
		features |= tfObject
	case typeAddrImplements(t, objectMarshalerType):
		features |= tfObjectMarshalerAddr
		features |= tfObjectAddr
	}

	// SettingMarshaler
	switch {
	case typeImplements(t, settingMarshalerType):
		features |= tfSettingMarshaler
	case typeAddrImplements(t, settingMarshalerType):
		features |= tfSettingMarshalerAddr
	}

	// TextMarshaler
	switch {
	case typeImplements(t, textMarshalerType):
		features |= tfTextMarshaler
	case typeAddrImplements(t, textMarshalerType):
		features |= tfTextMarshalerAddr
	}

	// PreambleMarshaler
	switch {
	case typeImplements(t, preambleMarshalerType):
		features |= tfPreambleMarshaler
	case typeAddrImplements(t, preambleMarshalerType):
		features |= tfPreambleMarshalerAddr
	}

	// Sectioner
	switch {
	case typeImplements(t, sectionerType):
		features |= tfSectioner
	case typeAddrImplements(t, sectionerType):
		features |= tfSectionerAddr
	}

	// Indenter
	switch {
	case typeImplements(t, indenterType):
		features |= tfIndenter
	case typeAddrImplements(t, indenterType):
		features |= tfIndenterAddr
	}

	// MultiValue
	switch t.Kind() {
	case reflect.Slice:
		// Slices of non-blocks get written out as multiple values with the
		// same field name, unless they perform their own marshaling.
		elemFeatures := typeFeatures(t.Elem())
		if !features.Marshaler() && !elemFeatures.Block() {
			features |= tfMultiValue
		}
		if !features.Marshaler() && !features.MarshalerAddr() && !elemFeatures.Block() && !elemFeatures.BlockAddr() {
			features |= tfMultiValueAddr
		}
	}

	// Block
	switch t.Kind() {
	case reflect.Struct, reflect.Slice:
		// Structs and slices are considered blocks unless they perform
		// their own marshaling.
		if !features.Marshaler() {
			features |= tfBlock
		}
		if !features.Marshaler() && !features.MarshalerAddr() {
			features |= tfBlockAddr
		}
	case reflect.Ptr:
		// Pointers are considered blocks if their elements are, unless
		// they perform their own marshaling.
		elemFeatures := typeFeatures(t.Elem())
		if !features.Marshaler() && !elemFeatures.Block() {
			features |= tfBlock
		}
		if !features.Marshaler() && !features.MarshalerAddr() && !elemFeatures.Block() && !elemFeatures.BlockAddr() {
			features |= tfBlockAddr
		}
	}

	return features
}

func (tf typeFeature) BlockMarshaler() bool {
	return tf&tfBlockMarshaler != 0
}

func (tf typeFeature) BlockMarshalerAddr() bool {
	return tf&tfBlockMarshalerAddr != 0
}

func (tf typeFeature) ObjectMarshaler() bool {
	return tf&tfObjectMarshaler != 0
}

func (tf typeFeature) ObjectMarshalerAddr() bool {
	return tf&tfObjectMarshalerAddr != 0
}

func (tf typeFeature) SettingMarshaler() bool {
	return tf&tfSettingMarshaler != 0
}

func (tf typeFeature) SettingMarshalerAddr() bool {
	return tf&tfSettingMarshalerAddr != 0
}

func (tf typeFeature) TextMarshaler() bool {
	return tf&tfTextMarshaler != 0
}

func (tf typeFeature) TextMarshalerAddr() bool {
	return tf&tfTextMarshalerAddr != 0
}

func (tf typeFeature) PreambleMarshaler() bool {
	return tf&tfPreambleMarshaler != 0
}

func (tf typeFeature) PreambleMarshalerAddr() bool {
	return tf&tfPreambleMarshalerAddr != 0
}

func (tf typeFeature) Sectioner() bool {
	return tf&tfSectioner != 0
}

func (tf typeFeature) SectionerAddr() bool {
	return tf&tfSectionerAddr != 0
}

func (tf typeFeature) Indenter() bool {
	return tf&tfIndenter != 0
}

func (tf typeFeature) IndenterAddr() bool {
	return tf&tfIndenterAddr != 0
}

func (tf typeFeature) MultiValue() bool {
	return tf&tfMultiValue != 0
}

func (tf typeFeature) MultiValueAddr() bool {
	return tf&tfMultiValueAddr != 0
}

func (tf typeFeature) Block() bool {
	return tf&tfBlock != 0
}

func (tf typeFeature) BlockAddr() bool {
	return tf&tfBlockAddr != 0
}

func (tf typeFeature) Object() bool {
	return tf&tfObject != 0
}

func (tf typeFeature) ObjectAddr() bool {
	return tf&tfObjectAddr != 0
}

func (tf typeFeature) Omit() bool {
	return tf&tfOmit != 0
}

func (tf typeFeature) OmitEmtpy() bool {
	return tf&tfOmitEmpty != 0
}

func (tf typeFeature) Marshaler() bool {
	return tf&tfMarshalerMask != 0
}

func (tf typeFeature) MarshalerAddr() bool {
	return tf&tfMarshalerAddrMask != 0
}
