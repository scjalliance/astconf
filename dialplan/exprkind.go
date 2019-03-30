package dialplan

// ExprKind describes the kind of an expression.
type ExprKind int

// Types of expressions.
const (
	Ref ExprKind = iota
	NumLit
	StringLit
	Op
)
