package pjsip_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestAORMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.AOR{Name: "100-aor"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[100-aor]\ntype = aor\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestOverlayAORs(t *testing.T) {
	base := pjsip.AOR{MaxContacts: astval.NewInt(5), RemoveExisting: astval.Yes, Contacts: []string{"a"}}
	override := pjsip.AOR{Name: "100-aor", MaxContacts: astval.NewInt(1), Contacts: []string{"b"}}
	got := pjsip.OverlayAORs(base, override)
	want := pjsip.AOR{
		Name:           "100-aor",
		MaxContacts:    astval.NewInt(1),
		RemoveExisting: astval.Yes,
		Contacts:       []string{"b"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OverlayAORs() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestMergeAORs(t *testing.T) {
	base := pjsip.AOR{MaxContacts: astval.NewInt(5), Contacts: []string{"a"}}
	override := pjsip.AOR{Name: "100-aor", Contacts: []string{"b", "a"}, Mailboxes: []string{"100@default"}}
	got := pjsip.MergeAORs(base, override)
	want := pjsip.AOR{
		Name:        "100-aor",
		MaxContacts: astval.NewInt(5),
		Contacts:    []string{"b", "a"}, // Higher priority values come first
		Mailboxes:   []string{"100@default"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeAORs() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestAORMap(t *testing.T) {
	var m pjsip.AORMap

	if m.Contains("100-aor") {
		t.Error("empty map Contains() = true")
	}
	if _, ok := m.AOR("100-aor"); ok {
		t.Error("empty map AOR() ok = true")
	}

	if !m.Add(pjsip.AOR{Name: "100-aor", MaxContacts: astval.NewInt(1)}) {
		t.Error("first Add() = false")
	}
	if m.Add(pjsip.AOR{Name: "100-AOR", MaxContacts: astval.NewInt(9)}) {
		t.Error("duplicate Add() = true; name lookup should be case-insensitive")
	}
	if a, ok := m.AOR("100-aor"); !ok || a.MaxContacts.Value() != 1 {
		t.Errorf("AOR() = %+v, %v; want max_contacts 1", a, ok)
	}

	m.Overlay(pjsip.AOR{Name: "100-aor", Contacts: []string{"x"}, RemoveExisting: astval.Yes})
	if a, _ := m.AOR("100-aor"); a.MaxContacts.Value() != 1 || !a.RemoveExisting.True() || !reflect.DeepEqual(a.Contacts, []string{"x"}) {
		t.Errorf("Overlay() produced %+v", a)
	}

	m.Merge(pjsip.AOR{Name: "100-aor", Contacts: []string{"y"}})
	if a, _ := m.AOR("100-aor"); !reflect.DeepEqual(a.Contacts, []string{"y", "x"}) {
		t.Errorf("Merge() produced contacts %v", a.Contacts)
	}

	m.Overlay(pjsip.AOR{Name: "200-aor"})
	m.Merge(pjsip.AOR{Name: "300-aor"})

	var names []string
	for _, a := range m.AORs() {
		names = append(names, a.Name)
	}
	want := []string{"100-aor", "200-aor", "300-aor"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("AORs() order = %v, want %v", names, want)
	}
}
