package dialplan_test

import (
	"errors"
	"testing"

	"github.com/scjalliance/astconf"
	. "github.com/scjalliance/astconf/dialplan"
)

// marshalExtension renders a single extension and returns the error.
func marshalExtension(number string, actions ...Action) error {
	section := Section{
		Context:    "test",
		Extensions: []Extension{{Number: number, Actions: actions}},
	}
	_, err := astconf.Marshal(&section)
	return err
}

func TestInvalidDialplanContent(t *testing.T) {
	tests := []struct {
		name   string
		number string
		action Action
	}{
		{name: "comma in extension number", number: "100,2", action: Noop("")},
		{name: "paren in extension number", number: "100)", action: Noop("")},
		{name: "variable in extension number", number: "${EXTEN}", action: Noop("")},
		{name: "close paren in noop", number: "100", action: Noop("Call Fred)")},
		{name: "open paren in noop", number: "100", action: Noop("Call (Fred")},
		{name: "variable in noop", number: "100", action: Noop("Call ${SIPDOMAIN}")},
		{name: "expression in noop", number: "100", action: Noop("Call $[1]")},
		{name: "comma in noop", number: "100", action: Noop("Flintstone, Fred")},
		{name: "comma in sip header content", number: "100", action: SIPAddHeader("Alert-Info", "a,b")},
		{name: "paren in sip header name", number: "100", action: SIPAddHeader("Alert)", "a")},
		{name: "paren in playback file", number: "100", action: Playback("beep)")},
		{name: "ampersand in playback file", number: "100", action: Playback("beep&boop")},
		{name: "comma in dial device", number: "100", action: Dial(SIP("a,b"), 20)},
		{name: "ampersand in dial device", number: "100", action: Dial(SIP("a&b"), 20)},
		{name: "slash in dial device", number: "100", action: Dial(SIP("a/b"), 20)},
		{name: "variable in dial device", number: "100", action: DialMany([]Device{SIP("a"), SIP("${X}")}, 20)},
		{name: "paren in page device", number: "100", action: Page(SIP("a)"))},
		{name: "paren in page announcement", number: "100", action: PageApp{Recipients: []Device{SIP("a")}, Announcement: "x)"}},
		{name: "comma in hangup cause", number: "100", action: HangupApp{Cause: "16,17"}},
		{name: "comma in gosub argument", number: "100", action: Gosub("ctx", "s", 1, "a,b")},
		{name: "paren in gosub context", number: "100", action: Gosub("ctx)", "s", 1)},
		{name: "comma in macro name", number: "100", action: Macro("pre,answer")},
		{name: "bracket in macro argument", number: "100", action: Macro("pre-answer", String("a]"))},
		{name: "nested noop in execif", number: "100", action: ExecIf(Equal(Int(1), Int(1)), Noop("x)"))},
		{name: "nested noop in execif else", number: "100", action: ExecIfElse(Equal(Int(1), Int(1)), Noop("ok"), Noop("x)"))},
		{name: "string in execif comparison", number: "100", action: ExecIf(Equal(DeviceState(SIP("a")), String("NOT_INUSE]")), Noop("ok"))},
		{name: "device in execif comparison", number: "100", action: ExecIf(Equal(DeviceState(SIP("a,b")), String("NOT_INUSE")), Noop("ok"))},
		{name: "device in or comparison", number: "100", action: ExecIf(Or(Equal(DeviceState(SIP("a")), String("x")), Equal(DeviceState(SIP("b)")), String("x"))), Noop("ok"))},
		{name: "string in set value", number: "100", action: Set(Var("X"), String("a$b"))},
		{name: "string in if function", number: "100", action: Set(Var("X"), IfElse(Equal(Int(1), Int(1)), String("a"), String("b)")))},
		{name: "presence state provider", number: "100", action: ExecIf(Equal(PresenceState("hint:100)", "subtype"), String("x")), Noop("ok"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := marshalExtension(tt.number, tt.action)
			var ice astconf.InvalidContentError
			if !errors.As(err, &ice) {
				t.Fatalf("Marshal() error = %v, want InvalidContentError", err)
			}
		})
	}
}

func TestValidDialplanContent(t *testing.T) {
	tests := []struct {
		name   string
		number string
		action Action
	}{
		{name: "pattern with brackets", number: "_[2-9]XX", action: Noop("ok")},
		{name: "quoted comparison nested in a quoted comparison", number: "100", action: ExecIf(Equal(CallerID("a"), Equal(CallerID("b"), String("x"))), Noop("ok"))},
		{name: "pattern with wildcard", number: "_1XX.", action: Noop("ok")},
		{name: "special extension", number: "s", action: Noop("ok")},
		{name: "spaces and punctuation in noop", number: "100", action: Noop("Call Fred Flintstone - ext. 100")},
		{name: "device with dot", number: "100", action: Dial(SIP("fred.flintstone"), 20)},
		{name: "nested execif", number: "100", action: ExecIf(Equal(DeviceState(SIP("a")), String("NOT_INUSE")), Dial(SIP("a"), 20))},
		{name: "or of device states", number: "100", action: ExecIf(Or(Equal(DeviceState(SIP("a")), String("x")), Equal(DeviceState(SIP("b")), String("y"))), Noop("ok"))},
		{name: "set from if function", number: "100", action: Set(CallerID("name"), IfElse(GreaterThan(Len(CallerID("num")), Int(0)), CallerID("name"), String("Intercom")))},
		{name: "macro with variable", number: "100", action: Macro("pre-answer", Var("EXTEN"))},
		{name: "gosub", number: "100", action: Gosub("ctx", "s", 1, "arg")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := marshalExtension(tt.number, tt.action); err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
		})
	}
}

