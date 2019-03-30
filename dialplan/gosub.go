package dialplan

import "strconv"

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
