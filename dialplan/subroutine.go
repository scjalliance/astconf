package dialplan

import (
	"fmt"
	"strconv"
	"strings"
)

// Subroutine returns a dialplan application that calls a subroutine with
// Gosub, passing args in parentheses after the priority. The subroutine
// reads them as ${ARG1}, ${ARG2} and so on, and must end in Return.
//
// Unlike GosubApp, which writes each argument after the priority as a
// separate application argument, Subroutine writes them the way Asterisk
// passes them, and each argument is an expression, so it may be a variable.
func Subroutine(context, extension string, priority int, args ...Expression) SubroutineApp {
	return SubroutineApp{
		Context:   context,
		Extension: extension,
		Priority:  priority,
		Args:      args,
	}
}

// SubroutineApp is a dialplan application that calls a subroutine with
// arguments.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Gosub/
type SubroutineApp struct {
	Context   string
	Extension string
	Priority  int // Should be >= 1
	Args      []Expression
}

// App returns the assembled application call.
func (sub SubroutineApp) App() AppCall {
	target := strconv.Itoa(sub.Priority)
	if len(sub.Args) > 0 {
		args := make([]string, 0, len(sub.Args))
		for _, arg := range sub.Args {
			args = append(args, arg.Expr().String())
		}
		target += "(" + strings.Join(args, ",") + ")"
	}
	return AppCall{
		Name: "Gosub",
		Args: []string{sub.Context, sub.Extension, target},
	}
}

// Validate returns an error if the priority is below one, or if the context,
// extension or any argument is invalid.
func (sub SubroutineApp) Validate() error {
	if sub.Priority < 1 {
		return fmt.Errorf("gosub priority must be at least 1, got %d", sub.Priority)
	}
	if err := errorIfAny("gosub context", sub.Context, invalidArgChars); err != nil {
		return err
	}
	if err := errorIfAny("gosub extension", sub.Extension, invalidArgChars); err != nil {
		return err
	}
	for _, arg := range sub.Args {
		if err := validate(arg); err != nil {
			return err
		}
	}
	return nil
}
