package astgen_test

import (
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
	"github.com/scjalliance/astconf/sip"
)

func TestEndpoints(t *testing.T) {
	base := pjsip.Endpoint{
		Context: "from-internal",
	}
	endpoints := astgen.Endpoints(testData(), base, "default")

	var values []interface{}
	for i := range endpoints {
		values = append(values, &endpoints[i])
	}
	got := render(t, values...)

	want := "[fred](desk)\n" +
		"type = endpoint\n" +
		"context = from-internal\n" +
		"auth = fred-auth\n" +
		"aors = fred\n" +
		"callerid = \"Fred Flintstone\" <100>\n" +
		"mailboxes = 100@default\n" +
		"set_var = VMCODE=1234\n" +
		"set_var = MOBILE=5555551234\n" +
		"set_var = USERNAME=fred\n" +
		"set_var = USER_LOCATION=Quarry\n" +
		"set_var = OUTBOUND_CALLERID=5555550100\n" +
		"set_var = AREACODE=555\n" +
		"\n" +
		"[lobby]\n" +
		"type = endpoint\n" +
		"context = from-internal\n" +
		"auth = lobby-auth\n" +
		"aors = lobby\n" +
		"callerid = \"Lobby\" <200>\n" +
		"mailboxes = 200@default\n" +
		"set_var = USER_LOCATION=Quarry\n" +
		"set_var = OUTBOUND_CALLERID=5555550100\n" +
		"set_var = AREACODE=555\n" +
		"\n" +
		"[000fd3000003]\n" +
		"type = endpoint\n" +
		"context = from-internal\n" +
		"auth = 000fd3000003-auth\n" +
		"aors = 000fd3000003\n" +
		"callerid = \"QRY-000FD3000003\" <UNAVAILABLE>\n" +
		"set_var = USER_LOCATION=Quarry\n" +
		"set_var = OUTBOUND_CALLERID=5555550100\n" +
		"set_var = AREACODE=555\n" +
		"\n" +
		"[fred.soft](soft)\n" +
		"type = endpoint\n" +
		"context = from-internal\n" +
		"auth = fred.soft-auth\n" +
		"aors = fred.soft\n" +
		"callerid = \"Fred Flintstone\" <100>\n" +
		"set_var = VMCODE=1234\n" +
		"set_var = MOBILE=5555551234\n" +
		"set_var = USERNAME=fred\n" +
		"set_var = USER_LOCATION=Quarry\n" +
		"set_var = OUTBOUND_CALLERID=5555550100\n" +
		"set_var = AREACODE=555\n" +
		"\n" +
		"[lobby.soft]\n" +
		"type = endpoint\n" +
		"context = from-internal\n" +
		"auth = lobby.soft-auth\n" +
		"aors = lobby.soft\n" +
		"callerid = \"Lobby\" <200>\n"
	if got != want {
		t.Errorf("Endpoints() rendered:\n%s\nwant:\n%s", got, want)
	}
}

func TestAuths(t *testing.T) {
	base := pjsip.Auth{
		AuthType: "userpass",
	}
	auths := astgen.Auths(testData(), base)

	var values []interface{}
	for i := range auths {
		values = append(values, &auths[i])
	}
	got := render(t, values...)

	want := "[fred-auth]\n" +
		"type = auth\n" +
		"auth_type = userpass\n" +
		"username = fred\n" +
		"password = fred-secret\n" +
		"\n" +
		"[lobby-auth]\n" +
		"type = auth\n" +
		"auth_type = userpass\n" +
		"username = lobby\n" +
		"password = lobby-secret\n" +
		"\n" +
		"[000fd3000003-auth]\n" +
		"type = auth\n" +
		"auth_type = userpass\n" +
		"username = 000fd3000003\n" +
		"password = spare-secret\n" +
		"\n" +
		"[fred.soft-auth]\n" +
		"type = auth\n" +
		"auth_type = userpass\n" +
		"username = fred.soft\n" +
		"password = fred-soft-secret\n" +
		"\n" +
		"[lobby.soft-auth]\n" +
		"type = auth\n" +
		"auth_type = userpass\n" +
		"username = lobby.soft\n" +
		"password = lobby-soft-secret\n"
	if got != want {
		t.Errorf("Auths() rendered:\n%s\nwant:\n%s", got, want)
	}
}

