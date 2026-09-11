package pjsip_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/pjsip"
)

func TestAuthMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.Auth{Name: "100-auth"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[100-auth]\ntype = auth\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestOverlayAuths(t *testing.T) {
	base := pjsip.Auth{AuthType: "userpass", Realm: "asterisk"}
	override := pjsip.Auth{Name: "100-auth", Username: "100", Password: "secret"}
	got := pjsip.OverlayAuths(base, override)
	want := pjsip.Auth{
		Name:     "100-auth",
		AuthType: "userpass",
		Username: "100",
		Password: "secret",
		Realm:    "asterisk",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OverlayAuths() =\n%+v\nwant:\n%+v", got, want)
	}
}

func TestAuthMap(t *testing.T) {
	var m pjsip.AuthMap

	if m.Contains("100-auth") {
		t.Error("empty map Contains() = true")
	}
	if _, ok := m.Auth("100-auth"); ok {
		t.Error("empty map Auth() ok = true")
	}

	if !m.Add(pjsip.Auth{Name: "100-auth", Username: "100"}) {
		t.Error("first Add() = false")
	}
	if m.Add(pjsip.Auth{Name: "100-AUTH", Username: "x"}) {
		t.Error("duplicate Add() = true; name lookup should be case-insensitive")
	}
	if a, ok := m.Auth("100-auth"); !ok || a.Username != "100" {
		t.Errorf("Auth() = %+v, %v; want username 100", a, ok)
	}

	m.Overlay(pjsip.Auth{Name: "100-auth", Password: "secret"})
	if a, _ := m.Auth("100-auth"); a.Username != "100" || a.Password != "secret" {
		t.Errorf("Overlay() produced %+v", a)
	}

	m.Overlay(pjsip.Auth{Name: "200-auth"})

	var names []string
	for _, a := range m.Auths() {
		names = append(names, a.Name)
	}
	want := []string{"100-auth", "200-auth"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("Auths() order = %v, want %v", names, want)
	}
}
