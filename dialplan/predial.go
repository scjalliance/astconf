package dialplan

import (
	"fmt"
	"strconv"
	"strings"
)

// PredialTarget is a subroutine that an application runs on each callee
// channel before dialing it, such as the target of Page's b() option.
//
// It is written as context^extension^priority(arg1^arg2), the form Asterisk
// uses inside an option, where commas would end the option early.
type PredialTarget struct {
	Context   string
	Extension string
	Priority  int // Should be >= 1
	Args      []string
}

// String returns the target in option form.
func (t PredialTarget) String() string {
	var b strings.Builder
	b.WriteString(t.Context)
	b.WriteString("^")
	b.WriteString(t.Extension)
	b.WriteString("^")
	b.WriteString(strconv.Itoa(t.Priority))
	if len(t.Args) > 0 {
		b.WriteString("(")
		b.WriteString(strings.Join(t.Args, "^"))
		b.WriteString(")")
	}
	return b.String()
}

// Validate returns an error if the priority is below one, or if any part
// contains a character that would change the meaning of the option.
func (t PredialTarget) Validate() error {
	if t.Priority < 1 {
		return fmt.Errorf("predial priority must be at least 1, got %d", t.Priority)
	}
	if err := errorIfAny("predial context", t.Context, invalidArgChars+"^"); err != nil {
		return err
	}
	if err := errorIfAny("predial extension", t.Extension, invalidArgChars+"^"); err != nil {
		return err
	}
	for _, arg := range t.Args {
		if err := errorIfAny("predial argument", arg, invalidArgChars+"^"); err != nil {
			return err
		}
	}
	return nil
}
