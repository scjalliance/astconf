package astgen

import (
	"fmt"

	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/dialplan"
)

// PhoneExtensions generates a dialplan section that includes all phone
// extensions in a dataset.
func PhoneExtensions(data *astorg.DataSet, context string) dialplan.Section {
	section := dialplan.Section{Context: context}

	for _, person := range data.People {
		if person.Extension == "" || person.Username == "" {
			continue
		}

		device := dialplan.SIP(person.Username)
		presence := fmt.Sprintf("CustomPresence:%s", person.Username)

		ext := dialplan.Extension{
			Comment: person.FullName,
			Number:  person.Extension,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Call %s", person.Username)),
				dialplan.Macro("pre-answer", dialplan.Var("EXTEN")),
				dialplan.ExecIf(dialplan.Equal(dialplan.DeviceState(device), dialplan.String("NOT_INUSE")), dialplan.Dial(device, 20)),
				dialplan.ExecIf(dialplan.Equal(dialplan.PresenceState(presence, "subtype"), dialplan.String("with Call Waiting Enabled")), dialplan.Dial(device, 20)),
				dialplan.Macro("no-answer", dialplan.Var("EXTEN")),
			},
		}

		section.Extensions = append(section.Extensions, ext)
	}

	for _, role := range data.PhoneRoles {
		if role.Extension == "" || role.Username == "" {
			continue
		}

		device := dialplan.SIP(role.Username)
		presence := fmt.Sprintf("CustomPresence:%s", role.Username)

		ext := dialplan.Extension{
			Comment: role.DisplayName,
			Number:  role.Extension,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Call %s", role.Username)),
				dialplan.Macro("pre-answer", dialplan.Var("EXTEN")),
				dialplan.ExecIf(dialplan.Equal(dialplan.DeviceState(device), dialplan.String("NOT_INUSE")), dialplan.Dial(device, 20)),
				dialplan.ExecIf(dialplan.Equal(dialplan.PresenceState(presence, "subtype"), dialplan.String("with Call Waiting Enabled")), dialplan.Dial(device, 20)),
				dialplan.Congestion(),
			},
		}

		section.Extensions = append(section.Extensions, ext)
	}

	return section
}
