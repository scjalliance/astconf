package astgen

import (
	"github.com/scjalliance/astconf/astorg"
	"github.com/scjalliance/astconf/astorg/astorgvm"
	"github.com/scjalliance/astconf/voicemail"
)

// Mailboxes generates a voicemail configuration section for the dataset.
func Mailboxes(data *astorg.DataSet, context string) voicemail.Section {
	section := voicemail.Section{Context: context}

	for _, person := range data.People {
		if person.Extension == "" {
			continue
		}

		if person.VoicemailCode == "" {
			continue
		}

		mailbox := voicemail.Box{
			Extension:   person.Extension,
			Name:        person.Username,
			Password:    person.VoicemailCode,
			SayCallerID: true,
		}

		switch person.VoicemailAccess {
		case astorgvm.Phone:
			mailbox.SendToPager = true
		case astorgvm.Email:
			mailbox.EmailOnly = true
		}

		switch person.VoicemailAccess {
		case astorgvm.Default, astorgvm.Email, astorgvm.PhoneAndEmail:
			for _, addr := range person.EmailAddresses {
				if addr.Primary {
					mailbox.EmailAddresses = []string{addr.Address}
					break
				}
			}
		}

		section.Mailboxes = append(section.Mailboxes, mailbox)
	}

	return section
}
