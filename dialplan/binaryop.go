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
	if op.quoted() {
		return ExprDef{Content: `"` + e1 + `"` + op.Operator + `"` + e2 + `"`, Kind: Op}
	}
	return ExprDef{Content: e1 + op.Operator + e2, Kind: Op}
}

// quoted reports whether Expr will wrap each operand in double quotes, which
// it does when either operand asks for it.
func (op BinaryOp) quoted() bool {
	if q1, ok := op.E1.(Quoteable); ok && q1.QuotedContent() {
		return true
	}
	q2, ok := op.E2.(Quoteable)
	return ok && q2.QuotedContent()
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

// Validate returns an error if the operator or either operand is invalid.
//
// When Expr will wrap the operands in double quotes, a quote in either
// rendered operand would close it early and change the comparison, so it is
// rejected here rather than in the operand's own validation, where a quote is
// harmless outside a quoted comparison.
func (op BinaryOp) Validate() error {
	if err := errorIfAny("binary operator", op.Operator, invalidOperatorChars); err != nil {
		return err
	}
	if err := validate(op.E1); err != nil {
		return err
	}
	if err := validate(op.E2); err != nil {
		return err
	}
	if op.quoted() {
		if err := errorIfAny("quoted operand", op.E1.Expr().String(), `"`); err != nil {
			return err
		}
		if err := errorIfAny("quoted operand", op.E2.Expr().String(), `"`); err != nil {
			return err
		}
	}
	return nil
}