func TestDeviceValidate(t *testing.T) {
	if err := SIP("fred.flintstone").Validate(); err != nil {
		t.Errorf("Validate() = %v for a valid device", err)
	}
	for _, resource := range []string{"a,b", "a(b", "a)b", "a$b", "a[b", "a]b", "a&b", "a/b"} {
		if err := SIP(resource).Validate(); err == nil {
			t.Errorf("Validate() = nil for resource %q", resource)
		}
	}
}

// TestInvalidDialplanOperators covers the operator and technology fields,
// which are rendered into a line but were not validated alongside the
// operands and resources next to them.
func TestInvalidDialplanOperators(t *testing.T) {
	tests := []struct {
		name   string
		action Action
	}{
		{name: "quote in binary operator", action: ExecIf(BinaryOp{E1: Int(1), E2: Int(1), Operator: "=\""}, Noop("ok"))},
		{name: "paren in binary operator", action: ExecIf(BinaryOp{E1: Int(1), E2: Int(1), Operator: "=)"}, Noop("ok"))},
		{name: "paren in logical operator", action: ExecIf(MultiOp{Expressions: []Expression{Equal(Int(1), Int(1)), Equal(Int(2), Int(2))}, Operator: "|("}, Noop("ok"))},
		{name: "comma in dial device technology", action: Dial(Device{Technology: "SIP,Local", Resource: "100"}, 20)},
		{name: "slash in dial device technology", action: Dial(Device{Technology: "SIP/x", Resource: "100"}, 20)},
		{name: "quote in operand of a quoted comparison", action: ExecIf(Equal(CallerID("num"), String(`a"b`)), Noop("ok"))},
		{name: "quote in operand against presence state", action: ExecIf(Equal(PresenceState("hint:100", "subtype"), String(`a"b`)), Noop("ok"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := marshalExtension("100", tt.action)
			var ice astconf.InvalidContentError
			if !errors.As(err, &ice) {
				t.Fatalf("Marshal() error = %v, want InvalidContentError", err)
			}
		})
	}
}

// TestDialplanStructuralValidation covers the two cases that are wrong by
// shape rather than by character. Both previously reached rendering, and the
// operand count panicked there.
func TestDialplanStructuralValidation(t *testing.T) {
	tests := []struct {
		name   string
		action Action
	}{
		{name: "gosub priority below one", action: Gosub("ctx", "s", 0)},
		{name: "empty binary operator", action: ExecIf(BinaryOp{E1: Int(1), E2: Int(1)}, Noop("ok"))},
		{name: "empty logical operator", action: ExecIf(MultiOp{Expressions: []Expression{Equal(Int(1), Int(1)), Equal(Int(2), Int(2))}}, Noop("ok"))},
		{name: "negative gosub priority", action: Gosub("ctx", "s", -1)},
		{name: "logical operator with no operands", action: ExecIf(MultiOp{Operator: "|"}, Noop("ok"))},
		{name: "logical operator with one operand", action: ExecIf(MultiOp{Expressions: []Expression{Equal(Int(1), Int(1))}, Operator: "|"}, Noop("ok"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := marshalExtension("100", tt.action); err == nil {
				t.Fatal("Marshal() error = nil, want an error")
			}
		})
	}
}
