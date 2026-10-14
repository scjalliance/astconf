package dialplan_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/dialplan"
)

// TestHint renders an extension's hint before its first priority, and a
// hint on an extension with no actions on its own.
func TestHint(t *testing.T) {
	section := dialplan.Section{
		Context: "test",
		Extensions: []dialplan.Extension{
			{
				Comment: "Fred",
				Number:  "100",
				Hint:    dialplan.Hint{dialplan.PJSIP("fred"), dialplan.PJSIP("fred.soft")},
				Actions: []dialplan.Action{dialplan.Noop("Call fred")},
			},
			{
				Number: "101",
				Hint:   dialplan.Hint{dialplan.PJSIP("barney")},
			},
			{
				Number:  "102",
				Actions: []dialplan.Action{dialplan.Noop("No hint")},
			},
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&section); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[test]\n" +
		"\n" +
		"; Fred\n" +
		"exten => 100,hint,PJSIP/fred&PJSIP/fred.soft\n" +
		"exten => 100,1,Noop(Call fred)\n" +
		"\n" +
		"exten => 101,hint,PJSIP/barney\n" +
		"\n" +
		"exten => 102,1,Noop(No hint)\n"
	if got := buf.String(); got != want {
		t.Errorf("Encode() =\n%s\nwant:\n%s", got, want)
	}
}

// TestHintInvalid checks that a hint device that would break the line is
// rejected.
func TestHintInvalid(t *testing.T) {
	for _, resource := range []string{"a&b", "a,b", "a/b"} {
		t.Run(resource, func(t *testing.T) {
			section := dialplan.Section{Context: "test", Extensions: []dialplan.Extension{{Number: "1", Hint: dialplan.Hint{dialplan.PJSIP(resource)}}}}
			var buf bytes.Buffer
			err := astconf.NewEncoder(&buf).Encode(&section)
			var ice astconf.InvalidContentError
			if !errors.As(err, &ice) {
				t.Fatalf("Encode() error = %v, want InvalidContentError", err)
			}
		})
	}
}
