package astgen_test

import (
	"strings"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/pjsip"
	"github.com/scjalliance/astconf/sip"
)

// TestCallerIDQuotes checks that double quotes in a person's or role's
// display name are removed, because they would end the quoted name part
// of the callerid value early.
func TestCallerIDQuotes(t *testing.T) {
	data := testData()
	data.People[0].FullName = "Fred \"Freddy\" Flintstone"
	data.PhoneRoles[0].DisplayName = "Lobby \"Front\""

	entities := astgen.SIP(data, sip.Entity{}, "default")
	endpoints := astgen.Endpoints(data, pjsip.Endpoint{}, "default")

	var got []string
	for _, e := range entities {
		got = append(got, e.CallerID)
	}
	for _, e := range endpoints {
		got = append(got, e.CallerID)
	}
	for _, callerID := range got {
		if strings.Count(callerID, "\"") > 2 {
			t.Errorf("callerid %q has more than one quoted part", callerID)
		}
	}
	if want := "\"Fred Freddy Flintstone\" <100>"; entities[0].CallerID != want {
		t.Errorf("SIP callerid = %q, want %q", entities[0].CallerID, want)
	}
	if want := "\"Lobby Front\" <200>"; endpoints[1].CallerID != want {
		t.Errorf("Endpoints callerid = %q, want %q", endpoints[1].CallerID, want)
	}
}
