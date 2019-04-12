package astorg

import "github.com/scjalliance/astconf/astorg/astorgvm"

// Mailbox is a voicemail mailbox.
type Mailbox struct {
	Name       string
	Number     string
	AccessCode string
	Email      string
	AccessMode astorgvm.Access
}
