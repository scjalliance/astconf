package sip_test

import (
	"bytes"
	"fmt"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/sip"
)

func Example() {
	entity := sip.Entity{
		Username:   "fred.flintstone",
		Templates:  []string{"scjdefault"},
		Type:       sip.Friend,
		Disallow:   []string{"all"},
		Allow:      []string{"ulaw", "g722"},
		AllowGuest: astval.No,
		CallerID:   "\"Fred Flintstone\" <100>",
		Context:    "from-internal",
		Host:       "dynamic",
		Transport:  []string{"tcp", "udp"},
		Mailbox:    "100@default",
		Secret:     "bedrock",
		Variables: []astval.Var{
			astval.NewVar("USER_LOCATION", "QRY"),
			astval.NewVar("AREACODE", "360"),
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&entity); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(buf.String())
	// Output:
	// [fred.flintstone](scjdefault)
	// type = friend
	// disallow = all
	// allow = ulaw
	// allow = g722
	// allowguest = no
	// callerid = "Fred Flintstone" <100>
	// context = from-internal
	// host = dynamic
	// transport = tcp,udp
	// mailbox = 100@default
	// secret = bedrock
	// setvar = USER_LOCATION=QRY
	// setvar = AREACODE=360
}
