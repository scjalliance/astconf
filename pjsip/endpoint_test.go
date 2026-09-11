package pjsip_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestEndpointMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.Endpoint{Name: "100"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[100]\ntype = endpoint\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestOverlayEndpoints(t *testing.T) {
	base := pjsip.Endpoint{
		Templates:    []string{"base"},
		Transport:    "transport-udp",
		Context:      "from-internal",
		Allow:        []string{"ulaw"},
		RTPSymmetric: astval.Yes,
		Variables:    []astval.Var{astval.NewVar("A", "1")},
	}
	override := pjsip.Endpoint{
		Name:      "100",
		Templates: []string{"phone"},
		Transport: "transport-tcp",
		Allow:     []string{"g722"},
		Auth:      "100-auth",
		AORs:      []string{"100-aor"},
		Variables: []astval.Var{astval.NewVar("B", "2")},
	}
	got := pjsip.OverlayEndpoints(base, override)
	want := pjsip.Endpoint{
		Name:         "100",
		Templates:    []string{"phone"},
		Transport:    "transport-tcp",
		Context:      "from-internal",
		Allow:        []string{"g722"}, // Vectors are replaced, not combined
		Auth:         "100-auth",
		AORs:         []string{"100-aor"},
		RTPSymmetric: astval.Yes,
		Variables:    []astval.Var{astval.NewVar("B", "2")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OverlayEndpoints() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestMergeEndpoints(t *testing.T) {
	base := pjsip.Endpoint{
		Templates: []string{"base"},
		Context:   "from-internal",
		Disallow:  []string{"all"},
		Allow:     []string{"ulaw"},
		Variables: []astval.Var{astval.NewVar("A", "1")},
	}
	override := pjsip.Endpoint{
		Name:      "100",
		Templates: []string{"phone"},
		Allow:     []string{"g722", "ulaw"},
		Mailboxes: []string{"100@default"},
		Variables: []astval.Var{astval.NewVar("B", "2")},
	}
	got := pjsip.MergeEndpoints(base, override)
	want := pjsip.Endpoint{
		Name:      "100",
		Templates: []string{"phone", "base"}, // Higher priority values come first
		Context:   "from-internal",
		Disallow:  []string{"all"},
		Allow:     []string{"g722", "ulaw"}, // Vectors are combined and deduplicated
		Mailboxes: []string{"100@default"},
		Variables: []astval.Var{astval.NewVar("B", "2"), astval.NewVar("A", "1")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeEndpoints() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestEndpointMap(t *testing.T) {
	var m pjsip.EndpointMap

	if m.Contains("100") {
		t.Error("empty map Contains() = true")
	}
	if _, ok := m.Endpoint("100"); ok {
		t.Error("empty map Endpoint() ok = true")
	}

	if !m.Add(pjsip.Endpoint{Name: "100", Context: "a"}) {
		t.Error("first Add() = false")
	}
	if m.Add(pjsip.Endpoint{Name: "100", Context: "b"}) {
		t.Error("duplicate Add() = true")
	}
	if e, ok := m.Endpoint("100"); !ok || e.Context != "a" {
		t.Errorf("Endpoint() = %+v, %v; want context a", e, ok)
	}

	// Names are matched case-insensitively
	m.Add(pjsip.Endpoint{Name: "Fred"})
	if !m.Contains("fred") || !m.Contains("FRED") {
		t.Error("Contains() is not case-insensitive")
	}

	m.Overlay(pjsip.Endpoint{Name: "100", Allow: []string{"g722"}, Auth: "100-auth"})
	if e, _ := m.Endpoint("100"); e.Context != "a" || e.Auth != "100-auth" || !reflect.DeepEqual(e.Allow, []string{"g722"}) {
		t.Errorf("Overlay() produced %+v", e)
	}

	m.Merge(pjsip.Endpoint{Name: "100", Allow: []string{"ulaw"}})
	if e, _ := m.Endpoint("100"); !reflect.DeepEqual(e.Allow, []string{"ulaw", "g722"}) {
		t.Errorf("Merge() produced allow %v", e.Allow)
	}

	m.Overlay(pjsip.Endpoint{Name: "200"})
	m.Merge(pjsip.Endpoint{Name: "300"})

	var names []string
	for _, e := range m.Endpoints() {
		names = append(names, e.Name)
	}
	want := []string{"100", "Fred", "200", "300"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("Endpoints() order = %v, want %v", names, want)
	}
}
