package astval_test

import (
	"testing"

	"github.com/scjalliance/astconf/astval"
)

func TestYesNo(t *testing.T) {
	if got := astval.YesNo(true).String(); got != "yes" {
		t.Errorf("YesNo(true).String() = %q, want %q", got, "yes")
	}
	if got := astval.YesNo(false).String(); got != "no" {
		t.Errorf("YesNo(false).String() = %q, want %q", got, "no")
	}
	tests := []struct {
		text    string
		want    astval.YesNo
		wantErr bool
	}{
		{text: "yes", want: true},
		{text: "no", want: false},
		{text: "", wantErr: true},
		{text: "true", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			var got astval.YesNo
			err := got.UnmarshalText([]byte(tt.text))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("UnmarshalText() = %v, want %v", got, tt.want)
			}
			text, err := got.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error: %v", err)
			}
			if string(text) != tt.text {
				t.Errorf("MarshalText() = %q, want %q", text, tt.text)
			}
		})
	}
}

func TestOnOff(t *testing.T) {
	if got := astval.OnOff(true).String(); got != "on" {
		t.Errorf("OnOff(true).String() = %q, want %q", got, "on")
	}
	if got := astval.OnOff(false).String(); got != "off" {
		t.Errorf("OnOff(false).String() = %q, want %q", got, "off")
	}
	tests := []struct {
		text    string
		want    astval.OnOff
		wantErr bool
	}{
		{text: "on", want: true},
		{text: "off", want: false},
		{text: "yes", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			var got astval.OnOff
			err := got.UnmarshalText([]byte(tt.text))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("UnmarshalText() = %v, want %v", got, tt.want)
			}
			text, err := got.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error: %v", err)
			}
			if string(text) != tt.text {
				t.Errorf("MarshalText() = %q, want %q", text, tt.text)
			}
		})
	}
}

func TestOneZero(t *testing.T) {
	if got := astval.OneZero(true).String(); got != "1" {
		t.Errorf("OneZero(true).String() = %q, want %q", got, "1")
	}
	if got := astval.OneZero(false).String(); got != "0" {
		t.Errorf("OneZero(false).String() = %q, want %q", got, "0")
	}
	tests := []struct {
		text    string
		want    astval.OneZero
		wantErr bool
	}{
		{text: "1", want: true},
		{text: "0", want: false},
		{text: "2", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			var got astval.OneZero
			err := got.UnmarshalText([]byte(tt.text))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("UnmarshalText() = %v, want %v", got, tt.want)
			}
			text, err := got.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText() error: %v", err)
			}
			if string(text) != tt.text {
				t.Errorf("MarshalText() = %q, want %q", text, tt.text)
			}
		})
	}
}
