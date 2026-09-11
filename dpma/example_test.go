package dpma_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astval"
	"github.com/scjalliance/astconf/dpma"
)

func ExamplePhone() {
	phone := dpma.Phone{
		Username:       "fred.flintstone",
		Networks:       []string{"quarry"},
		MAC:            "000FD3000001",
		Lines:          []string{"fred.flintstone", "shared"},
		FullName:       "Fred Flintstone",
		Timezone:       "America/Los_Angeles",
		Ringtones:      []string{"Alarm"},
		ActiveRingtone: "Alarm",
		WebUIEnabled:   astval.No,
		Brightness:     astval.NewInt(5),
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&phone); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(buf.String())
	// Output:
	// [fred.flintstone]
	// type = phone
	// network = quarry
	// mac = 000FD3000001
	// line = fred.flintstone
	// line = shared
	// full_name = Fred Flintstone
	// timezone = America/Los_Angeles
	// ringtone = Alarm
	// active_ringtone = Alarm
	// web_ui_enabled = no
	// brightness = 5
}

func ExampleLine() {
	line := dpma.Line{
		Name:                "fred.flintstone",
		Extension:           "100",
		Label:               "Fred",
		Mailbox:             "100@default",
		Transport:           "tcp",
		RegistrationTimeout: astval.NewSeconds(5 * time.Minute),
		Secret:              "bedrock",
		Context:             "from-internal",
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&line); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(buf.String())
	// Output:
	// [fred.flintstone]
	// type = line
	// exten = 100
	// line_label = Fred
	// mailbox = 100@default
	// transport = tcp
	// reregistration_timeout = 300
	// secret = bedrock
	// context = from-internal
}
