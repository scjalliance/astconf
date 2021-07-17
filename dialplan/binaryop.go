package dialplan

// BinaryOp is a binary operation on two expressions.
type BinaryOp struct {
	E1       Expression
	E2       Expression
	Operator string
}

// Expr returns the binary operation as an expression.
func (op BinaryOp) Expr() ExprDef {
	e1 := op.E1.Expr().String()
	e2 := op.E2.Expr().String()
	quoted := false
	if q1, ok := op.E1.(Quoteable); ok {
		quoted = quoted || q1.QuotedContent()
	}
	if q2, ok := op.E2.(Quoteable); ok {
		quoted = quoted || q2.QuotedContent()
	}
	if quoted {
		return ExprDef{Content: `"` + e1 + `"` + op.Operator + `"` + e2 + `"`, Kind: Op}
	}
	return ExprDef{Content: e1 + op.Operator + e2, Kind: Op}
}

// Equal returns an equality operation.
func Equal(e1 Expression, e2 Expression) BinaryOp {
	return BinaryOp{E1: e1, E2: e2, Operator: "="}
}

// GreaterThan returns a greater-than operation.
func GreaterThan(e1 Expression, e2 Expression) BinaryOp {
	return BinaryOp{E1: e1, E2: e2, Operator: ">"}
}

// GreaterThanOrEqual returns a greater-than-or-equal operation.
func GreaterThanOrEqual(e1 Expression, e2 Expression) BinaryOp {
	return BinaryOp{E1: e1, E2: e2, Operator: ">="}
}

// LessThan returns a less-than operation.
func LessThan(e1 Expression, e2 Expression) BinaryOp {
	return BinaryOp{E1: e1, E2: e2, Operator: "<"}
}

// LessThanOrEqual returns a less-than-or-equal operation.
func LessThanOrEqual(e1 Expression, e2 Expression) BinaryOp {
	return BinaryOp{E1: e1, E2: e2, Operator: "<="}
}

// NotEqual returns an inequality operation.
func NotEqual(e1 Expression, e2 Expression) BinaryOp {
	return BinaryOp{E1: e1, E2: e2, Operator: "!="}
}
