package dialplan

// If returns a dialplan function that conditionally evaluates expressions.
//
// When expr is true, retValue will be returned.
func If(expr Expression, retValue Expression) IfFunc {
	return IfFunc{
		Expression: expr,
		IfTrue:     retValue,
	}
}

// IfElse returns a dialplan function that conditionally evaluates
// expressions.
//
// When expr is true, trueValue will be returned.
//
// When expr is false, falseValue will be returned.
func IfElse(expr Expression, trueValue, falseValue Expression) IfFunc {
	return IfFunc{
		Expression: expr,
		IfTrue:     trueValue,
		IfFalse:    falseValue,
	}
}

// IfFunc is a dialplan function that conditionally evaluates expressions.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Function_IF
type IfFunc struct {
	Expression
	IfTrue  Expression
	IfFalse Expression
}

// Func returns the assembled function call.
func (function IfFunc) Func() FuncCall {
	call := FuncCall{
		Name: "IF",
		Args: []string{function.Expression.Expr().String()},
	}
	if function.IfTrue != nil {
		call.TrueArgs = []string{function.IfTrue.Expr().String()}
	}
	if function.IfFalse != nil {
		call.FalseArgs = []string{function.IfFalse.Expr().String()}
	}
	return call
}

// Expr returns the function as an expression.
func (function IfFunc) Expr() ExprDef {
	return ExprDef{Content: function.Func().String()}
}
