package astorgconv

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/digium"
)

// ContactBuilder constructs a slice of dpma contact entries for
// astorg types.
//
// An empty contact builder is ready for use.
type ContactBuilder struct {
	contacts []digium.Contact
	locs     map[string]astorg.Location
}

// AddLocation adds a location to the builder.
func (builder *ContactBuilder) AddLocation(locations ...astorg.Location) {
	for _, loc := range locations {
		builder.locs[loc.Name] = loc
	}
}

// AddPerson adds a contact entry for each person to the builder.
//
// TODO: Let the caller provide some sort of non-default conversion function?
func (builder *ContactBuilder) AddPerson(people ...astorg.Person) {
	for _, person := range people {
		contact := digium.Contact{
			//ServerUUID:      person.Server,
			Type:            "sip",
			FirstName:       person.FirstName,
			LastName:        person.LastName,
			Organization:    person.Organization,
			JobTitle:        person.Title,
			Location:        person.Location,
			SubscriptionURI: "auto_hint_" + person.Username,
			AccountID:       person.Extension,
			Actions: []digium.ContactAction{
				digium.ContactAction{
					ID:    "primary",
					Label: "Desk Phone",
					Name:  "Desk Phone",
					Dial:  person.Extension,
				},
			},
		}
		for _, number := range person.ContactNumbers {
			contact.Actions = append(contact.Actions, digium.ContactAction{
				ID:         sanitizeNameUnderscored(number.Label),
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
		builder.contacts = append(builder.contacts, contact)
	}
}

// AddPhoneRole adds a contact entry for each phone role to the builder.
//
// TODO: Let the caller provide some sort of non-default conversion function?
func (builder *ContactBuilder) AddPhoneRole(roles ...astorg.PhoneRole) {
	for _, role := range roles {
		contact := digium.Contact{
			//ServerUUID:      person.Server,
			Type:            "sip",
			FirstName:       role.DisplayName,
			Location:        role.Location,
			SubscriptionURI: "auto_hint_" + role.Username,
			AccountID:       role.Extension,
			Actions: []digium.ContactAction{
				digium.ContactAction{
					ID:    "primary",
					Label: "Desk Phone",
					Name:  "Desk Phone",
					Dial:  role.Extension,
				},
			},
		}
		builder.contacts = append(builder.contacts, contact)
	}
}

// AddCustom adds custom contact entries to the builder.
func (builder *ContactBuilder) AddCustom(contacts ...digium.Contact) {
	builder.contacts = append(builder.contacts, contacts...)
}

// Group returns a compiled group from the builder.
func (builder *ContactBuilder) Group(name, id string) digium.ContactGroup {
	return digium.ContactGroup{
		Name:     name,
		ID:       id,
		Contacts: builder.contacts,
	}
}

// Reset purges the current contact set from the builder.
func (builder *ContactBuilder) Reset() {
	builder.contacts = nil
}
