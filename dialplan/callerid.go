package dialplan

// CallerID returns a dialplan function that gets Caller ID data for the
// channel.
func CallerID(dataType string) CallerIDFunc {
	return CallerIDFunc{
		DataType: dataType,
	}
}

// CallerIDFunc is a dialplan function that gets Caller ID data for the
// channel.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Function_CALLERID
type CallerIDFunc struct {
	DataType string
}

// Func returns the assembled function call.
func (f CallerIDFunc) Func() FuncCall {
	return FuncCall{
		Name: "CALLERID",
		Args: []string{f.DataType},
	}
}

// Expr returns the caller ID function as an expression.
func (f CallerIDFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}

// Ref returns a reference to the caller ID function.
func (f CallerIDFunc) Ref() string {
	return f.Func().String()
}
