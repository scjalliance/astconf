package astorgconv

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/digium"
)

// FieldBuilder constructs a slice of digium busy lamp fields for
// astorg types.
//
// An empty field builder is ready for use.
type FieldBuilder struct {
	fields digium.FieldList
	locs   map[string]astorg.Location
}

// AddLocation adds a location to the builder.
func (builder *FieldBuilder) AddLocation(locations ...astorg.Location) {
	for _, loc := range locations {
		builder.locs[loc.Name] = loc
	}
}

// AddPerson adds a contact entry for each person to the builder.
//
// TODO: Let the caller provide some sort of non-default conversion function?
func (builder *FieldBuilder) AddPerson(location string, index int, people ...astorg.Person) (next int) {
	for _, person := range people {
		field := digium.Field{
			Location:  location,
			Index:     index,
			ContactID: person.Username,
		}
		builder.fields = append(builder.fields, field)
		index++
	}
	return index
}

// AddPhoneRole adds a busy lamp field entry for each phone role to
// the builder, starting at index. It returns the next index.
//
// TODO: Let the caller provide some sort of non-default conversion function?
func (builder *FieldBuilder) AddPhoneRole(location string, index int, roles ...astorg.PhoneRole) (next int) {
	for _, role := range roles {
		field := digium.Field{
			Location:  location,
			Index:     index,
			ContactID: role.Username,
		}
		builder.fields = append(builder.fields, field)
		index++
	}
	return index
}

// AddCustom adds custom busy lamp field entries to the builder.
func (builder *FieldBuilder) AddCustom(fields ...digium.Field) (next int) {
	builder.fields = append(builder.fields, fields...)
	if len(fields) == 0 {
		return 0
	}
	return fields[len(fields)-1].Index + 1
}

// Fields returns a compiled field list from the builder.
func (builder *FieldBuilder) Fields() digium.FieldList {
	return builder.fields
}

// Reset purges the current field list from the builder.
func (builder *FieldBuilder) Reset() {
	builder.fields = nil
}
