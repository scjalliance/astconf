package astgen_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/dpma"
	"github.com/scjalliance/astconf/pjsip"
)

// TestServerDeviceSections checks which devices each server renders
// endpoints, aors, lines and phones for. The rendered text of each entry is
// covered by the golden files.
func TestServerDeviceSections(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	tests := []struct {
		server    string
		endpoints []string
		aors      []string
		lines     []string
		phones    []string
	}{
		{
			server:    "den",
			endpoints: []string{"lobby", "barney", "000fd3000009"},
			aors:      []string{"lobby-aor", "barney-aor", "000fd3000009-aor"},
			lines:     []string{"lobby", "barney", "000fd3000009"},
			phones:    []string{"000fd3000002", "000fd3000004", "000fd3000009"},
		},
		{
			// No location names the hub, so it renders no spare.
			server:    "hub",
			endpoints: []string{"fred.soft", "wilma.soft", "yard.soft"},
			aors:      []string{"fred.soft-aor", "wilma.soft-aor", "yard.soft-aor"},
		},
		{
			server:    "qry",
			endpoints: []string{"fred", "000fd3000003", "000fd3000009"},
			aors:      []string{"fred-aor", "000fd3000003-aor", "000fd3000009-aor"},
			lines:     []string{"fred", "000fd3000003", "000fd3000009"},
			phones:    []string{"000fd3000001", "000fd3000003", "000fd3000009"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.server, func(t *testing.T) {
			var endpoints, aors, lines, phones []string
			for _, e := range astgen.ServerEndpoints(data, p, tt.server, pjsip.Endpoint{}, "default") {
				endpoints = append(endpoints, e.Name)
			}
			for _, a := range astgen.ServerAORs(data, p, tt.server, pjsip.AOR{}) {
				aors = append(aors, a.Name)
			}
			for _, l := range astgen.ServerLines(data, p, tt.server) {
				lines = append(lines, l.Name)
			}
			for _, ph := range astgen.ServerPhones(data, p, tt.server, dpma.Phone{}, "") {
				phones = append(phones, ph.MAC)
			}
			if !reflect.DeepEqual(endpoints, tt.endpoints) {
				t.Errorf("ServerEndpoints() = %v, want %v", endpoints, tt.endpoints)
			}
			if !reflect.DeepEqual(aors, tt.aors) {
				t.Errorf("ServerAORs() = %v, want %v", aors, tt.aors)
			}
			if !reflect.DeepEqual(lines, tt.lines) {
				t.Errorf("ServerLines() = %v, want %v", lines, tt.lines)
			}
			if !reflect.DeepEqual(phones, tt.phones) {
				t.Errorf("ServerPhones() = %v, want %v", phones, tt.phones)
			}
		})
	}
}

// TestServerEndpointsNoVoicemailCode checks that no server renders a
// voicemail PIN into an endpoint, while the existing generator still does.
func TestServerEndpointsNoVoicemailCode(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	existing := astgen.Endpoints(data, pjsip.Endpoint{}, "default")
	var existingValues []interface{}
	for i := range existing {
		existingValues = append(existingValues, &existing[i])
	}
	if got := render(t, existingValues...); !strings.Contains(got, "VMCODE") {
		t.Fatalf("Endpoints() has no VMCODE, so this test proves nothing:\n%s", got)
	}

	for _, server := range []string{"den", "hub", "qry"} {
		endpoints := astgen.ServerEndpoints(data, p, server, pjsip.Endpoint{}, "default")
		var values []interface{}
		for i := range endpoints {
			values = append(values, &endpoints[i])
		}
		got := render(t, values...)
		if strings.Contains(got, "VMCODE") || strings.Contains(got, "1234") || strings.Contains(got, "4321") {
			t.Errorf("ServerEndpoints(%s) holds a voicemail code:\n%s", server, got)
		}
	}
}

// TestServerMailboxes checks that voicemail renders on the hub alone, and
// that no mailbox carries a password.
func TestServerMailboxes(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)

	office := astgen.ServerMailboxes(data, p, "qry", "default")
	if got, want := render(t, &office), "[default]\n"; got != want {
		t.Errorf("ServerMailboxes(qry) rendered:\n%s\nwant:\n%s", got, want)
	}

	hub := astgen.ServerMailboxes(data, p, "hub", "default")
	got := render(t, &hub)
	want := "[default]\n" +
		"300 => -x,Sales,sales@example.com,,attach=yes|saycid=yes|delete=yes\n" +
		"100 => -x,fred,fred@example.com,,attach=yes|saycid=yes|delete=no\n" +
		"101 => -x,barney,,,attach=yes|saycid=yes|delete=no\n"
	if got != want {
		t.Errorf("ServerMailboxes(hub) rendered:\n%s\nwant:\n%s", got, want)
	}
	if len(hub.Mailboxes) != 3 {
		t.Fatalf("ServerMailboxes(hub) has %d mailboxes, want 3", len(hub.Mailboxes))
	}
	for _, box := range hub.Mailboxes {
		if box.Password != astgen.LockedMailboxPassword {
			t.Errorf("mailbox %s has password %q, want the locked placeholder", box.Extension, box.Password)
		}
	}
	if strings.Contains(got, "3000") || strings.Contains(got, "1234") || strings.Contains(got, "4321") {
		t.Errorf("ServerMailboxes(hub) holds a PIN:\n%s", got)
	}
}
