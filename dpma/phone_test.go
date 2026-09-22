package dpma_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
)

func TestPhoneSectionName(t *testing.T) {
	tests := []struct {
		name  string
		phone dpma.Phone
		want  string
	}{
		{name: "username", phone: dpma.Phone{Username: "fred", MAC: "00:0F:D3:00:00:01"}, want: "fred"},
		{name: "mac without username", phone: dpma.Phone{MAC: "000FD3000001"}, want: "000FD3000001"},
		{name: "colons removed from mac", phone: dpma.Phone{MAC: "00:0F:D3:00:00:01"}, want: "000FD3000001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.phone.SectionName(); got != tt.want {
				t.Errorf("SectionName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOverlayPhones(t *testing.T) {
	base := dpma.Phone{
		Networks:     []string{"a"},
		Lines:        []string{"line1"},
		Timezone:     "UTC",
		WebUIEnabled: astval.No,
		Brightness:   astval.NewInt(3),
	}
	override := dpma.Phone{
		Username: "fred",
		Lines:    []string{"line2"},
		Timezone: "America/Los_Angeles",
	}
	got := dpma.OverlayPhones(base, override)
	want := dpma.Phone{
		Username:     "fred",
		Networks:     []string{"a"},
		Lines:        []string{"line2"},
		Timezone:     "America/Los_Angeles",
		WebUIEnabled: astval.No,
		Brightness:   astval.NewInt(3),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OverlayPhones() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestMergePhones(t *testing.T) {
	base := dpma.Phone{
		Networks: []string{"a"},
		Firmware: []string{"old"},
		Lines:    []string{"line1"},
		Alerts:   []string{"page"},
	}
	override := dpma.Phone{
		Username: "fred",
		Networks: []string{"b", "a"},
		Firmware: []string{"new"},
		Lines:    []string{"line2"},
	}
	got := dpma.MergePhones(base, override)
	want := dpma.Phone{
		Username: "fred",
		Networks: []string{"b", "a"}, // Higher priority values come first
		Firmware: []string{"new"},    // Firmware is overlayed, not merged
		Lines:    []string{"line2", "line1"},
		Alerts:   []string{"page"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergePhones() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestPhoneMap(t *testing.T) {
	var m dpma.PhoneMap

	if m.Contains("000FD3000001") {
		t.Error("empty map Contains() = true")
	}
	if _, ok := m.Phone("000FD3000001"); ok {
		t.Error("empty map Phone() ok = true")
	}

	if !m.Add(dpma.Phone{MAC: "000FD3000001", Timezone: "a"}) {
		t.Error("first Add() = false")
	}
	if m.Add(dpma.Phone{MAC: "000fd3000001", Timezone: "b"}) {
		t.Error("duplicate Add() = true; mac lookup should be case-insensitive")
	}
	if p, ok := m.Phone("000fd3000001"); !ok || p.Timezone != "a" {
		t.Errorf("Phone() = %+v, %v; want timezone a", p, ok)
	}

	m.Overlay(dpma.Phone{MAC: "000FD3000001", Lines: []string{"x"}, FullName: "n"})
	if p, _ := m.Phone("000FD3000001"); p.Timezone != "a" || p.FullName != "n" || !reflect.DeepEqual(p.Lines, []string{"x"}) {
		t.Errorf("Overlay() produced %+v", p)
	}

	m.Merge(dpma.Phone{MAC: "000FD3000001", Lines: []string{"y"}})
	if p, _ := m.Phone("000FD3000001"); !reflect.DeepEqual(p.Lines, []string{"y", "x"}) {
		t.Errorf("Merge() produced lines %v", p.Lines)
	}

	m.Overlay(dpma.Phone{MAC: "000FD3000002"})
	m.Merge(dpma.Phone{MAC: "000FD3000003"})

	var macs []string
	for _, p := range m.Phones() {
		macs = append(macs, p.MAC)
	}
	want := []string{"000FD3000001", "000FD3000002", "000FD3000003"}
	if !reflect.DeepEqual(macs, want) {
		t.Errorf("Phones() order = %v, want %v", macs, want)
	}
}
