package dialplan

import "strconv"

// DialExpr returns a dial application whose dial string is computed when
// the call runs, such as a variable that earlier actions built up.
func DialExpr(target Expression, seconds int) DialExprApp {
	return DialExprApp{
		Target:  target,
		Timeout: seconds,
	}
}

// DialExprApp is a dialplan application that dials a computed dial string.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Dial/
type DialExprApp struct {
	Target  Expression
	Timeout int // Seconds
}

// App returns the assembled application call.
func (dial DialExprApp) App() AppCall {
	app := AppCall{
		Name: "Dial",
		Args: []string{dial.Target.Expr().String()},
	}
	if dial.Timeout > 0 {
		app.Args = append(app.Args, strconv.Itoa(dial.Timeout))
	}
	return app
}

// Validate returns an error if the target is invalid.
func (dial DialExprApp) Validate() error {
	return validate(dial.Target)
}
