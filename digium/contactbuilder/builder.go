package contactbuilder

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/digium"
)

// Builder constructs a slice of digium contact entries.
//
// An empty contact builder is ready for use.
type Builder struct {
	contacts []digium.Contact
}

// Add adds custom contact entries to the builder.
func (builder *Builder) Add(contacts ...digium.Contact) {
	builder.contacts = append(builder.contacts, contacts...)
}

// AddPerson adds a contact entry for each person to the builder.
//
// The PersonFactory function will be used to convert people into contacts.
func (builder *Builder) AddPerson(people ...astorg.Person) {
	builder.AddCustomPerson(PersonFactory, people...)
}

// AddCustomPerson adds a contact entry for each person to the builder.
//
// The function f will be used to convert people into contacts.
func (builder *Builder) AddCustomPerson(f PersonFunc, people ...astorg.Person) {
	for _, person := range people {
		var contact digium.Contact
		f(&person, &contact)
		builder.contacts = append(builder.contacts, contact)
	}
}

// AddPhoneRole adds a contact entry for each person to the builder.
//
// The PhoneRoleFactory function will be used to convert people into contacts.
func (builder *Builder) AddPhoneRole(roles ...astorg.PhoneRole) {
	builder.AddCustomPhoneRole(PhoneRoleFactory, roles...)
}

// AddCustomPhoneRole adds a contact entry for each phone role to the builder.
//
// The function f will be used to convert phone roles into contacts.
func (builder *Builder) AddCustomPhoneRole(f PhoneRoleFunc, roles ...astorg.PhoneRole) {
	for _, role := range roles {
		var contact digium.Contact
		f(&role, &contact)
		builder.contacts = append(builder.contacts, contact)
	}
}

// Group returns a compiled group from the builder.
func (builder *Builder) Group(name, id string) digium.ContactGroup {
	return digium.ContactGroup{
		Name:     name,
		ID:       id,
		Contacts: builder.contacts,
	}
}

// Reset purges the current contact list from the builder.
func (builder *Builder) Reset() {
	builder.contacts = nil
}
