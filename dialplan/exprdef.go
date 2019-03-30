package dialplan

// ExprDef is a dialplan expresssion definition.
type ExprDef struct {
	Content string
	Kind    ExprKind
}

// String returns the invocation string for the expression.
func (expr ExprDef) String() string {
	switch expr.Kind {
	default:
		return "${" + expr.Content + "}"
	case NumLit, StringLit:
		return expr.Content
	case Op:
		return "$[" + expr.Content + "]"
	}
}
