package pjsip_test

import (
	"errors"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/pjsip"
)

func TestInboundPublicationMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.InboundPublication{Name: "trunk-den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-den]\ntype = inbound-publication\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestInboundPublicationFull(t *testing.T) {
	publication := pjsip.InboundPublication{
		Name:     "trunk-hub",
		Endpoint: "trunk-hub",
		Events: pjsip.PublicationEvents{
			{Event: "asterisk-devicestate", Handler: "hub"},
			{Event: "asterisk-mwi", Handler: "hub"},
		},
	}
	got, err := astconf.Marshal(&publication)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[trunk-hub]\n" +
		"type = inbound-publication\n" +
		"endpoint = trunk-hub\n" +
		"event_asterisk-devicestate = hub\n" +
		"event_asterisk-mwi = hub\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

// TestInboundPublicationInvalidEvent checks that an event name that would
// break the setting name is rejected instead of written.
func TestInboundPublicationInvalidEvent(t *testing.T) {
	publication := pjsip.InboundPublication{
		Name:   "trunk-den",
		Events: pjsip.PublicationEvents{{Event: "a=b", Handler: "den"}},
	}
	_, err := astconf.Marshal(&publication)
	var ice astconf.InvalidContentError
	if !errors.As(err, &ice) {
		t.Fatalf("Marshal() error = %v, want InvalidContentError", err)
	}
}
