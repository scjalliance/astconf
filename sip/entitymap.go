package sip

import "strings"

// EntityMap is an ordered map of entities, indexed by entity username.
//
// It should not be copied by value.
//
// An empty entity map is ready for use.
type EntityMap struct {
	entities []Entity
	lookup   map[string]int // Maps entity usernames to indices in the entities slice
}

// Contains returns true if the map contains an entity with the given username.
func (m *EntityMap) Contains(username string) bool {
	if m.lookup == nil {
		return false
	}

	username = strings.ToLower(username)

	_, exists := m.lookup[username]
	return exists
}

// Add adds an entity to the map. If the map already contains an entry with
// the entity's username, the map is not changed.
//
// Add returns true if the entity was added to the map.
func (m *EntityMap) Add(entity Entity) bool {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	username := strings.ToLower(entity.Username)

	if _, exists := m.lookup[username]; exists {
		return false
	}

	m.entities = append(m.entities, entity)
	if username != "" {
		index := len(m.entities) - 1
		m.lookup[username] = index
	}

	return true
}

// Overlay adds an entity to the map. If the map already contains an entry with
// the entity's username, the entries are overlayed, with priority given
// to the new entry.
func (m *EntityMap) Overlay(entity Entity) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	username := strings.ToLower(entity.Username)

	index, exists := m.lookup[username]
	if exists {
		m.entities[index] = OverlayEntities(m.entities[index], entity)
		return
	}

	m.entities = append(m.entities, entity)
	index = len(m.entities) - 1
	if username != "" {
		m.lookup[username] = index
	}
}

// Merge adds an entity to the map. If the map already contains an entry with
// the entity's username, the entries are merged.
func (m *EntityMap) Merge(entity Entity) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	username := strings.ToLower(entity.Username)

	index, exists := m.lookup[username]
	if exists {
		m.entities[index] = MergeEntities(m.entities[index], entity)
		return
	}

	m.entities = append(m.entities, entity)
	index = len(m.entities) - 1
	if username != "" {
		m.lookup[username] = index
	}
}

// Entity returns the entity with the given username.
func (m *EntityMap) Entity(username string) (entity Entity, ok bool) {
	if m.lookup == nil {
		return
	}

	username = strings.ToLower(username)

	index, ok := m.lookup[username]
	if !ok {
		return
	}
	return m.entities[index], true
}

// Entities returns a slice of all entities in the map.
func (m *EntityMap) Entities() []Entity {
	return m.entities
}
