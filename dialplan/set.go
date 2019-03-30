package dialplan

import "fmt"

// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Set

// Set returns a dialplan application that sets channel and function
// variables.
func Set(name Reference, value Expression) SetApp {
	return SetApp{
		Name:  name,
		Value: value,
	}
}

// SetApp is a dialplan application that sets channel and function
// variables.
type SetApp struct {
	Name  Reference
	Value Expression
}

// App returns the assembled application call.
func (app SetApp) App() AppCall {
	return AppCall{
		Name: "Set",
		Args: []string{fmt.Sprintf("%s=%s", app.Name.Ref(), app.Value.Expr().String())},
	}
}
