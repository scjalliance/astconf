package dialplan

import (
	"fmt"
	"strings"
)

// FuncCall is a dialplan function call.
type FuncCall struct {
	Name      string
	Args      []string
	TrueArgs  []string // Args presented after a "?" character
	FalseArgs []string // Args presented after a ":" character
}

// String returns the invocation string for the function call.
//
// FIXME: Sanitize the function name and arguments.
func (call FuncCall) String() string {
	args := strings.Join(call.Args, ",")
	t := strings.Join(call.TrueArgs, ",")
	f := strings.Join(call.FalseArgs, ",")
	lt, lf := len(t), len(f)
	switch {
	case lt > 0 && lf > 0:
		return fmt.Sprintf("%s(%s?%s:%s)", call.Name, args, t, f)
	case lt > 0:
		return fmt.Sprintf("%s(%s?%s)", call.Name, args, t)
	case lf > 0:
		return fmt.Sprintf("%s(%s?:%s)", call.Name, args, f)
	default:
		return fmt.Sprintf("%s(%s)", call.Name, args)
	}
}
