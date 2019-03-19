package dpma

import (
	"strings"
)

// LineMap is an ordered map of dpma line entries, indexed by line name.
//
// It should not be copied by value.
//
// An empty line map is ready for use.
type LineMap struct {
	lines  []Line
	lookup map[string]int // Maps line names to indices in the lines slice
}

// Contains returns true if the map contains a line with the given name.
func (m *LineMap) Contains(name string) bool {
	if m.lookup == nil {
		return false
	}

	name = strings.ToLower(name)

	_, exists := m.lookup[name]
	return exists
}

// Add adds a line to the map. If the map already contains an entry with
// the line's name, the map is not changed.
//
// Add returns true if the line was added to the map.
func (m *LineMap) Add(line Line) bool {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(line.Name)

	if _, exists := m.lookup[name]; exists {
		return false
	}

	m.lines = append(m.lines, line)

	if name != "" {
		index := len(m.lines) - 1
		m.lookup[name] = index
	}

	return true
}

// Overlay adds a line to the map. If the map already contains an entry with
// the line's name, the entries are overlayed, with priority given
// to the new entry.
func (m *LineMap) Overlay(line Line) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(line.Name)

	index, exists := m.lookup[name]
	if exists {
		m.lines[index] = OverlayLines(m.lines[index], line)
		return
	}

	m.lines = append(m.lines, line)
	index = len(m.lines) - 1
	if name != "" {
		m.lookup[name] = index
	}
}

// Line returns the line entry with the given name.
func (m *LineMap) Line(name string) (line Line, ok bool) {
	if m.lookup == nil {
		return
	}

	name = strings.ToLower(name)

	index, ok := m.lookup[name]
	if !ok {
		return
	}
	return m.lines[index], true
}

// Lines returns a slice of all line entries in the map.
func (m *LineMap) Lines() []Line {
	return m.lines
}
