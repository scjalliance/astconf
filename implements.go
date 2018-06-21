package astconf

import "reflect"

func typeImplements(t reflect.Type, interfaces ...reflect.Type) bool {
	for _, u := range interfaces {
		if t.Implements(u) {
			return true
		}
	}
	return false
}

func typeAddrImplements(t reflect.Type, interfaces ...reflect.Type) bool {
	// If the type is already a pointer don't take the address of it
	if t.Kind() == reflect.Ptr {
		return false
	}
	t = reflect.PtrTo(t)
	return typeImplements(t, interfaces...)
}
