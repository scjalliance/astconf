package pjsip_test

import (
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestIdentifyMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.Identify{Name: "trunk-den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-den]\ntype = identify\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestIdentifyFull(t *testing.T) {
	identify := pjsip.Identify{
		Name:        "trunk-den",
		Endpoint:    "trunk-den",
		Match:       []string{"den-phone", "100.64.0.0/10"},
		SRVLookups:  astval.No,
		MatchHeader: "X-Trunk: den",
	}
	got, err := astconf.Marshal(&identify)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-den]\n" +
		"type = identify\n" +
		"endpoint = trunk-den\n" +
		"match = den-phone\n" +
		"match = 100.64.0.0/10\n" +
		"srv_lookups = no\n" +
		"match_header = X-Trunk: den\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
