package astconf_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/scjalliance/astconf"
)

// scalars exercises the built-in encoders for simple kinds.
type scalars struct {
	Str   string `astconf:"str"`
	Int   int    `astconf:"int"`
	Uint  uint   `astconf:"uint"`
	True  bool   `astconf:"true"`
	False bool   `astconf:"false"`
}

// omitted exercises the omitempty and "-" tag options.
type omitted struct {
	Kept     string   `astconf:"kept"`
	Str      string   `astconf:"str,omitempty"`
	Int      int      `astconf:"int,omitempty"`
	Bool     bool     `astconf:"bool,omitempty"`
	Slice    []string `astconf:"slice,omitempty"`
	Hidden   string   `astconf:"-"`
	Untagged string
}

// multi exercises multi-valued fields.
type multi struct {
	Allow  []string `astconf:"allow"`
	Codecs []string `astconf:"codecs,commaseparated"`
	Exten  []string `astconf:"exten,object"`
}

// inner is a struct without a section name of its own.
type inner struct {
	Value string `astconf:"value"`
}

// nested exercises a struct field that is not a section.
type nested struct {
	Name  string `astconf:"name"`
	Inner inner
}

// named exercises the SectionName helper type as a regular field. It must
// not be embedded, because embedding promotes its no-op MarshalAsterisk
// onto the outer struct.
type named struct {
	Section astconf.SectionName
	Value   string `astconf:"value"`
}

// templated exercises SectionTemplater.
type templated struct {
	Name  string `astconf:"-"`
	Value string `astconf:"value"`
}

func (t templated) SectionName() string        { return t.Name }
func (t templated) SectionTemplates() []string { return []string{"base", "extra"} }

// preambled exercises PreambleMarshaler on a pointer receiver.
type preambled struct {
	Name  string `astconf:"-"`
	Value string `astconf:"value"`
}

func (p *preambled) SectionName() string { return p.Name }
func (p *preambled) MarshalAsteriskPreamble(e *astconf.Encoder) error {
	return e.Printer().Setting("type", "preambled")
}

// upper is a TextMarshaler.
type upper string

func (u upper) MarshalText() ([]byte, error) {
	return bytes.ToUpper([]byte(u)), nil
}

// textual exercises TextMarshaler fields.
type textual struct {
	Name upper `astconf:"name"`
}

// quiet is a SettingMarshaler that writes nothing.
type quiet struct{}

func (quiet) MarshalAsteriskSetting(w io.Writer) error { return nil }

// withQuiet exercises a SettingMarshaler that writes nothing. The field
// name should be omitted along with the value.
type withQuiet struct {
	Before string `astconf:"before"`
	Quiet  quiet  `astconf:"quiet"`
	After  string `astconf:"after"`
}

// failing is a Marshaler that always fails.
type failing struct{}

var errFailing = errors.New("failing marshaler")

func (failing) MarshalAsterisk(e *astconf.Encoder) error { return errFailing }

