package dialplan_test

import (
	"bytes"
	"fmt"

	"github.com/scjalliance/astconf"
	. "github.com/scjalliance/astconf/dialplan"
)

// Employee holds contact information for Slate Rock and Gravel Company
// personnel.
type Employee struct {
	Number string
	Name   string
	Phone  Device
}

func Example() {
	section := Section{Context: "slate-employees"}
	users := []Employee{
		{Number: "100", Name: "Fred Flintstone", Phone: SIP("fred.flintstone")},
		{Number: "101", Name: "Barney Rubble", Phone: SIP("barney.rubble")},
	}
	for _, user := range users {
		section.Extensions = append(section.Extensions, Extension{
			Comment: user.Name,
			Number:  user.Number,
			Actions: []Action{
				Noop(fmt.Sprintf("Call %s", user.Name)),
				ExecIf(Equal(DeviceState(user.Phone), String("NOT_INUSE")), Dial(user.Phone, 20)),
				Congestion(),
			},
		})
	}
	var buf bytes.Buffer
	e := astconf.NewEncoder(&buf)
	err := e.Encode(&section)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(buf.String())
	}
	// Output:
	// [slate-employees]
	//
	// ; Fred Flintstone
	// exten => 100,1,Noop(Call Fred Flintstone)
	// same => n,ExecIf($[${DEVICE_STATE(SIP/fred.flintstone)}=NOT_INUSE]?Dial(SIP/fred.flintstone,20))
	// same => n,Congestion()
	//
	// ; Barney Rubble
	// exten => 101,1,Noop(Call Barney Rubble)
	// same => n,ExecIf($[${DEVICE_STATE(SIP/barney.rubble)}=NOT_INUSE]?Dial(SIP/barney.rubble,20))
	// same => n,Congestion()
}
