package dpma

import (
	"strings"
)

// PhoneMap is an ordered map of dpma phone entries, indexed by mac address.
//
// It should not be copied by value.
//
// An empty phone map is ready for use.
type PhoneMap struct {
	phones []Phone
	lookup map[string]int // Maps MAC addresses to indices in the phones slice
}

// Contains returns true if the map contains a phone with the given mac
// address.
func (m *PhoneMap) Contains(mac string) bool {
	if m.lookup == nil {
		return false
	}

	_, exists := m.lookup[strings.ToLower(mac)]
	return exists
}

// Add adds a phone to the map. If the map already contains an entry with
// the phone's mac address, the map is not changed.
//
// Add returns true if the phone was added to the map.
func (m *PhoneMap) Add(phone Phone) bool {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	mac := strings.ToLower(phone.MAC)

	if _, exists := m.lookup[mac]; exists {
		return false
	}

	m.phones = append(m.phones, phone)
	if mac != "" {
		index := len(m.phones) - 1
		m.lookup[mac] = index
	}

	return true
}

// Overlay adds a phone to the map. If the map already contains an entry with
// the phone's mac address, the entries are overlayed, with priority given
// to the new entry.
func (m *PhoneMap) Overlay(phone Phone) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	mac := strings.ToLower(phone.MAC)

	index, exists := m.lookup[mac]
	if exists {
		m.phones[index] = OverlayPhones(m.phones[index], phone)
		return
	}

	m.phones = append(m.phones, phone)
	index = len(m.phones) - 1
	if mac != "" {
		m.lookup[mac] = index
	}
}

// Merge adds a phone to the map. If the map already contains an entry with
// the phone's mac address, the entries are merged.
func (m *PhoneMap) Merge(phone Phone) {
	if m.lookup == nil {
		m.lookup = make(map[string]int)
	}

	mac := strings.ToLower(phone.MAC)

	index, exists := m.lookup[mac]
	if exists {
		m.phones[index] = MergePhones(m.phones[index], phone)
		return
	}

	m.phones = append(m.phones, phone)
	index = len(m.phones) - 1
	if mac != "" {
		m.lookup[mac] = index
	}
}

// Phone returns the phone entry with the given mac address.
func (m *PhoneMap) Phone(mac string) (phone Phone, ok bool) {
	if m.lookup == nil {
		return
	}

	mac = strings.ToLower(mac)

	index, ok := m.lookup[mac]
	if !ok {
		return
	}
	return m.phones[index], true
}

// Phones returns a slice of all phone entries in the map.
func (m *PhoneMap) Phones() []Phone {
	return m.phones
}
