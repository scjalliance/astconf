package dialplan

import (
	"fmt"
	"strconv"
)

// https://wiki.asterisk.org/wiki/display/AST/Gosub

// Gosub returns a gosub application.
//
// If context or extension are empty, they will be omitted.
func Gosub(context, extension string, priority int, args ...string) GosubApp {
	return GosubApp{
		Context:   context,
		Extension: extension,
		Priority:  priority,
		Args:      args,
	}
}

// GosubApp is a dialplan application that invokes subroutines.
type GosubApp struct {
	Context   string
	Extension string
	Priority  int // Should be >= 1
	Args      []string
}

// App returns the assembled application call.
func (gosub GosubApp) App() AppCall {
	app := AppCall{Name: "Gosub"}
	if gosub.Context != "" {
		app.Args = append(app.Args, gosub.Context)
	}
	if gosub.Extension != "" {
		app.Args = append(app.Args, gosub.Extension)
	}
	app.Args = append(app.Args, strconv.Itoa(gosub.Priority))
	app.Args = append(app.Args, gosub.Args...)
	return app
}

// Validate returns an error if the priority is below one, or if the context,
// extension or any argument contains invalid characters.
func (gosub GosubApp) Validate() error {
	if gosub.Priority < 1 {
		return fmt.Errorf("gosub priority must be at least 1, got %d", gosub.Priority)
	}
	if err := errorIfAny("gosub context", gosub.Context, invalidArgChars); err != nil {
		return err
	}
	if err := errorIfAny("gosub extension", gosub.Extension, invalidArgChars); err != nil {
		return err
	}
	for _, arg := range gosub.Args {
		if err := errorIfAny("gosub argument", arg, invalidArgChars); err != nil {
			return err
		}
	}
	return nil
}
