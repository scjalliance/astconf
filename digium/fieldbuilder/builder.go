package fieldbuilder

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/digium"
)

// Builder constructs a slice of digium busy lamp field entries.
//
// An empty field builder is ready for use.
type Builder struct {
	fields digium.FieldList
}

// Add adds custom busy lamp field entries to the builder.
func (builder *Builder) Add(fields ...digium.Field) (next int) {
	builder.fields = append(builder.fields, fields...)
	if len(fields) == 0 {
		return 0
	}
	return fields[len(fields)-1].Index + 1
}

// AddPerson adds a busy lamp field entry for each person to the builder,
// starting at index. It returns the next index.
func (builder *Builder) AddPerson(location string, index int, people ...astorg.Person) (next int) {
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
func (builder *Builder) AddPhoneRole(location string, index int, roles ...astorg.PhoneRole) (next int) {
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

// Fields returns a compiled field list from the builder.
func (builder *Builder) Fields() digium.FieldList {
	return builder.fields
}

// Reset purges the current field list from the builder.
func (builder *Builder) Reset() {
	builder.fields = nil
}
