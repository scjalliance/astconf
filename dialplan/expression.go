package dialplan

// Expression is an interface that can be evaluated by the dialplan
// as an expression.
type Expression interface {
	Expr() ExprDef
}

// Quoteable is an interface implemented by expressions that can request
// that their content should be wrapped in double quotes.
type Quoteable interface {
	QuotedContent() bool
}