func TestAORs(t *testing.T) {
	base := pjsip.AOR{
		MaxContacts:    astval.NewInt(5),
		RemoveExisting: astval.Yes,
	}
	aors := astgen.AORs(testData(), base)

	var values []interface{}
	for i := range aors {
		values = append(values, &aors[i])
	}
	got := render(t, values...)

	want := "[fred]\n" +
		"type = aor\n" +
		"max_contacts = 5\n" +
		"remove_existing = yes\n" +
		"\n" +
		"[lobby]\n" +
		"type = aor\n" +
		"max_contacts = 5\n" +
		"remove_existing = yes\n" +
		"\n" +
		"[000fd3000003]\n" +
		"type = aor\n" +
		"max_contacts = 5\n" +
		"remove_existing = yes\n" +
		"\n" +
		"[fred.soft]\n" +
		"type = aor\n" +
		"max_contacts = 5\n" +
		"remove_existing = yes\n" +
		"\n" +
		"[lobby.soft]\n" +
		"type = aor\n" +
		"max_contacts = 5\n" +
		"remove_existing = yes\n"
	if got != want {
		t.Errorf("AORs() rendered:\n%s\nwant:\n%s", got, want)
	}
}

// TestEndpointsSkipUnnamedPhones checks that a phone with no MAC address
// produces no endpoint, matching the SIP generator.
func TestEndpointsSkipUnnamedPhones(t *testing.T) {
	data := testData()
	data.Phones = append(data.Phones, data.Phones[0])
	data.Phones[len(data.Phones)-1].MAC = ""
	if got, want := len(astgen.Endpoints(data, pjsip.Endpoint{}, "default")), 5; got != want {
		t.Errorf("Endpoints() returned %d endpoints, want %d", got, want)
	}
}

// TestSoftphoneWithoutUsername verifies that a software phone carrying no
// username is skipped rather than rendered as an unnamed section. Asterisk
// rejects a "[]" section header, so emitting one breaks the whole file.
func TestSoftphoneWithoutUsername(t *testing.T) {
	data := &astorg.DataSet{
		Softphones: astorg.SoftphoneList{
			{Username: "", Secret: "unused"},
		},
	}

	if got := astgen.Endpoints(data, pjsip.Endpoint{}, "default"); len(got) != 0 {
		t.Errorf("Endpoints() returned %d entries for an unnamed softphone, want 0", len(got))
	}
	if got := astgen.Auths(data, pjsip.Auth{}); len(got) != 0 {
		t.Errorf("Auths() returned %d entries for an unnamed softphone, want 0", len(got))
	}
	if got := astgen.AORs(data, pjsip.AOR{}); len(got) != 0 {
		t.Errorf("AORs() returned %d entries for an unnamed softphone, want 0", len(got))
	}
	if got := astgen.SIP(data, sip.Entity{}, "default"); len(got) != 0 {
		t.Errorf("SIP() returned %d entries for an unnamed softphone, want 0", len(got))
	}
}

// TestDeviceAORIsItsEndpoint: every device endpoint names exactly one aor,
// with the endpoint's own name, which a phone registers under (#23).
func TestDeviceAORIsItsEndpoint(t *testing.T) {
	data := testData()
	aors := map[string]bool{}
	for _, a := range astgen.AORs(data, pjsip.AOR{}) {
		aors[a.Name] = true
	}
	for _, ep := range astgen.Endpoints(data, pjsip.Endpoint{}, "default") {
		if len(ep.AORs) != 1 || ep.AORs[0] != ep.Name {
			t.Errorf("endpoint %s: aors %v, want [%s]", ep.Name, ep.AORs, ep.Name)
		}
		if !aors[ep.Name] {
			t.Errorf("endpoint %s: no aor named %s", ep.Name, ep.Name)
		}
	}
}
