package dialplan

// ExecIf returns a dialplan application that conditionally executes other
// applications.
//
// When expr is true, app will be executed.
func ExecIf(expr Expression, app Application) ExecIfApp {
	return ExecIfApp{
		Expression: expr,
		IfTrue:     app,
	}
}

// ExecIfElse returns a dialplan application that conditionally executes other
// applications.
//
// When expr is true, trueApp will be executed.
//
// When expr is false, falseApp will be executed.
func ExecIfElse(expr Expression, trueApp, falseApp Application) ExecIfApp {
	return ExecIfApp{
		Expression: expr,
		IfTrue:     trueApp,
		IfFalse:    falseApp,
	}
}

// ExecIfApp is a dialplan application that conditionally executes other
// applications.
type ExecIfApp struct {
	Expression
	IfTrue  Application
	IfFalse Application
}

// App returns the assembled application call.
func (app ExecIfApp) App() AppCall {
	call := AppCall{
		Name: "ExecIf",
		Args: []string{app.Expr().String()},
	}
	if app.IfTrue != nil {
		call.TrueArgs = []string{app.IfTrue.App().String()}
	}
	if app.IfFalse != nil {
		call.FalseArgs = []string{app.IfFalse.App().String()}
	}
	return call
}
