package voicemail_test

import (
	"reflect"
	"testing"

	"github.com/scjalliance/astconf"
	"github.com/scjalliance/astconf/voicemail"
)

func TestBoxOptions(t *testing.T) {
	tests := []struct {
		name string
		box  voicemail.Box
		want []string
	}{
		{
			name: "defaults",
			box:  voicemail.Box{},
			want: []string{"attach=yes", "delete=no"},
		},
		{
			name: "all options",
			box: voicemail.Box{
				Timezone:     "pacific",
				Locale:       "en_US",
				SendToPager:  true,
				Format:       "wav49",
				SayCallerID:  true,
				SkipEnvelope: true,
				EmailOnly:    true,
			},
			want: []string{"tz=pacific", "locale=en_US", "attachfmt=wav49", "saycid=yes", "envelope=no", "delete=yes"},
		},
		{
			name: "separators stripped from values",
			box:  voicemail.Box{Timezone: "a|b,c"},
			want: []string{"tz=abc", "attach=yes", "delete=no"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.box.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSectionMarshal(t *testing.T) {
	section := voicemail.Section{
		Context: "default",
		Mailboxes: []voicemail.Box{
			{
				Extension:            "100",
				Name:                 "Fred Flintstone",
				Password:             "1234",
				PasswordIsChangeable: true,
				EmailAddresses:       []string{"fred@example.com", "wilma@example.com"},
				Timezone:             "pacific",
			},
			{
				Extension: "101",
				Name:      "Barney Rubble",
				Password:  "4321",
			},
			{
				// Trailing empty components are dropped
				Extension: "102",
				Password:  "0000",
				EmailOnly: true,
			},
		},
	}
	got, err := astconf.Marshal(&section)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	want := "[default]\n" +
		"100 => 1234,Fred Flintstone,fred@example.com|wilma@example.com,,tz=pacific|attach=yes|delete=no\n" +
		"101 => -4321,Barney Rubble,,,attach=yes|delete=no\n" +
		"102 => -0000,,,,attach=yes|delete=yes\n"
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
