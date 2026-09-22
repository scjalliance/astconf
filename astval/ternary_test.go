package astval_test

import (
	"bytes"
	"testing"

	"github.com/scjalliance/astconf/astval"
)

func TestYesNoNone(t *testing.T) {
	tests := []struct {
		name      string
		value     astval.YesNoNone
		specified bool
		isTrue    bool
		isFalse   bool
		text      string
	}{
		{name: "unspecified", value: astval.YesNoNone(0), text: ""},
		{name: "yes", value: astval.Yes, specified: true, isTrue: true, text: "yes"},
		{name: "no", value: astval.No, specified: true, isFalse: true, text: "no"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Specified(); got != tt.specified {
				t.Errorf("Specified() = %v, want %v", got, tt.specified)
			}
			if got := tt.value.True(); got != tt.isTrue {
				t.Errorf("True() = %v, want %v", got, tt.isTrue)
			}
			if got := tt.value.False(); got != tt.isFalse {
				t.Errorf("False() = %v, want %v", got, tt.isFalse)
			}
			text, err := tt.value.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error: %v", err)
			}
			if string(text) != tt.text {
				t.Errorf("MarshalText() = %q, want %q", text, tt.text)
			}
			var buf bytes.Buffer
			if err := tt.value.MarshalAsteriskSetting(&buf); err != nil {
				t.Fatalf("MarshalAsteriskSetting() error: %v", err)
			}
			if buf.String() != tt.text {
				t.Errorf("MarshalAsteriskSetting() = %q, want %q", buf.String(), tt.text)
			}
		})
	}
}

func TestYesNoNoneUnmarshalText(t *testing.T) {
	tests := []struct {
		text    string
		want    astval.YesNoNone
		wantErr bool
	}{
		{text: "yes", want: astval.Yes},
		{text: "no", want: astval.No},
		{text: "", want: astval.YesNoNone(0)},
		{text: "maybe", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			var got astval.YesNoNone
			err := got.UnmarshalText([]byte(tt.text))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("UnmarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnOffNone(t *testing.T) {
	tests := []struct {
		name      string
		value     astval.OnOffNone
		specified bool
		text      string
	}{
		{name: "unspecified", value: astval.OnOffNone(0), text: ""},
		{name: "on", value: astval.On, specified: true, text: "on"},
		{name: "off", value: astval.Off, specified: true, text: "off"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Specified(); got != tt.specified {
				t.Errorf("Specified() = %v, want %v", got, tt.specified)
			}
			text, err := tt.value.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error: %v", err)
			}
			if string(text) != tt.text {
				t.Errorf("MarshalText() = %q, want %q", text, tt.text)
			}
			var buf bytes.Buffer
			if err := tt.value.MarshalAsteriskSetting(&buf); err != nil {
				t.Fatalf("MarshalAsteriskSetting() error: %v", err)
			}
			if buf.String() != tt.text {
				t.Errorf("MarshalAsteriskSetting() = %q, want %q", buf.String(), tt.text)
			}
		})
	}
}

func TestOneZeroNone(t *testing.T) {
	tests := []struct {
		name      string
		value     astval.OneZeroNone
		specified bool
		text      string
	}{
		{name: "unspecified", value: astval.OneZeroNone(0), text: ""},
		{name: "one", value: astval.One, specified: true, text: "1"},
		{name: "zero", value: astval.Zero, specified: true, text: "0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Specified(); got != tt.specified {
				t.Errorf("Specified() = %v, want %v", got, tt.specified)
			}
			text, err := tt.value.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error: %v", err)
			}
			if string(text) != tt.text {
				t.Errorf("MarshalText() = %q, want %q", text, tt.text)
			}
			var buf bytes.Buffer
			if err := tt.value.MarshalAsteriskSetting(&buf); err != nil {
				t.Fatalf("MarshalAsteriskSetting() error: %v", err)
			}
			if buf.String() != tt.text {
				t.Errorf("MarshalAsteriskSetting() = %q, want %q", buf.String(), tt.text)
			}
		})
	}
}
