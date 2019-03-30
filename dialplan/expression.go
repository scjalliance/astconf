package dialplan

// Expression is an interface that can be evaluated by the dialplan
// as an expression.
type Expression interface {
	Expr() ExprDef
}
