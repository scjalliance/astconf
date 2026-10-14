package pjsip_test

import (
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestAsteriskPublicationMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.AsteriskPublication{Name: "den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[den]\ntype = asterisk-publication\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestAsteriskPublicationFull(t *testing.T) {
	publication := pjsip.AsteriskPublication{
		Name:                "den",
		DeviceStatePublish:  "publish-den",
		MailboxStatePublish: "publish-mwi-den",
		DeviceState:         astval.Yes,
		DeviceStateFilter:   "^PJSIP/(barney)$",
		MailboxState:        astval.No,
		MailboxStateFilter:  "^200@default$",
	}
	got, err := astconf.Marshal(&publication)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[den]\n" +
		"type = asterisk-publication\n" +
		"devicestate_publish = publish-den\n" +
		"mailboxstate_publish = publish-mwi-den\n" +
		"device_state = yes\n" +
		"device_state_filter = ^PJSIP/(barney)$\n" +
		"mailbox_state = no\n" +
		"mailbox_state_filter = ^200@default$\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
