package dialplan_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/dialplan"
)

// TestPerServerApps renders the applications and expressions that the
// per-server generators in astgen use to build a dial string at run time.
func TestPerServerApps(t *testing.T) {
	dialVar := dialplan.Var("ASTGEN_DIAL")
	section := dialplan.Section{
		Context: "test",
		Extensions: []dialplan.Extension{
			{
				Number: "100",
				Actions: []dialplan.Action{
					dialplan.Set(dialVar, dialplan.String("&PJSIP/dev-fred.soft@trunk-hub")),
					dialplan.ExecIf(
						dialplan.Equal(dialplan.DeviceState(dialplan.PJSIP("fred")), dialplan.String("NOT_INUSE")),
						dialplan.Set(dialVar, dialplan.Cat(dialVar, dialplan.String("&PJSIP/fred"))),
					),
					dialplan.ExecIf(
						dialplan.GreaterThan(dialplan.Len(dialVar), dialplan.Int(0)),
						dialplan.DialExpr(dialplan.Substring(dialVar, 1), 20),
					),
					dialplan.Busy(),
				},
			},
			{
				Number: "vm-100",
				Actions: []dialplan.Action{
					dialplan.VoiceMail("100", "default", "u"),
					dialplan.VoiceMail("200", "", ""),
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&section); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[test]\n" +
		"\n" +
		"exten => 100,1,Set(ASTGEN_DIAL=&PJSIP/dev-fred.soft@trunk-hub)\n" +
		"same => n,ExecIf($[${DEVICE_STATE(PJSIP/fred)}=NOT_INUSE]?Set(ASTGEN_DIAL=${ASTGEN_DIAL}&PJSIP/fred))\n" +
		"same => n,ExecIf($[${LEN(${ASTGEN_DIAL})}>0]?Dial(${ASTGEN_DIAL:1},20))\n" +
		"same => n,Busy()\n" +
		"\n" +
		"exten => vm-100,1,VoiceMail(100@default,u)\n" +
		"same => n,VoiceMail(200)\n"
	if got := buf.String(); got != want {
		t.Errorf("Encode() =\n%s\nwant:\n%s", got, want)
	}
}

// TestPerServerAppsInvalid checks that the new applications reject text
// that would change the structure of a dialplan line.
func TestPerServerAppsInvalid(t *testing.T) {
	tests := []struct {
		name   string
		action dialplan.Action
	}{
		{name: "comma in mailbox", action: dialplan.VoiceMail("1,2", "default", "")},
		{name: "paren in voicemail options", action: dialplan.VoiceMail("100", "default", "u)")},
		{name: "comma in cat part", action: dialplan.Set(dialplan.Var("X"), dialplan.Cat(dialplan.String("a,b")))},
		{name: "colon in substring name", action: dialplan.DialExpr(dialplan.Substring(dialplan.Var("A:B"), 1), 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := dialplan.Section{Context: "test", Extensions: []dialplan.Extension{{Number: "1", Actions: []dialplan.Action{tt.action}}}}
			var buf bytes.Buffer
			err := astconf.NewEncoder(&buf).Encode(&section)
			var ice astconf.InvalidContentError
			if !errors.As(err, &ice) {
				t.Fatalf("Encode() error = %v, want InvalidContentError", err)
			}
		})
	}
}

// TestPagePredialAndHeaders renders Page's predial option and the pieces of
// the subroutine it calls. A Page without a predial target renders as it
// always has.
func TestPagePredialAndHeaders(t *testing.T) {
	plain := dialplan.Page(dialplan.PJSIP("fred"))
	alerted := dialplan.Page(dialplan.PJSIP("fred"), dialplan.Device{Technology: "Local", Resource: "page-relay-500-den@astgen-page-relay"})
	alerted.Predial = &dialplan.PredialTarget{Context: "astgen-page-alert", Extension: "s", Priority: 1, Args: []string{"intercom"}}
	section := dialplan.Section{
		Context: "test",
		Extensions: []dialplan.Extension{
			{Number: "500", Actions: []dialplan.Action{plain, alerted}},
			{
				Number: "s",
				Actions: []dialplan.Action{
					dialplan.ExecIf(
						dialplan.Equal(dialplan.Channel("channeltype"), dialplan.String("PJSIP")),
						dialplan.Set(dialplan.PJSIPHeader("add", "Alert-Info"), dialplan.Var("ARG1")),
					),
					dialplan.Return(),
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&section); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[test]\n" +
		"\n" +
		"exten => 500,1,Page(PJSIP/fred,is)\n" +
		"same => n,Page(PJSIP/fred&Local/page-relay-500-den@astgen-page-relay,isb(astgen-page-alert^s^1(intercom)))\n" +
		"\n" +
		"exten => s,1,ExecIf($[\"${CHANNEL(channeltype)}\"=\"PJSIP\"]?Set(PJSIP_HEADER(add,Alert-Info)=${ARG1}))\n" +
		"same => n,Return()\n"
	if got := buf.String(); got != want {
		t.Errorf("Encode() =\n%s\nwant:\n%s", got, want)
	}
}

// TestPagePredialInvalid checks that a predial target that would break the
// option is rejected.
func TestPagePredialInvalid(t *testing.T) {
	tests := []struct {
		name   string
		target dialplan.PredialTarget
	}{
		{name: "caret in argument", target: dialplan.PredialTarget{Context: "a", Extension: "s", Priority: 1, Args: []string{"x^y"}}},
		{name: "paren in context", target: dialplan.PredialTarget{Context: "a)", Extension: "s", Priority: 1}},
		{name: "priority zero", target: dialplan.PredialTarget{Context: "a", Extension: "s"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := dialplan.Page(dialplan.PJSIP("fred"))
			target := tt.target
			page.Predial = &target
			section := dialplan.Section{Context: "test", Extensions: []dialplan.Extension{{Number: "1", Actions: []dialplan.Action{page}}}}
			var buf bytes.Buffer
			if err := astconf.NewEncoder(&buf).Encode(&section); err == nil {
				t.Fatalf("Encode() succeeded, want an error:\n%s", buf.String())
			}
		})
	}
}

// TestSubroutine renders Gosub calls with arguments in parentheses after the
// priority, which is how Asterisk passes them as ${ARG1}, ${ARG2}.
func TestSubroutine(t *testing.T) {
	section := dialplan.Section{
		Context: "test",
		Extensions: []dialplan.Extension{
			{
				Number: "100",
				Actions: []dialplan.Action{
					dialplan.Subroutine("pre-answer", "s", 1, dialplan.Var("EXTEN")),
					dialplan.Subroutine("no-answer", "s", 1, dialplan.Var("EXTEN"), dialplan.String("PJSIP/vm-100@trunk-hub")),
					dialplan.Subroutine("cleanup", "s", 1),
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := astconf.NewEncoder(&buf).Encode(&section); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[test]\n" +
		"\n" +
		"exten => 100,1,Gosub(pre-answer,s,1(${EXTEN}))\n" +
		"same => n,Gosub(no-answer,s,1(${EXTEN},PJSIP/vm-100@trunk-hub))\n" +
		"same => n,Gosub(cleanup,s,1)\n"
	if got := buf.String(); got != want {
		t.Errorf("Encode() =\n%s\nwant:\n%s", got, want)
	}
}

// TestSubroutineInvalid checks that a Gosub call that would break the line
// is rejected.
func TestSubroutineInvalid(t *testing.T) {
	tests := []struct {
		name   string
		action dialplan.Action
	}{
		{name: "comma in context", action: dialplan.Subroutine("a,b", "s", 1)},
		{name: "paren in argument", action: dialplan.Subroutine("a", "s", 1, dialplan.String("x)"))},
		{name: "priority zero", action: dialplan.Subroutine("a", "s", 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := dialplan.Section{Context: "test", Extensions: []dialplan.Extension{{Number: "1", Actions: []dialplan.Action{tt.action}}}}
			var buf bytes.Buffer
			if err := astconf.NewEncoder(&buf).Encode(&section); err == nil {
				t.Fatalf("Encode() succeeded, want an error:\n%s", buf.String())
			}
		})
	}
}
