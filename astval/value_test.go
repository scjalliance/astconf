package astval_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/scjalliance/astconf/astval"
)

func TestInt(t *testing.T) {
	var unspecified astval.Int
	if unspecified.Specified() {
		t.Error("zero Int reports Specified() = true")
	}
	if got := unspecified.Value(); got != 0 {
		t.Errorf("zero Int Value() = %d, want 0", got)
	}
	var buf bytes.Buffer
	if err := unspecified.MarshalAsteriskSetting(&buf); err != nil {
		t.Fatalf("MarshalAsteriskSetting() error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("zero Int marshaled as %q, want nothing", buf.String())
	}

	tests := []struct {
		value int
		want  string
	}{
		{value: 0, want: "0"},
		{value: 5, want: "5"},
		{value: -12, want: "-12"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			i := astval.NewInt(tt.value)
			if !i.Specified() {
				t.Error("NewInt() reports Specified() = false")
			}
			if got := i.Value(); got != tt.value {
				t.Errorf("Value() = %d, want %d", got, tt.value)
			}
			var buf bytes.Buffer
			if err := i.MarshalAsteriskSetting(&buf); err != nil {
				t.Fatalf("MarshalAsteriskSetting() error: %v", err)
			}
			if buf.String() != tt.want {
				t.Errorf("MarshalAsteriskSetting() = %q, want %q", buf.String(), tt.want)
			}

			var set astval.Int
			set.Set(tt.value)
			if set != i {
				t.Errorf("Set(%d) = %v, want %v", tt.value, set, i)
			}
		})
	}
}

func TestString(t *testing.T) {
	var unspecified astval.String
	if unspecified.Specified() {
		t.Error("zero String reports Specified() = true")
	}
	var buf bytes.Buffer
	if err := unspecified.MarshalAsteriskSetting(&buf); err != nil {
		t.Fatalf("MarshalAsteriskSetting() error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("zero String marshaled as %q, want nothing", buf.String())
	}

	// An explicitly specified empty string is still specified
	empty := astval.NewString("")
	if !empty.Specified() {
		t.Error("NewString(\"\") reports Specified() = false")
	}

	s := astval.NewString("hello")
	if got := s.Value(); got != "hello" {
		t.Errorf("Value() = %q, want %q", got, "hello")
	}
	buf.Reset()
	if err := s.MarshalAsteriskSetting(&buf); err != nil {
		t.Fatalf("MarshalAsteriskSetting() error: %v", err)
	}
	if buf.String() != "hello" {
		t.Errorf("MarshalAsteriskSetting() = %q, want %q", buf.String(), "hello")
	}

	var set astval.String
	set.Set("hello")
	if set != s {
		t.Errorf("Set() = %v, want %v", set, s)
	}
}

func TestSeconds(t *testing.T) {
	var unspecified astval.Seconds
	if unspecified.Specified() {
		t.Error("zero Seconds reports Specified() = true")
	}
	var buf bytes.Buffer
	if err := unspecified.MarshalAsteriskSetting(&buf); err != nil {
		t.Fatalf("MarshalAsteriskSetting() error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("zero Seconds marshaled as %q, want nothing", buf.String())
	}

	tests := []struct {
		value time.Duration
		want  string
	}{
		{value: 0, want: "0"},
		{value: 30 * time.Second, want: "30"},
		{value: 2 * time.Minute, want: "120"},
		{value: 1500 * time.Millisecond, want: "1"}, // Truncated to whole seconds
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			s := astval.NewSeconds(tt.value)
			if !s.Specified() {
				t.Error("NewSeconds() reports Specified() = false")
			}
			if got := s.Value(); got != tt.value {
				t.Errorf("Value() = %v, want %v", got, tt.value)
			}
			var buf bytes.Buffer
			if err := s.MarshalAsteriskSetting(&buf); err != nil {
				t.Fatalf("MarshalAsteriskSetting() error: %v", err)
			}
			if buf.String() != tt.want {
				t.Errorf("MarshalAsteriskSetting() = %q, want %q", buf.String(), tt.want)
			}

			var set astval.Seconds
			set.Set(tt.value)
			if set != s {
				t.Errorf("Set() = %v, want %v", set, s)
			}
		})
	}
}

func TestVar(t *testing.T) {
	v := astval.NewVar("AREACODE", "360")
	if got := v.String(); got != "AREACODE=360" {
		t.Errorf("String() = %q, want %q", got, "AREACODE=360")
	}
	text, err := v.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() error: %v", err)
	}
	if string(text) != "AREACODE=360" {
		t.Errorf("MarshalText() = %q, want %q", text, "AREACODE=360")
	}

	tests := []struct {
		text string
		want astval.Var
	}{
		{text: "AREACODE=360", want: astval.Var{Name: "AREACODE", Value: "360"}},
		{text: "NAME=a=b", want: astval.Var{Name: "NAME", Value: "a=b"}}, // Only the first = separates
		{text: "FLAG", want: astval.Var{Name: "FLAG"}},
		{text: "EMPTY=", want: astval.Var{Name: "EMPTY", Value: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			var got astval.Var
			if err := got.UnmarshalText([]byte(tt.text)); err != nil {
				t.Fatalf("UnmarshalText() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("UnmarshalText() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
