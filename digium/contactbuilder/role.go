package contactbuilder

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/digium"
)

// PhoneRoleFunc is a function that converts phone roles into contacts.
type PhoneRoleFunc func(from *astorg.PhoneRole, to *digium.Contact)

// Then returns a new function that will execute f followed by each member
// of series.
func (f PhoneRoleFunc) Then(series ...PhoneRoleFunc) PhoneRoleFunc {
	return func(from *astorg.PhoneRole, to *digium.Contact) {
		f(from, to)
		for _, g := range series {
			g(from, to)
		}
	}
}

// PhoneRoleFactory converts phone roles into contacts.
var PhoneRoleFactory PhoneRoleFunc = func(role *astorg.PhoneRole, contact *digium.Contact) {
	//ServerUUID=      person.Server
	contact.ID = role.Username
	contact.Type = "sip"
	contact.FirstName = role.DisplayName
	contact.Location = role.Location
	contact.SubscriptionURI = "auto_hint_" + role.Username
	contact.AccountID = role.Extension
	contact.Actions = append(contact.Actions, digium.ContactAction{
		ID:    "primary",
		Label: "Desk Phone",
		Name:  "Desk Phone",
		Dial:  role.Extension,
	})
}
