package contactbuilder

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/digium"
)

// PersonFunc is a function that converts people into contacts.
type PersonFunc func(from *astorg.Person, to *digium.Contact)

// PersonContact applies f to from and to.
func (f PersonFunc) PersonContact(from *astorg.Person, to *digium.Contact) {
	f(from, to)
}

// Then returns a new function that will execute f followed by each member
// of series.
func (f PersonFunc) Then(series ...PersonFunc) PersonFunc {
	return func(from *astorg.Person, to *digium.Contact) {
		f(from, to)
		for _, g := range series {
			g(from, to)
		}
	}
}

// PersonFactory converts people into contacts.
var PersonFactory PersonFunc = func(person *astorg.Person, contact *digium.Contact) {
	//contact.ServerUUID:      person.Server
	contact.ID = person.Username
	contact.Type = "sip"
	contact.FirstName = person.FirstName
	contact.LastName = person.LastName
	contact.Organization = person.Organization
	contact.JobTitle = person.Title
	contact.Location = person.Location
	contact.SubscriptionURI = "auto_hint_" + person.Username
	contact.AccountID = person.Extension

	contact.Actions = append(contact.Actions, digium.ContactAction{
		ID:    "primary",
		Label: "Desk Phone",
		Name:  "Desk Phone",
		Dial:  person.Extension,
	})

	for _, number := range person.ContactNumbers {
		contact.Actions = append(contact.Actions, digium.ContactAction{
			ID:         sanitize(number.Label, "_"),
			Label:      number.Label,
			Name:       number.Label,
			Dial:       number.Dial,
			DialPrefix: number.DialPrefix,
		})
	}
	if person.VoicemailExtension != "" {
		contact.Actions = append(contact.Actions, digium.ContactAction{
			ID:    "send_to_vm",
			Label: "Voicemail",
			Name:  "Voicemail",
			Dial:  person.VoicemailExtension,
			Headers: []digium.ContactHeader{
				digium.ContactHeader{
					Key:   "X-Digium-Call-Feature",
					Value: "feature_send_to_vm",
				},
				digium.ContactHeader{
					Key:   "Diversion",
					Value: `<sip:%_ACCOUNT_USERNAME_%@%_ACCOUNT_SERVER_%:%_ACCOUNT_PORT_%>;reason="send_to_vm"`,
				},
			},
		})
	}
	for _, email := range person.EmailAddresses {
		contact.Emails = append(contact.Emails, digium.ContactEmail{
			Address: email.Address,
			Label:   email.Label,
			Primary: email.Primary,
		})
	}
}
