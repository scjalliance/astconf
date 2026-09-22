package pjsip

import "strings"

// AORMap is an ordered map of aors, indexed by aor name.
//
// It should not be copied by value.
//
// An empty aor map is ready for use.
type AORMap struct {
	aors   []AOR
	lookup map[string]int // Maps aor names to indices in the aors slice
}

// Contains returns true if the map contains an aor with the given name.
func (m *AORMap) Contains(name string) bool {
	if m.lookup == nil {
		return false
	}

	name = strings.ToLower(name)

	_, exists := m.lookup[name]
	return exists
}

// Add adds an aor to the map. If the map already contains an entry with
// the aor's name, the map is not changed.
//
// Add returns true if the aor was added to the map.
func (m *AORMap) Add(aor AOR) bool {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(aor.Name)

	if _, exists := m.lookup[name]; exists {
		return false
	}

	m.aors = append(m.aors, aor)
	if name != "" {
		index := len(m.aors) - 1
		m.lookup[name] = index
	}

	return true
}

// Overlay adds an aor to the map. If the map already contains an entry
// with the aor's name, the entries are overlayed, with priority given
// to the new entry.
func (m *AORMap) Overlay(aor AOR) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(aor.Name)

	index, exists := m.lookup[name]
	if exists {
		m.aors[index] = OverlayAORs(m.aors[index], aor)
		return
	}

	m.aors = append(m.aors, aor)
	index = len(m.aors) - 1
	if name != "" {
		m.lookup[name] = index
	}
}

// Merge adds an aor to the map. If the map already contains an entry
// with the aor's name, the entries are merged.
func (m *AORMap) Merge(aor AOR) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(aor.Name)

	index, exists := m.lookup[name]
	if exists {
		m.aors[index] = MergeAORs(m.aors[index], aor)
		return
	}

	m.aors = append(m.aors, aor)
	index = len(m.aors) - 1
	if name != "" {
		m.lookup[name] = index
	}
}

// AOR returns the aor with the given name.
func (m *AORMap) AOR(name string) (aor AOR, ok bool) {
	if m.lookup == nil {
		return
	}

	name = strings.ToLower(name)

	index, ok := m.lookup[name]
	if !ok {
		return
	}
	return m.aors[index], true
}

// AORs returns a slice of all aors in the map.
func (m *AORMap) AORs() []AOR {
	return m.aors
}
