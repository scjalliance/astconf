package astconf

import (
	"reflect"
	"sort"
)

// Field is an asterisk configuration field. It provides contextual information
// for custom marshalers.
//type Field struct {
//	Name string
//}

type fieldSet []field

func (fs fieldSet) Len() int {
	return len(fs)
}

func (fs fieldSet) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs fieldSet) Less(i, j int) bool {
	a := &fs[i]
	b := &fs[j]

	// Non-blocks before blocks
	if !a.Block() && b.Block() {
		return true
	}
	if a.Block() && !b.Block() {
		return false
	}

	// If both are blocks, non-sectioners before sectioners
	if a.Block() && b.Block() {
		if !a.Sectioner() && b.Sectioner() {
			return true
		}
		if a.Sectioner() && !b.Sectioner() {
			return false
		}
	}

	return false
}

type field struct {
	index    []int
	name     string
	embedded bool
	typeFeature
}

// TODO: Cache fieldsets for each type
//var tfieldsMap sync.Map // map[reflect.Type]*fieldSet

// typeFields returns a list of struct fields that the encoder/decoder
// should process.
func typeFields(typ reflect.Type, allowAddr bool) (fields fieldSet) {
	if typ.Kind() != reflect.Struct {
		return nil
	}

	for i, n := 0, typ.NumField(); i < n; i++ {
		sf := typ.Field(i)
		tag := parseTag(sf.Tag.Get("astconf"))
		tagFeatures := tagFeatures(&tag)
		features := typeFeatures(sf.Type) | tagFeatures

		name := tag.name
		if name == "" {
			name = sf.Name
		}

		fields = append(fields, field{
			index:       []int{i},
			name:        name,
			embedded:    sf.Anonymous,
			typeFeature: features,
		})

		// Embed fields from anonymous structs
		if sf.Anonymous {
			t := sf.Type
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() == reflect.Struct {
				for _, child := range typeFields(t, allowAddr) {
					child.index = append([]int{i}, child.index...)
					child.typeFeature |= tagFeatures
					fields = append(fields, child)
				}
				continue
			}
		}

		// FIXME: Hide hidden fields

		/*
			finfo := structFieldInfo(typ, &sf)
			if finfo.Hidden() {
				continue
			}
		*/

		// FIXME: Avoid infinite recursion
	}

	// Make sure blocks are put at the end of the field list
	sort.Stable(fields)

	return fields
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

/*
	// Anonymous fields to explore at the current level and the next.
	current := []fieldInfo{}
	next := []fieldInfo{{typ: t}}
	var (
		//count = map[reflect.Type]int{}

		visited = map[reflect.Type]struct{}{}
	)

	func typeFields(typ reflect.Type) []field {
		var fi []fieldInfo

		for len(current) > 0 {
			for _, f := range current {
				if _, ok := visited[f.typ]; ok {
					continue
				}
				visited[f.typ] = struct{}{}

				if f.
				fi = structFields(f.typ, f.index, fi)
			}
			current, next = next, current[:0]
		}

	return nil
}
*/

/*
type fieldInfo struct {
	field
	typ reflect.Type

	anonymous  bool
	unexported bool
	tag        fieldTag
	target     reflect.Type
}

func structFieldInfo(container reflect.Type, sf *reflect.StructField) fieldInfo {
	t := sf.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	omit, name, opts := parseTag(sf.Tag.Get("astconf"))
	if name == "" {
		name = sf.Name
	}

	if ft.Implements(sectionNamerType) {
		if se.section == nil || len(se.section) >= len(f.index) {
			se.section = f.index
		}
	}

	mode := typeMode(t)
	if opts.Contains("object") && mode != modeBlock {
		mode = modeObject
	}

	return fieldInfo{
		field: field{
			name:  name,
			index: sf.Index,
			mode:  mode,
		},
		anonymous:  sf.Anonymous,
		unexported: sf.PkgPath != "",
		omit:       omit,
		options:    opts,
		target:     t,
	}
}

func (f *fieldInfo) Hidden() bool {
	if f.anonymous {
		if f.unexported && f.target.Kind() != reflect.Struct {
			// Ignore embedded fields of unexported non-struct types.
			return true
		}
	} else if f.unexported {
		// Ignore unexported non-embedded fields.
		return true
	}
	if f.omit {
		return true
	}
	return false
}
*/
