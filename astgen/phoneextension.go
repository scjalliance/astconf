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

	// Build a map of SIP identities that we can use later to determine
	// which ones are usable, and thus which ones we should check for
	// in-use status
	usablePhones := make(map[string]bool)
	{
		lookup := data.Lookup()
		for _, phone := range data.Phones {
			username := lineUsername(phone.MAC, lookup)
			if username != "" {
				usablePhones[username] = true
			}
		}
	}

	// Step 1: Add extensions for people
	for _, person := range data.People {
		if person.Extension == "" || person.Username == "" {
			continue
		}

		allDevices := makeDevices(append([]string{person.Username}, person.Softphones...)...)
		usableDevices := allDevices
		if !usablePhones[person.Username] && len(person.Softphones) > 0 {
			// If someone has a softphone but no physical desk phone, don't
			// check whether the desk phone is in use when dialing
			usableDevices = makeDevices(person.Softphones...)
		}
		presence := fmt.Sprintf("CustomPresence:%s", person.Username)

		ext := dialplan.Extension{
			Comment: person.FullName,
			Number:  person.Extension,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Call %s", person.Username)),
				dialplan.Macro("pre-answer", dialplan.Var("EXTEN")),
				dialplan.ExecIf(devicesNotInUse(usableDevices), dialplan.DialMany(allDevices, 20)),
				dialplan.ExecIf(dialplan.Equal(dialplan.PresenceState(presence, "subtype"), dialplan.String("with Call Waiting Enabled")), dialplan.DialMany(allDevices, 20)),
				dialplan.Macro("no-answer", dialplan.Var("EXTEN")),
			},
		}

		section.Extensions = append(section.Extensions, ext)
	}

	// Step 2: Add extensions for phone roles
	for _, role := range data.PhoneRoles {
		if role.Extension == "" || role.Username == "" {
			continue
		}

		allDevices := makeDevices(append([]string{role.Username}, role.Softphones...)...)
		usableDevices := allDevices
		if !usablePhones[role.Username] && len(role.Softphones) > 0 {
			// If someone has a softphone but no physical desk phone, don't
			// check whether the desk phone is in use when dialing
			usableDevices = makeDevices(role.Softphones...)
		}
		presence := fmt.Sprintf("CustomPresence:%s", role.Username)

		ext := dialplan.Extension{
			Comment: role.DisplayName,
			Number:  role.Extension,
			Actions: []dialplan.Action{
				dialplan.Noop(fmt.Sprintf("Call %s", role.Username)),
				dialplan.Macro("pre-answer", dialplan.Var("EXTEN")),
				dialplan.ExecIf(devicesNotInUse(usableDevices), dialplan.DialMany(allDevices, 20)),
				dialplan.ExecIf(dialplan.Equal(dialplan.PresenceState(presence, "subtype"), dialplan.String("with Call Waiting Enabled")), dialplan.DialMany(allDevices, 20)),
				dialplan.Congestion(),
			},
		}

		section.Extensions = append(section.Extensions, ext)
	}

	return section
}
