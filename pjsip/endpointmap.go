package pjsip

import "strings"

// EndpointMap is an ordered map of endpoints, indexed by endpoint name.
//
// It should not be copied by value.
//
// An empty endpoint map is ready for use.
type EndpointMap struct {
	endpoints []Endpoint
	lookup    map[string]int // Maps endpoint names to indices in the endpoints slice
}

// Contains returns true if the map contains an endpoint with the given name.
func (m *EndpointMap) Contains(name string) bool {
	if m.lookup == nil {
		return false
	}

	name = strings.ToLower(name)

	_, exists := m.lookup[name]
	return exists
}

// Add adds an endpoint to the map. If the map already contains an entry with
// the endpoint's name, the map is not changed.
//
// Add returns true if the endpoint was added to the map.
func (m *EndpointMap) Add(endpoint Endpoint) bool {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(endpoint.Name)

	if _, exists := m.lookup[name]; exists {
		return false
	}

	m.endpoints = append(m.endpoints, endpoint)
	if name != "" {
		index := len(m.endpoints) - 1
		m.lookup[name] = index
	}

	return true
}

// Overlay adds an endpoint to the map. If the map already contains an entry
// with the endpoint's name, the entries are overlayed, with priority given
// to the new entry.
func (m *EndpointMap) Overlay(endpoint Endpoint) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(endpoint.Name)

	index, exists := m.lookup[name]
	if exists {
		m.endpoints[index] = OverlayEndpoints(m.endpoints[index], endpoint)
		return
	}

	m.endpoints = append(m.endpoints, endpoint)
	index = len(m.endpoints) - 1
	if name != "" {
		m.lookup[name] = index
	}
}

// Merge adds an endpoint to the map. If the map already contains an entry
// with the endpoint's name, the entries are merged.
func (m *EndpointMap) Merge(endpoint Endpoint) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	name := strings.ToLower(endpoint.Name)

	index, exists := m.lookup[name]
	if exists {
		m.endpoints[index] = MergeEndpoints(m.endpoints[index], endpoint)
		return
	}

	m.endpoints = append(m.endpoints, endpoint)
	index = len(m.endpoints) - 1
	if name != "" {
		m.lookup[name] = index
	}
}

// Endpoint returns the endpoint with the given name.
func (m *EndpointMap) Endpoint(name string) (endpoint Endpoint, ok bool) {
	if m.lookup == nil {
		return
	}

	name = strings.ToLower(name)

	index, ok := m.lookup[name]
	if !ok {
		return
	}
	return m.endpoints[index], true
}

// Endpoints returns a slice of all endpoints in the map.
func (m *EndpointMap) Endpoints() []Endpoint {
	return m.endpoints
}
