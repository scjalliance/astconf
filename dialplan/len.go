package dialplan

// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Function_LEN

// Len returns a dialplan function that returns the length of a string.
func Len(value Expression) LenFunc {
	return LenFunc{Value: value}
}

// LenFunc is a dialplan function that returns the length of a string.
type LenFunc struct {
	Value Expression
}

// Func returns the assembled function call.
func (f LenFunc) Func() FuncCall {
	return FuncCall{
		Name: "LEN",
		Args: []string{f.Value.Expr().String()},
	}
}

// Expr returns the presence state function as an expression.
func (f LenFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}
