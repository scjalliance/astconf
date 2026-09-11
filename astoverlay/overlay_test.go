package astoverlay_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/astoverlay"
	"github.com/scjalliance/astconf/astval"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		from, to string
		want     string
	}{
		{name: "from replaces to", from: "a", to: "b", want: "a"},
		{name: "empty from keeps to", from: "", to: "b", want: "b"},
		{name: "both empty", from: "", to: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astoverlay.String(&tt.from, &tt.to)
			if tt.to != tt.want {
				t.Errorf("to = %q, want %q", tt.to, tt.want)
			}
		})
	}
}

func TestSectionName(t *testing.T) {
	from, to := astconf.SectionName("a"), astconf.SectionName("b")
	astoverlay.SectionName(&from, &to)
	if to != "a" {
		t.Errorf("to = %q, want %q", to, "a")
	}
	from, to = "", "b"
	astoverlay.SectionName(&from, &to)
	if to != "b" {
		t.Errorf("to = %q, want %q", to, "b")
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name     string
		from, to int
		want     int
	}{
		{name: "from replaces to", from: 1, to: 2, want: 1},
		{name: "zero from keeps to", from: 0, to: 2, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astoverlay.Int(&tt.from, &tt.to)
			if tt.to != tt.want {
				t.Errorf("to = %d, want %d", tt.to, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	from, to := time.Second, time.Minute
	astoverlay.Duration(&from, &to)
	if to != time.Second {
		t.Errorf("to = %v, want %v", to, time.Second)
	}
	from, to = 0, time.Minute
	astoverlay.Duration(&from, &to)
	if to != time.Minute {
		t.Errorf("to = %v, want %v", to, time.Minute)
	}
}

func TestStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		from, to []string
		want     []string
	}{
		{name: "from replaces to", from: []string{"a"}, to: []string{"b", "c"}, want: []string{"a"}},
		{name: "nil from keeps to", from: nil, to: []string{"b"}, want: []string{"b"}},
		{name: "empty non-nil from replaces to", from: []string{}, to: []string{"b"}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astoverlay.StringSlice(&tt.from, &tt.to)
			if !reflect.DeepEqual(tt.to, tt.want) {
				t.Errorf("to = %v, want %v", tt.to, tt.want)
			}
		})
	}
}

func TestAstVarSlice(t *testing.T) {
	a := []astval.Var{astval.NewVar("A", "1")}
	b := []astval.Var{astval.NewVar("B", "2")}
	astoverlay.AstVarSlice(&a, &b)
	if !reflect.DeepEqual(b, []astval.Var{astval.NewVar("A", "1")}) {
		t.Errorf("to = %v, want %v", b, a)
	}
	var none []astval.Var
	astoverlay.AstVarSlice(&none, &b)
	if !reflect.DeepEqual(b, []astval.Var{astval.NewVar("A", "1")}) {
		t.Errorf("nil from changed to = %v", b)
	}
}

func TestAstYesNoNone(t *testing.T) {
	from, to := astval.Yes, astval.No
	astoverlay.AstYesNoNone(&from, &to)
	if to != astval.Yes {
		t.Errorf("to = %v, want Yes", to)
	}
	from, to = astval.YesNoNone(0), astval.No
	astoverlay.AstYesNoNone(&from, &to)
	if to != astval.No {
		t.Errorf("unspecified from changed to = %v", to)
	}
}

func TestAstInt(t *testing.T) {
	from, to := astval.NewInt(1), astval.NewInt(2)
	astoverlay.AstInt(&from, &to)
	if to.Value() != 1 {
		t.Errorf("to = %d, want 1", to.Value())
	}
	from, to = astval.Int{}, astval.NewInt(2)
	astoverlay.AstInt(&from, &to)
	if to.Value() != 2 {
		t.Errorf("unspecified from changed to = %d", to.Value())
	}
	// A specified zero still overlays
	from, to = astval.NewInt(0), astval.NewInt(2)
	astoverlay.AstInt(&from, &to)
	if to.Value() != 0 {
		t.Errorf("specified zero did not overlay, to = %d", to.Value())
	}
}

func TestAstString(t *testing.T) {
	from, to := astval.NewString("a"), astval.NewString("b")
	astoverlay.AstString(&from, &to)
	if to.Value() != "a" {
		t.Errorf("to = %q, want %q", to.Value(), "a")
	}
	from, to = astval.String{}, astval.NewString("b")
	astoverlay.AstString(&from, &to)
	if to.Value() != "b" {
		t.Errorf("unspecified from changed to = %q", to.Value())
	}
}

func TestAstSeconds(t *testing.T) {
	from, to := astval.NewSeconds(time.Second), astval.NewSeconds(time.Minute)
	astoverlay.AstSeconds(&from, &to)
	if to.Value() != time.Second {
		t.Errorf("to = %v, want %v", to.Value(), time.Second)
	}
	from, to = astval.Seconds{}, astval.NewSeconds(time.Minute)
	astoverlay.AstSeconds(&from, &to)
	if to.Value() != time.Minute {
		t.Errorf("unspecified from changed to = %v", to.Value())
	}
}
