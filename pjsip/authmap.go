package pjsip

import "strings"

// AuthMap is an ordered map of auths, indexed by auth name.
//
// It should not be copied by value.
//
// An empty auth map is ready for use.
type AuthMap struct {
	auths  []Auth
	lookup map[string]int // Maps auth names to indices in the auths slice
}

// Contains returns true if the map contains an auth with the given name.
func (m *AuthMap) Contains(name string) bool {
	if m.lookup == nil {
		return false
	}

	name = strings.ToLower(name)

	_, exists := m.lookup[name]
	return exists
}

// Add adds an auth to the map. If the map already contains an entry with
// the auth's name, the map is not changed.
//
// Add returns true if the auth was added to the map.
func (m *AuthMap) Add(auth Auth) bool {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(auth.Name)

	if _, exists := m.lookup[name]; exists {
		return false
	}

	m.auths = append(m.auths, auth)
	if name != "" {
		index := len(m.auths) - 1
		m.lookup[name] = index
	}

	return true
}

// Overlay adds an auth to the map. If the map already contains an entry
// with the auth's name, the entries are overlayed, with priority given
// to the new entry.
func (m *AuthMap) Overlay(auth Auth) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(auth.Name)

	index, exists := m.lookup[name]
	if exists {
		m.auths[index] = OverlayAuths(m.auths[index], auth)
		return
	}

	m.auths = append(m.auths, auth)
	index = len(m.auths) - 1
	if name != "" {
		m.lookup[name] = index
	}
}

// Auth returns the auth with the given name.
func (m *AuthMap) Auth(name string) (auth Auth, ok bool) {
	if m.lookup == nil {
		return
	}

	name = strings.ToLower(name)

	index, ok := m.lookup[name]
	if !ok {
		return
	}
	return m.auths[index], true
}

// Auths returns a slice of all auths in the map.
func (m *AuthMap) Auths() []Auth {
	return m.auths
}

// Merge adds an auth to the map. If the map already contains an entry
// with the auth's name, the entries are merged.
func (m *AuthMap) Merge(auth Auth) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(auth.Name)

	index, exists := m.lookup[name]
	if exists {
		m.auths[index] = MergeAuths(m.auths[index], auth)
		return
	}

	m.auths = append(m.auths, auth)
	index = len(m.auths) - 1
	if name != "" {
		m.lookup[name] = index
	}
}
