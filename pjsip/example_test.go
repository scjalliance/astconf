package pjsip_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/pjsip"
)

func Example() {
	endpoint := pjsip.Endpoint{
		Name:      "132",
		Templates: []string{"scjdefault"},
		Auth:      "132-auth",
		AORs:      []string{"132-aor"},
		CallerID:  "\"Dusty Wilson\" <132>",
		Mailboxes: []string{"132@default"},
		Variables: []astval.Var{
			astval.NewVar("USER_LOCATION", "LCY"),
			astval.NewVar("AREACODE", "360"),
		},
	}
	auth := pjsip.Auth{
		Name:     "132-auth",
		AuthType: "userpass",
		Username: "132",
		Password: "secret",
	}
	aor := pjsip.AOR{
		Name:           "132-aor",
		MaxContacts:    astval.NewInt(5),
		RemoveExisting: astval.Yes,
	}

	var buf bytes.Buffer
	e := astconf.NewEncoder(&buf)
	for _, v := range []interface{}{&endpoint, &auth, &aor} {
		if err := e.Encode(v); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Print(buf.String())
	// Output:
	// [132](scjdefault)
	// type = endpoint
	// auth = 132-auth
	// aors = 132-aor
	// callerid = "Dusty Wilson" <132>
	// mailboxes = 132@default
	// set_var = USER_LOCATION=LCY
	// set_var = AREACODE=360
	//
	// [132-auth]
	// type = auth
	// auth_type = userpass
	// username = 132
	// password = secret
	//
	// [132-aor]
	// type = aor
	// max_contacts = 5
	// remove_existing = yes
}

func ExampleEndpoint_template() {
	endpoint := pjsip.Endpoint{
		Name:                       "scjbase",
		Transport:                  "transport-tcp",
		Context:                    "from-internal",
		Disallow:                   []string{"all"},
		Allow:                      []string{"g722", "ulaw", "slin16"},
		DirectMediaGlareMitigation: "none",
		RTPSymmetric:               astval.Yes,
		ForceRport:                 astval.Yes,
		RewriteContact:             astval.Yes,
		Variables: []astval.Var{
			astval.NewVar("ATTENDED_TRANSFER_COMPLETE_SOUND", "beep"),
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&endpoint); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(buf.String())
	// Output:
	// [scjbase]
	// type = endpoint
	// transport = transport-tcp
	// context = from-internal
	// disallow = all
	// allow = g722
	// allow = ulaw
	// allow = slin16
	// direct_media_glare_mitigation = none
	// rtp_symmetric = yes
	// force_rport = yes
	// rewrite_contact = yes
	// set_var = ATTENDED_TRANSFER_COMPLETE_SOUND=beep
}

func ExampleTransport() {
	transport := pjsip.Transport{
		Name:                     "transport-tls",
		Protocol:                 "tls",
		Bind:                     "0.0.0.0:5061",
		ExternalMediaAddress:     "198.51.100.10",
		ExternalSignalingAddress: "198.51.100.10",
		LocalNet:                 []string{"10.100.0.0/16", "100.64.0.0/10"},
		CertFile:                 "/etc/asterisk/keys/asterisk.pem",
		PrivKeyFile:              "/etc/asterisk/keys/asterisk.key",
		CAListFile:               "/etc/asterisk/keys/ca.crt",
		Method:                   "tlsv1_2",
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&transport); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(buf.String())
	// Output:
	// [transport-tls]
	// type = transport
	// protocol = tls
	// bind = 0.0.0.0:5061
	// external_media_address = 198.51.100.10
	// external_signaling_address = 198.51.100.10
	// local_net = 10.100.0.0/16
	// local_net = 100.64.0.0/10
	// cert_file = /etc/asterisk/keys/asterisk.pem
	// priv_key_file = /etc/asterisk/keys/asterisk.key
	// ca_list_file = /etc/asterisk/keys/ca.crt
	// method = tlsv1_2
}

func ExampleAOR_qualify() {
	aor := pjsip.AOR{
		Name:             "twilio-aor",
		Contacts:         []string{"sip:example.pstn.twilio.com"},
		QualifyFrequency: astval.NewSeconds(time.Minute),
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&aor); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(buf.String())
	// Output:
	// [twilio-aor]
	// type = aor
	// contact = sip:example.pstn.twilio.com
	// qualify_frequency = 60
}
