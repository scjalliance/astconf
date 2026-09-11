package astgen_test

import (
	"bytes"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astorg/astorgvm"
)

// testData returns a small dataset that exercises each step of the
// generators: an unassigned phone, a phone assigned to a person, a phone
// assigned to a role, a softphone assigned to a person, and a softphone
// assigned to a role.
func testData() *astorg.DataSet {
	return &astorg.DataSet{
		Locations: astorg.LocationList{
			{
				Name:         "Quarry",
				Abbreviation: "QRY",
				Timezone:     "America/Los_Angeles",
				CallerID:     "3605550100",
				AreaCode:     "360",
			},
		},
		People: astorg.PersonList{
			{
				Username:        "fred",
				FullName:        "Fred Flintstone",
				Extension:       "100",
				Location:        "Quarry",
				VoicemailCode:   "1234",
				VoicemailAccess: astorgvm.PhoneAndEmail,
				ContactNumbers: []astorg.Number{
					{Label: "Mobile", Dial: "3605551234"},
				},
				EmailAddresses: []astorg.Email{
					{Address: "fred@example.com", Primary: true},
				},
				Phones:     []string{"000fd3000001"},
				Softphones: []string{"fred.soft"},
			},
		},
		PhoneRoles: astorg.PhoneRoleList{
			{
				Username:      "lobby",
				DisplayName:   "Lobby",
				Extension:     "200",
				MailboxNumber: "200",
				Phones:        []string{"000fd3000002"},
				Softphones:    []string{"lobby.soft"},
			},
		},
		Phones: astorg.PhoneList{
			{MAC: "000fd3000001", Location: "Quarry", Secret: "fred-secret", Templates: []string{"desk"}},
			{MAC: "000fd3000002", Location: "Quarry", Secret: "lobby-secret"},
			{MAC: "000fd3000003", Location: "Quarry", Secret: "spare-secret"},
		},
		Softphones: astorg.SoftphoneList{
			{Username: "fred.soft", Location: "Quarry", Secret: "fred-soft-secret", Templates: []string{"soft"}},
			{Username: "lobby.soft", Secret: "lobby-soft-secret"},
		},
		Mailboxes: astorg.MailboxList{
			{Number: "300", Name: "Sales", AccessCode: "3000", AccessMode: astorgvm.Email, Email: "sales@example.com"},
		},
	}
}

// render encodes each value in turn and returns the combined output.
func render(t *testing.T, values ...interface{}) string {
	t.Helper()
	var buf bytes.Buffer
	e := astconf.NewEncoder(&buf)
	for _, v := range values {
		if err := e.Encode(v); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}
	return buf.String()
}
