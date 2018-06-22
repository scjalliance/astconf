package astconf

import (
	"reflect"
	"sort"
)

type fieldSet []field

func (fs fieldSet) Len() int {
	return len(fs)
}

func (fs fieldSet) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs fieldSet) Less(i, j int) bool {
	// Any field that could possibly start its own section should always go
	// at the end. This helps to avoid accidental mixing of fields that
	// belong to different sections.
	//
	// Right now we do this in a hacky way by assuming that every "block" is
	// a potential section starter. The term "block" here is nebulous, but
	// right now blocks are things that could print multiple lines, including
	// structs, slices, and types with custom marshalers.
	//
	// Probably the more correct thing to do is to accurately determine
	// whether a type really could be a section starter by examining it and
	// all its descedent elements. A type with any of these qualities should
	// be considered:
	//  1. The type itself is a Sectioner
	//  2. The type is a struct that contains an embedded sectioner
	//  3. The type itself has a MarshalAsterisk() function
	//  4. The type is a struct that ontains an embedded MarshalAsterisk() function
	//  5. The type is a struct with any non-embedded member that meets any of 1-4
	//
	// Such an analysis should probably be perform by typeFeatures().

	a := &fs[i]
	b := &fs[j]

	// Note: We skip multi-value slices because their elements aren't blocks
	aIsBlock := a.Block() && !a.MultiValue()
	bIsBlock := b.Block() && !b.MultiValue()

	// Non-blocks before blocks
	if !aIsBlock && bIsBlock {
		return true
	}
	if aIsBlock && !bIsBlock {
		return false
	}

	// If both are blocks, non-sectioners before sectioners
	if aIsBlock && bIsBlock {
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
