package astconf_test

import (
	"bytes"
	"testing"

	"github.com/scjalliance/astconf"
)

func TestPrinter(t *testing.T) {
	var buf bytes.Buffer
	p := astconf.NewPrinter(&buf)

	steps := []struct {
		name string
		call func() error
		want string
	}{
		{name: "include", call: func() error { return p.Include("pjsip.d/*.conf") }, want: "#include pjsip.d/*.conf\n"},
		{name: "comment", call: func() error { return p.Comment("hello") }, want: "; hello\n"},
		{name: "first section", call: func() error { return p.Section("general") }, want: "[general]\n"},
		{name: "setting", call: func() error { return p.Setting("key", "value") }, want: "key = value\n"},
		{name: "object", call: func() error { return p.Object("exten", "100,1,Noop()") }, want: "exten => 100,1,Noop()\n"},
		{name: "second section", call: func() error { return p.Section("phone", "base") }, want: "\n[phone](base)\n"},
		{name: "templated section", call: func() error { return p.Section("other", "a", "b") }, want: "\n[other](a,b)\n"},
		{name: "start", call: func() error { return p.Start("name", " = ") }, want: "name = "},
	}
	for _, step := range steps {
		buf.Reset()
		if err := step.call(); err != nil {
			t.Fatalf("%s: error %v", step.name, err)
		}
		if buf.String() != step.want {
			t.Errorf("%s = %q, want %q", step.name, buf.String(), step.want)
		}
	}
}

func TestPrinterBreak(t *testing.T) {
	var buf bytes.Buffer
	p := astconf.NewPrinter(&buf)
	p.Break()
	if buf.Len() != 0 {
		t.Errorf("Break() before start wrote %q", buf.String())
	}
	if err := p.Setting("a", "1"); err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	p.Break()
	if buf.String() != "\n" {
		t.Errorf("Break() after start wrote %q, want newline", buf.String())
	}
}

func TestPrinterAlignment(t *testing.T) {
	tests := []struct {
		name      string
		alignment astconf.Alignment
		want      string
	}{
		{name: "none", alignment: astconf.AlignmentNone, want: "ab = 1\n"},
		{name: "left", alignment: astconf.AlignmentLeft, want: "ab    = 1\n"},
		{name: "right", alignment: astconf.AlignmentRight, want: "   ab = 1\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := astconf.Printer{Writer: &buf, Indent: 5, Alignment: tt.alignment}
			if err := p.Setting("ab", "1"); err != nil {
				t.Fatal(err)
			}
			if buf.String() != tt.want {
				t.Errorf("Setting() = %q, want %q", buf.String(), tt.want)
			}
		})
	}
}
