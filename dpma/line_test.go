package dpma_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
)

func TestOverlayLines(t *testing.T) {
	base := dpma.Line{
		Transport:           "tcp",
		Context:             "from-internal",
		RegistrationTimeout: astval.NewSeconds(time.Minute),
	}
	override := dpma.Line{
		Name:      "fred",
		Extension: "100",
		Transport: "tls",
	}
	got := dpma.OverlayLines(base, override)
	want := dpma.Line{
		Name:                "fred",
		Extension:           "100",
		Transport:           "tls",
		Context:             "from-internal",
		RegistrationTimeout: astval.NewSeconds(time.Minute),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OverlayLines() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestLineMap(t *testing.T) {
	var m dpma.LineMap

	if m.Contains("fred") {
		t.Error("empty map Contains() = true")
	}
	if _, ok := m.Line("fred"); ok {
		t.Error("empty map Line() ok = true")
	}

	if !m.Add(dpma.Line{Name: "fred", Extension: "100"}) {
		t.Error("first Add() = false")
	}
	if m.Add(dpma.Line{Name: "Fred", Extension: "200"}) {
		t.Error("duplicate Add() = true; name lookup should be case-insensitive")
	}
	if l, ok := m.Line("FRED"); !ok || l.Extension != "100" {
		t.Errorf("Line() = %+v, %v; want extension 100", l, ok)
	}

	m.Overlay(dpma.Line{Name: "fred", Label: "Fred"})
	if l, _ := m.Line("fred"); l.Extension != "100" || l.Label != "Fred" {
		t.Errorf("Overlay() produced %+v", l)
	}

	m.Overlay(dpma.Line{Name: "barney"})

	var names []string
	for _, l := range m.Lines() {
		names = append(names, l.Name)
	}
	want := []string{"fred", "barney"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("Lines() order = %v, want %v", names, want)
	}
}
