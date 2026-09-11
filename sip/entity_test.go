package sip_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/sip"
)

func TestEntityMinimal(t *testing.T) {
	got, err := astconf.Marshal(&sip.Entity{Username: "100"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	// The zero Type is written as peer; everything else is omitted.
	want := "[100]\ntype = peer\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestType(t *testing.T) {
	tests := []struct {
		value     sip.Type
		specified bool
		text      string
	}{
		{value: sip.Default, specified: false, text: "peer"},
		{value: sip.Peer, specified: true, text: "peer"},
		{value: sip.User, specified: true, text: "user"},
		{value: sip.Friend, specified: true, text: "friend"},
		{value: sip.Type(99), specified: true, text: ""},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			if got := tt.value.Specified(); got != tt.specified {
				t.Errorf("Specified() = %v, want %v", got, tt.specified)
			}
			if got := tt.value.String(); got != tt.text {
				t.Errorf("String() = %q, want %q", got, tt.text)
			}
			text, err := tt.value.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error: %v", err)
			}
			if string(text) != tt.text {
				t.Errorf("MarshalText() = %q, want %q", text, tt.text)
			}
		})
	}
}

func TestOverlayEntities(t *testing.T) {
	base := sip.Entity{
		Type:      sip.Friend,
		Templates: []string{"base"},
		Allow:     []string{"ulaw"},
		Context:   "from-internal",
		Host:      "dynamic",
		Variables: []astval.Var{astval.NewVar("A", "1")},
	}
	override := sip.Entity{
		Username:  "100",
		Templates: []string{"phone"},
		Allow:     []string{"g722"},
		Secret:    "s",
		Variables: []astval.Var{astval.NewVar("B", "2")},
	}
	got := sip.OverlayEntities(base, override)
	want := sip.Entity{
		Username:  "100",
		Type:      sip.Friend, // Default type in override does not replace Friend
		Templates: []string{"phone"},
		Allow:     []string{"g722"}, // Vectors are replaced, not combined
		Context:   "from-internal",
		Host:      "dynamic",
		Secret:    "s",
		Variables: []astval.Var{astval.NewVar("B", "2")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OverlayEntities() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestMergeEntities(t *testing.T) {
	base := sip.Entity{
		Type:      sip.Peer,
		Templates: []string{"base"},
		Allow:     []string{"ulaw"},
		Context:   "from-internal",
		Variables: []astval.Var{astval.NewVar("A", "1")},
	}
	override := sip.Entity{
		Username:  "100",
		Type:      sip.Friend,
		Templates: []string{"phone"},
		Allow:     []string{"g722", "ulaw"},
		Variables: []astval.Var{astval.NewVar("B", "2")},
	}
	got := sip.MergeEntities(base, override)
	want := sip.Entity{
		Username:  "100",
		Type:      sip.Friend,
		Templates: []string{"phone", "base"}, // Higher priority values come first
		Allow:     []string{"g722", "ulaw"},  // Vectors are combined and deduplicated
		Context:   "from-internal",
		Variables: []astval.Var{astval.NewVar("B", "2"), astval.NewVar("A", "1")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeEntities() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestEntityMap(t *testing.T) {
	var m sip.EntityMap

	if m.Contains("100") {
		t.Error("empty map Contains() = true")
	}
	if _, ok := m.Entity("100"); ok {
		t.Error("empty map Entity() ok = true")
	}

	if !m.Add(sip.Entity{Username: "100", Context: "a"}) {
		t.Error("first Add() = false")
	}
	if m.Add(sip.Entity{Username: "100", Context: "b"}) {
		t.Error("duplicate Add() = true")
	}
	if !m.Contains("100") {
		t.Error("Contains() = false after Add()")
	}
	if !m.Contains("100") || !m.Contains("100") {
		t.Error("Contains() is not case-insensitive")
	}
	if e, ok := m.Entity("100"); !ok || e.Context != "a" {
		t.Errorf("Entity() = %+v, %v; want context a", e, ok)
	}

	// Usernames are matched case-insensitively
	m.Add(sip.Entity{Username: "Fred"})
	if !m.Contains("fred") || !m.Contains("FRED") {
		t.Error("Contains() is not case-insensitive")
	}

	// Overlay replaces vectors and fills scalars
	m.Overlay(sip.Entity{Username: "100", Allow: []string{"g722"}, Secret: "s"})
	if e, _ := m.Entity("100"); e.Context != "a" || e.Secret != "s" || !reflect.DeepEqual(e.Allow, []string{"g722"}) {
		t.Errorf("Overlay() produced %+v", e)
	}

	// Merge combines vectors
	m.Merge(sip.Entity{Username: "100", Allow: []string{"ulaw"}})
	if e, _ := m.Entity("100"); !reflect.DeepEqual(e.Allow, []string{"ulaw", "g722"}) {
		t.Errorf("Merge() produced allow %v", e.Allow)
	}

	// Overlay and Merge add missing entries
	m.Overlay(sip.Entity{Username: "200"})
	m.Merge(sip.Entity{Username: "300"})

	// Entities preserves insertion order
	var names []string
	for _, e := range m.Entities() {
		names = append(names, e.Username)
	}
	want := []string{"100", "Fred", "200", "300"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("Entities() order = %v, want %v", names, want)
	}
}