// unsupported exercises an unsupported kind.
type unsupported struct {
	M map[string]string `astconf:"m"`
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "scalars",
			value: scalars{Str: "text", Int: -3, Uint: 4, True: true, False: false},
			want: "str = text\n" +
				"int = -3\n" +
				"uint = 4\n" +
				"true = yes\n" +
				"false = false\n", // Note: false is written as "false", not "no"
		},
		{
			name:  "omitempty with empty values",
			value: omitted{Kept: "", Hidden: "hidden", Untagged: "untagged"},
			want: "kept = \n" +
				"Untagged = untagged\n",
		},
		{
			name:  "omitempty with values",
			value: omitted{Str: "s", Int: 1, Bool: true, Slice: []string{"a"}},
			want: "kept = \n" +
				"str = s\n" +
				"int = 1\n" +
				"bool = yes\n" +
				"slice = a\n" +
				"Untagged = \n",
		},
		{
			name:  "multi-valued fields",
			value: multi{Allow: []string{"ulaw", "g722"}, Codecs: []string{"ulaw", "g722"}, Exten: []string{"100,1,Noop()", "101,1,Noop()"}},
			want: "allow = ulaw\n" +
				"allow = g722\n" +
				"codecs = ulaw,g722\n" +
				"exten => 100,1,Noop()\n" +
				"exten => 101,1,Noop()\n",
		},
		{
			name:  "empty multi-valued fields",
			value: multi{},
			want:  "",
		},
		{
			name:  "nested struct without section",
			value: nested{Name: "outer", Inner: inner{Value: "inner"}},
			want: "name = outer\n" +
				"value = inner\n",
		},
		{
			name:  "section name helper",
			value: named{Section: "general", Value: "v"},
			want: "[general]\n" +
				"value = v\n",
		},
		{
			name:  "section templates",
			value: templated{Name: "phone", Value: "v"},
			want: "[phone](base,extra)\n" +
				"value = v\n",
		},
		{
			name:  "preamble on pointer receiver",
			value: &preambled{Name: "p", Value: "v"},
			want: "[p]\n" +
				"type = preambled\n" +
				"value = v\n",
		},
		{
			name:  "text marshaler",
			value: textual{Name: "shout"},
			want:  "name = SHOUT\n",
		},
		{
			name:  "setting marshaler that writes nothing",
			value: withQuiet{Before: "b", After: "a"},
			want: "before = b\n" +
				"after = a\n",
		},
		{
			name:  "slice of sections",
			value: []templated{{Name: "a", Value: "1"}, {Name: "b", Value: "2"}},
			want: "[a](base,extra)\n" +
				"value = 1\n" +
				"\n" +
				"[b](base,extra)\n" +
				"value = 2\n",
		},
		{
			name:  "nil pointer",
			value: (*preambled)(nil),
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := astconf.Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal() =\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestMarshalErrors(t *testing.T) {
	t.Run("unsupported type", func(t *testing.T) {
		_, err := astconf.Marshal(unsupported{M: map[string]string{"a": "b"}})
		var ute *astconf.UnsupportedTypeError
		if !errors.As(err, &ute) {
			t.Fatalf("Marshal() error = %v, want UnsupportedTypeError", err)
		}
	})
	t.Run("invalid value", func(t *testing.T) {
		_, err := astconf.Marshal(nil)
		var ive *astconf.InvalidValueError
		if !errors.As(err, &ive) {
			t.Fatalf("Marshal() error = %v, want InvalidValueError", err)
		}
	})
	t.Run("marshaler error", func(t *testing.T) {
		_, err := astconf.Marshal(failing{})
		var me *astconf.MarshalerError
		if !errors.As(err, &me) {
			t.Fatalf("Marshal() error = %v, want MarshalerError", err)
		}
		if me.Err != errFailing {
			t.Errorf("MarshalerError.Err = %v, want %v", me.Err, errFailing)
		}
	})
}

func TestAlignment(t *testing.T) {
	value := struct {
		A   string `astconf:"a"`
		Bcd string `astconf:"bcd"`
	}{A: "1", Bcd: "2"}
	tests := []struct {
		name string
		opts []astconf.EncOpt
		want string
	}{
		{name: "default", want: "a = 1\nbcd = 2\n"},
		{name: "unaligned", opts: []astconf.EncOpt{astconf.Unaligned}, want: "a = 1\nbcd = 2\n"},
		{name: "left", opts: []astconf.EncOpt{astconf.AlignLeft}, want: "a   = 1\nbcd = 2\n"},
		{name: "right", opts: []astconf.EncOpt{astconf.AlignRight}, want: "  a = 1\nbcd = 2\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := astconf.NewEncoder(&buf, tt.opts...).Encode(value); err != nil {
				t.Fatalf("Encode() error: %v", err)
			}
			if buf.String() != tt.want {
				t.Errorf("Encode() =\n%s\nwant:\n%s", buf.String(), tt.want)
			}
		})
	}
}

func TestMarshalTo(t *testing.T) {
	var buf bytes.Buffer
	if err := astconf.MarshalTo(scalars{Str: "x"}, &buf); err != nil {
		t.Fatalf("MarshalTo() error: %v", err)
	}
	want := "str = x\nint = 0\nuint = 0\ntrue = false\nfalse = false\n"
	if buf.String() != want {
		t.Errorf("MarshalTo() =\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestEncoderMultipleSections(t *testing.T) {
	var buf bytes.Buffer
	e := astconf.NewEncoder(&buf)
	if err := e.Encode(templated{Name: "a", Value: "1"}); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	if err := e.Encode(templated{Name: "b", Value: "2"}); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}
	want := "[a](base,extra)\nvalue = 1\n\n[b](base,extra)\nvalue = 2\n"
	if buf.String() != want {
		t.Errorf("Encode() twice =\n%s\nwant:\n%s", buf.String(), want)
	}
}

// pointers exercises pointer fields to scalar values.
type pointers struct {
	Str    *string `astconf:"str"`
	Int    *int    `astconf:"int"`
	OptInt *int    `astconf:"optint,omitempty"`
	Block  *inner
}

func TestMarshalPointerFields(t *testing.T) {
	str := "text"
	zero := 0
	seven := 7
	tests := []struct {
		name  string
		value pointers
		want  string
	}{
		{
			name:  "nil pointers are omitted",
			value: pointers{},
			want:  "",
		},
		{
			name:  "pointers are dereferenced",
			value: pointers{Str: &str, Int: &seven, OptInt: &seven, Block: &inner{Value: "v"}},
			want: "str = text\n" +
				"int = 7\n" +
				"optint = 7\n" +
				"value = v\n",
		},
		{
			name:  "omitempty keeps a non-nil pointer to a zero value",
			value: pointers{Int: &zero, OptInt: &zero},
			want:  "int = 0\noptint = 0\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := astconf.Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal() =\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

// TestMarshalSliceOfPointers covers a slice of pointers to section types,
// which previously broke the marshaler.
func TestMarshalSliceOfPointers(t *testing.T) {
	value := struct {
		Name  string `astconf:"name"`
		Items []*preambled
	}{
		Name:  "n",
		Items: []*preambled{{Name: "a", Value: "1"}, nil, {Name: "b", Value: "2"}},
	}
	want := "name = n\n" +
		"\n" +
		"[a]\n" +
		"type = preambled\n" +
		"value = 1\n" +
		"\n" +
		"[b]\n" +
		"type = preambled\n" +
		"value = 2\n"
	got, err := astconf.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}
	if string(got) != want {
		t.Errorf("Marshal() =\n%s\nwant:\n%s", got, want)
	}
}
