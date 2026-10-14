package pjsip_test

import (
	"testing"
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func TestOutboundPublishMinimal(t *testing.T) {
	got, err := astconf.Marshal(&pjsip.OutboundPublish{Name: "publish-den"})
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[publish-den]\ntype = outbound-publish\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}

func TestOutboundPublishFull(t *testing.T) {
	publish := pjsip.OutboundPublish{
		Name:               "publish-den",
		ServerURI:          "sip:trunk-lcy@den-phone:5060;transport=tcp",
		FromURI:            "sip:lcy@lcy-phone",
		ToURI:              "sip:den@den-phone",
		Event:              "asterisk-devicestate",
		OutboundAuth:       "auth-trunk-den",
		Expiration:         astval.NewSeconds(time.Hour),
		Transport:          "transport-tcp",
		OutboundProxy:      "sip:proxy.example.com",
		MaxAuthAttempts:    astval.NewInt(3),
		DeviceStateFilter:  "^PJSIP/(fred)$",
		MailboxStateFilter: "^100@default$",
	}
	got, err := astconf.Marshal(&publish)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	// A semicolon starts a comment in asterisk configuration files, so the
	// encoder escapes each one as "\;". Asterisk reads "\;" back as ";".
	want := "[publish-den]\n" +
		"type = outbound-publish\n" +
		"server_uri = sip:trunk-lcy@den-phone:5060\\;transport=tcp\n" +
		"from_uri = sip:lcy@lcy-phone\n" +
		"to_uri = sip:den@den-phone\n" +
		"event = asterisk-devicestate\n" +
		"outbound_auth = auth-trunk-den\n" +
		"expiration = 3600\n" +
		"transport = transport-tcp\n" +
		"outbound_proxy = sip:proxy.example.com\n" +
		"max_auth_attempts = 3\n" +
		"@device_state_filter = ^PJSIP/(fred)$\n" +
		"@mailbox_state_filter = ^100@default$\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
