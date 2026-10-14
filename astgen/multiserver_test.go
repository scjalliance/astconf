package astgen_test

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astorg/astorgvm"
)

// multiServerData returns a data set with two office servers and a hub.
//
// Server IDs sort differently than server names, so a test that orders by
// name instead of ID fails. Quarry and Den have servers; Remote is an
// office with no server, so extensions there live on the hub.
//
//   - fred (100) has a desk phone at Quarry and a softphone, so his
//     extension lives on qry and rings one device on the hub.
//   - barney (101) has a desk phone at Den.
//   - wilma (102) works at Remote and has only a softphone, so her
//     extension and her one device both live on the hub.
//   - The lobby role (200) has a desk phone at Den and a mailbox.
//   - The yard role (201) is at Remote with only a softphone.
//   - A spare phone at Quarry is assigned to nobody.
//   - The All Call paging group (500) covers Quarry and Den, so none of its
//     members are on the hub.
func multiServerData() *astorg.DataSet {
	return &astorg.DataSet{
		Servers: astorg.ServerList{
			{ID: "qry", Name: "Quarry", Address: "qry-phone"},
			{ID: "den", Name: "Rock Den", Address: "den-phone"},
			{ID: "hub", Name: "Hub", Hub: true, Address: "hub-phone"},
		},
		Locations: astorg.LocationList{
			{Name: "Quarry", Abbreviation: "QRY", Server: "qry", CallerID: "5555550100", AreaCode: "555", PagingGroups: []string{"500"}},
			{Name: "Den", Abbreviation: "DEN", Server: "den", CallerID: "5555550200", AreaCode: "555", PagingGroups: []string{"500"}},
			{Name: "Remote", Abbreviation: "RMT", CallerID: "5555550300", AreaCode: "555"},
		},
		People: astorg.PersonList{
			{
				Username:        "fred",
				FullName:        "Fred Flintstone",
				Extension:       "100",
				Location:        "Quarry",
				VoicemailCode:   "1234",
				VoicemailAccess: astorgvm.PhoneAndEmail,
				EmailAddresses:  []astorg.Email{{Address: "fred@example.com", Primary: true}},
				Phones:          []string{"000fd3000001"},
				Softphones:      []string{"fred.soft"},
			},
			{
				Username:      "barney",
				FullName:      "Barney Rubble",
				Extension:     "101",
				Location:      "Den",
				VoicemailCode: "4321",
				Phones:        []string{"000fd3000004"},
			},
			{
				Username:   "wilma",
				FullName:   "Wilma Flintstone",
				Extension:  "102",
				Location:   "Remote",
				Softphones: []string{"wilma.soft"},
			},
		},
		PhoneRoles: astorg.PhoneRoleList{
			{
				Username:      "lobby",
				DisplayName:   "Lobby",
				Extension:     "200",
				Location:      "Den",
				MailboxNumber: "300",
				Phones:        []string{"000fd3000002"},
			},
			{
				Username:    "yard",
				DisplayName: "Yard",
				Extension:   "201",
				Location:    "Remote",
				Softphones:  []string{"yard.soft"},
			},
		},
		Phones: astorg.PhoneList{
			{MAC: "000fd3000001", Location: "Quarry", Secret: "fred-secret"},
			{MAC: "000fd3000002", Location: "Den", Secret: "lobby-secret"},
			{MAC: "000fd3000003", Location: "Quarry", Secret: "spare-secret"},
			{MAC: "000fd3000004", Location: "Den", Secret: "barney-secret"},
		},
		Softphones: astorg.SoftphoneList{
			{Username: "fred.soft", Location: "Quarry", Secret: "fred-soft-secret"},
			{Username: "wilma.soft", Location: "Remote", Secret: "wilma-soft-secret"},
			{Username: "yard.soft", Location: "Remote", Secret: "yard-soft-secret"},
		},
		PagingGroups: astorg.PagingGroupList{
			{Name: "All Call", Extension: "500", Alert: "intercom"},
		},
		Mailboxes: astorg.MailboxList{
			{Number: "300", Name: "Sales", AccessCode: "3000", AccessMode: astorgvm.Email, Email: "sales@example.com"},
		},
	}
}
