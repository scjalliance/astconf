package astconf

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Config file format references:
//
// https://wiki.asterisk.org/wiki/display/AST/Config+File+Format
// https://github.com/asterisk/asterisk/blob/master/utils/extconf.c

var (
	sectionStart     = []byte("[")
	sectionEnd       = []byte("]")
	templateStart    = []byte("(")
	templateEnd      = []byte(")")
	includeStart     = []byte("#include ")
	comma            = []byte(",")
	commentStart     = []byte("; ")
	settingSeparator = []byte(" = ")
	objectSeparator  = []byte(" => ")
	newline          = []byte("\n")
)

// Characters that cannot appear in each component of a configuration file.
// Newlines would start a new line, brackets and parentheses would be read
// as section headers, and an equals sign would end a setting name early.
const (
	invalidNameChars     = "\r\n=[]"
	invalidValueChars    = "\r\n"
	invalidSectionChars  = "\r\n[]()"
	invalidTemplateChars = "\r\n[](),"
	invalidCommentChars  = "\r\n"
)

// Printer is capable of printing asterisk configuration data to an
// underlying io.Writer.
//
// Printer performs no buffering on its own. It is recommended that the
// underlying writer be buffered.
//
// To marshal custom types via reflection, use Encoder or Marshal.
type Printer struct {
	io.Writer
	Indent    int
	Alignment Alignment
	Started   bool // If true a newline will be printed before each section
}

// NewPrinter returns a new printer that writes to w.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{
		Writer: w,
	}
}

// Include will print a file include construct to p.Writer for the given path.
func (p *Printer) Include(path string) error {
	if err := errorIfAny("include path", path, invalidCommentChars); err != nil {
		return err
	}
	wg := writegroup{Writer: p.Writer}
	wg.Write(includeStart)
	wg.Write([]byte(path))
	wg.Write(newline)
	return wg.Err()
}

//func (p *Printer) Template(name string, templates ...string) error {
//}

// Comment writes a single line comment to p.Writer.
func (p *Printer) Comment(comment string) error {
	if err := errorIfAny("comment", comment, invalidCommentChars); err != nil {
		return err
	}
	wg := writegroup{Writer: p.Writer}
	wg.Write(commentStart)
	wg.Write([]byte(comment))
	wg.Write(newline)
	return wg.Err()
}

// Break causes the printer to insert a newline if the printer has started.
func (p *Printer) Break() {
	if p.Started {
		p.Writer.Write(newline)
	}
}

// Start begins a new field by writing its name and optional
// separator.
//
// An InvalidContentError will be returned if the name contains an
// invalid character.
func (p *Printer) Start(name string, sep string) error {
	if err := errorIfAny("field name", name, invalidNameChars); err != nil {
		return err
	}
	return p.start([]byte(name), []byte(sep))
}

func (p *Printer) start(name []byte, sep []byte) error {
	needed := p.Indent - len(name) // FIXME: Count multi-byte characters correctly
	wg := writegroup{Writer: p.Writer}

	// Write the name and any indentation
	switch p.Alignment {
	case AlignmentRight:
		if needed > 0 {
			wg.Write(makeIndent(needed))
		}
		wg.Write(name)
	case AlignmentLeft:
		wg.Write(name)
		if needed > 0 {
			wg.Write(makeIndent(needed))
		}
	default:
		wg.Write(name)
	}

	// Write the separator
	if len(sep) > 0 {
		wg.Write(sep)
	}

	p.Started = true

	return wg.Err()
}

// Section writes a header starting a new section.
func (p *Printer) Section(section string, templates ...string) error {
	if err := errorIfAny("section name", section, invalidSectionChars); err != nil {
		return err
	}
	for _, template := range templates {
		if err := errorIfAny("section template", template, invalidTemplateChars); err != nil {
			return err
		}
	}
	wg := writegroup{Writer: p.Writer}

	if p.Started {
		wg.Write(newline)
	}
	wg.Write(sectionStart)
	wg.Write([]byte(section))
	wg.Write(sectionEnd)

	if len(templates) > 0 {
		wg.Write(templateStart)
		for i, template := range templates {
			if i > 0 {
				wg.Write(comma)
			}
			wg.Write([]byte(template))
		}
		wg.Write(templateEnd)
	}

	wg.Write(newline)

	return wg.Err()
}

/*
func (p *Printer) section(section []byte) error {
	return nil
}
*/

// Setting will write a setting to p.Writer.
//
// An InvalidContentError will be returned if the setting or its value
// contain an invalid character.
//
// If the underlying write operation fails an error will be returned.
func (p *Printer) Setting(setting, value string) error {
	if err := errorIfAny("setting name", setting, invalidNameChars); err != nil {
		return err
	}
	if err := errorIfAny("value", value, invalidValueChars); err != nil {
		return err
	}
	if err := p.start([]byte(setting), settingSeparator); err != nil {
		return err
	}
	wg := writegroup{Writer: p.Writer}
	wg.Write(escapeValue([]byte(value)))
	wg.Write(newline)
	return wg.Err()
}

// Object will print an object to p.Writer.
func (p *Printer) Object(object, value string) error {
	if err := errorIfAny("object name", object, invalidNameChars); err != nil {
		return err
	}
	if err := errorIfAny("value", value, invalidValueChars); err != nil {
		return err
	}
	if err := p.start([]byte(object), objectSeparator); err != nil {
		return err
	}
	wg := writegroup{Writer: p.Writer}
	wg.Write(escapeValue([]byte(value)))
	wg.Write(newline)
	return wg.Err()
}

// Write will write v to the underlying writer.
func (p *Printer) Write(v []byte) (n int, err error) {
	n, err = p.Writer.Write(v)
	if n > 0 {
		p.Started = true
	}
	return n, err
}

// Finish tells the printer to end a section.
//func (p *Printer) Finish() error {
//	_, err := w.w.Write([]byte("\n"))
//	return err
//}

// escapeValue escapes characters in a setting or object value that would
// otherwise change its meaning. An unescaped semicolon starts a comment.
//
// The returned slice is p itself when nothing needs escaping.
func escapeValue(p []byte) []byte {
	if bytes.IndexByte(p, ';') < 0 {
		return p
	}
	return bytes.ReplaceAll(p, []byte(";"), []byte("\\;"))
}

func errorIfAny(component, value, badChars string) error {
	if pos := strings.IndexAny(value, badChars); pos >= 0 {
		return InvalidContentError{
			Component: component,
			Value:     value,
			Pos:       pos,
		}
	}
	return nil
}

func errorIfAnyBytes(component string, value []byte, badChars string) error {
	if pos := bytes.IndexAny(value, badChars); pos >= 0 {
		return InvalidContentError{
			Component: component,
			Value:     string(value),
			Pos:       pos,
		}
	}
	return nil
}

// InvalidContentError is returned when configuration data contains invalid characters.
type InvalidContentError struct {
	Component string // The kind of component, such as "value" or "section name"
	Value     string // The offending content
	Pos       int    // The position of the first invalid character in Value
}

func (err InvalidContentError) Error() string {
	var kind = func(b byte) string {
		switch b {
		case '\n':
			return "a newline"
		case '\r':
			return "a carriage return"
		default:
			return fmt.Sprintf("an invalid character %q", b)
		}
	}
	return fmt.Sprintf("asterisk configuration %s %q contains %s", err.Component, err.Value, kind(err.Value[err.Pos]))
}
