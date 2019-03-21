package astconf

import (
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

//var newlineReplacer = strings.NewReplacer("")

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

//func (p *Printer) Template(name string, templates ...string) error {
//}

// Comment writes a single line comment to p.Writer.
func (p *Printer) Comment(comment string) error {
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
func (p *Printer) Start(name string, sep string) error {
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
	//if err := errorIfAny("section", section, "[]\n"); err != nil {
	//	return err
	//}
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
	/*
		if err := errorIfAny("setting", setting, "\n="); err != nil {
			return err
		}
		if err := errorIfAny("value", value, "\n"); err != nil {
			return err
		}
	*/
	if err := p.start([]byte(setting), settingSeparator); err != nil {
		return err
	}
	wg := writegroup{Writer: p.Writer}
	wg.Write([]byte(value))
	wg.Write(newline)
	return wg.Err()
}

// Include will print an include construct to p.Writer for the given path.
func (p *Printer) Include(path string) error {
	return nil
}

// Object will print an object to p.Writer.
func (p *Printer) Object(object, value string) error {
	if err := p.start([]byte(object), objectSeparator); err != nil {
		return err
	}
	wg := writegroup{Writer: p.Writer}
	wg.Write([]byte(value))
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

func stripChars(str, chars string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chars, r) < 0 {
			return r
		}
		return -1
	}, str)
}

func errorIfAny(component, value, badChars string) error {
	if pos := strings.IndexAny(value, badChars); pos >= 0 {
		return InvalidContentError{
			Component: component,
			Value:     stripChars(value, badChars),
			Pos:       pos,
		}
	}
	return nil
}

// InvalidContentError is returned when configuration data contains invalid characters.
type InvalidContentError struct {
	Component string
	Value     string
	Pos       int
}

func (err InvalidContentError) Error() string {
	var kind = func(b byte) string {
		switch b {
		case '\n':
			return "newline character"
		default:
			return "an invalid character"
		}
	}
	return fmt.Sprintf("asterisk configuration %s \"%s\" contains %s", err.Component, err.Value, kind(err.Value[err.Pos]))
}
