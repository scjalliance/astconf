package astorg

// MailboxList is a slice of mailboxes.
type MailboxList []Mailbox

// ByName returns a map of mailboxes indexed by mailbox name.
func (mailboxes MailboxList) ByName() map[string]Mailbox {
	lookup := make(map[string]Mailbox, len(mailboxes))
	for _, mailbox := range mailboxes {
		if mailbox.Name == "" {
			continue
		}
		lookup[mailbox.Name] = mailbox
	}
	return lookup
}
