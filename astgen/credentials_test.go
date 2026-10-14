package astgen_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/pjsip"
)

// TestCredentialsHostedPhones checks that a server lists a phone
// credential for each desk phone it hosts and each spare, sorted by MAC
// address after the trunk credentials, with the auth section each line's
// endpoint names.
func TestCredentialsHostedPhones(t *testing.T) {
	p, _ := place(t, multiServerData())
	want := []astgen.Credential{
		{Kind: "trunk", ID: "den~qry", Section: "auth-trunk-den"},
		{Kind: "trunk", ID: "hub~qry", Section: "auth-trunk-hub"},
		{Kind: "phone", ID: "000fd3000001", Section: "fred-auth"},
		{Kind: "phone", ID: "000fd3000003", Section: "000fd3000003-auth"},
		{Kind: "phone", ID: "000fd3000009", Section: "000fd3000009-auth"},
	}
	if got := p.Credentials("qry"); !reflect.DeepEqual(got, want) {
		t.Errorf("Credentials(qry) =\n%v\nwant:\n%v", got, want)
	}
}

// TestCredentialsSharedLine checks that two desk phones on one line get
// one credential, derived from the first phone in data set order, since
// they share one auth section.
func TestCredentialsSharedLine(t *testing.T) {
	data := multiServerData()
	data.Phones = append(data.Phones, astorg.Phone{MAC: "000fd3000000", Location: "Quarry", Secret: "fred-second-secret"})
	data.People[0].Phones = append(data.People[0].Phones, "000fd3000000")
	p, problems := place(t, data)
	if len(problems) != 0 {
		t.Fatalf("Place() problems: %v", problems)
	}
	var phones []astgen.Credential
	for _, cred := range p.Credentials("qry") {
		if cred.Kind == "phone" && cred.Section == "fred-auth" {
			phones = append(phones, cred)
		}
	}
	want := []astgen.Credential{{Kind: "phone", ID: "000fd3000001", Section: "fred-auth"}}
	if !reflect.DeepEqual(phones, want) {
		t.Errorf("credentials for fred-auth = %v, want %v", phones, want)
	}
}

// TestCredentialsMatchEndpoints checks that, on every server, the phone
// credentials name exactly the auth sections of the desk phone endpoints
// the server renders, so the node agent writes every section and no more.
func TestCredentialsMatchEndpoints(t *testing.T) {
	data := multiServerData()
	p, _ := place(t, data)
	softphones := make(map[string]bool)
	for _, softphone := range data.Softphones {
		softphones[softphone.Username] = true
	}
	for _, serverID := range p.ServerIDs() {
		want := make(map[string]bool)
		for _, endpoint := range astgen.ServerEndpoints(data, p, serverID, pjsip.Endpoint{}, "default") {
			if !softphones[endpoint.Name] {
				want[endpoint.Auth] = true
			}
		}
		got := make(map[string]bool)
		for _, cred := range p.Credentials(serverID) {
			if cred.Kind == "phone" {
				if got[cred.Section] {
					t.Errorf("Credentials(%s) lists %s twice", serverID, cred.Section)
				}
				got[cred.Section] = true
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Credentials(%s) phone sections = %v, want the endpoint auths %v", serverID, got, want)
		}
	}
}
