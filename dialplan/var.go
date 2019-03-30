package dialplan

/*
// Var returns a variable reference.
func Var(v string) string {
	return "${" + v + "}"
}
*/

// Var is a variable reference.
type Var string

// String returns a string representation of the variable.
func (v Var) String() string {
	//return "${" + string(v) + "}"
	return string(v)
}

// Expr returns the variable as an expression.
func (v Var) Expr() ExprDef {
	return ExprDef{Content: v.String()}
}

// Ref returns the variable as a reference.
func (v Var) Ref() string {
	return string(v)
}
