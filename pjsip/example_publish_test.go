package pjsip_test

import (
	"bytes"
	"fmt"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

// ExampleAsteriskPublication shares device state with one other server,
// den, over a trunk that is recognized by its source address.
func ExampleAsteriskPublication() {
	identify := pjsip.Identify{
		Name:     "trunk-den",
		Endpoint: "trunk-den",
		Match:    []string{"den-phone"},
	}
	publish := pjsip.OutboundPublish{
		Name:              "publish-den",
		ServerURI:         "sip:trunk-lcy@den-phone:5060;transport=tcp",
		Event:             "asterisk-devicestate",
		OutboundAuth:      "auth-trunk-den",
		DeviceStateFilter: "^PJSIP/(fred)$",
	}
	inbound := pjsip.InboundPublication{
		Name:     "trunk-den",
		Endpoint: "trunk-den",
		Events:   pjsip.PublicationEvents{{Event: "asterisk-devicestate", Handler: "den"}},
	}
	publication := pjsip.AsteriskPublication{
		Name:               "den",
		DeviceStatePublish: "publish-den",
		DeviceState:        astval.Yes,
		DeviceStateFilter:  "^PJSIP/(barney)$",
	}

	var buf bytes.Buffer
	e := astconf.NewEncoder(&buf)
	for _, v := range []interface{}{&identify, &publish, &inbound, &publication} {
		if err := e.Encode(v); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Print(buf.String())
	// Output:
	// [trunk-den]
	// type = identify
	// endpoint = trunk-den
	// match = den-phone
	//
	// [publish-den]
	// type = outbound-publish
	// server_uri = sip:trunk-lcy@den-phone:5060\;transport=tcp
	// event = asterisk-devicestate
	// outbound_auth = auth-trunk-den
	// @device_state_filter = ^PJSIP/(fred)$
	//
	// [trunk-den]
	// type = inbound-publication
	// endpoint = trunk-den
	// event_asterisk-devicestate = den
	//
	// [den]
	// type = asterisk-publication
	// devicestate_publish = publish-den
	// device_state = yes
	// device_state_filter = ^PJSIP/(barney)$
}
