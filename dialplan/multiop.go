package dialplan

import (
	"fmt"
	"strings"
)

// MultiOp is a operation on two or more expressions.
type MultiOp struct {
	Expressions []Expression
	Operator    string
}

// Expr returns the binary operation as an expression.
func (op MultiOp) Expr() ExprDef {
	switch len(op.Expressions) {
	case 0, 1:
		panic("invalid number of expressions in logical operator")
	default:
		var builder strings.Builder
		for i := range op.Expressions {
			if i > 0 {
				builder.WriteString(op.Operator)
			}
			builder.WriteString(op.Expressions[i].Expr().String())
		}
		return ExprDef{Content: builder.String(), Kind: Op}
	}
}

// Or returns a logical or operation for two or more expressions.
func Or(e1 Expression, e2 Expression, extra ...Expression) MultiOp {
	return MultiOp{
		Expressions: append([]Expression{e1, e2}, extra...),
		Operator:    "|",
	}
}

// And returns a logical and operation for two or more expressions.
func And(e1 Expression, e2 Expression, extra ...Expression) MultiOp {
	return MultiOp{
		Expressions: append([]Expression{e1, e2}, extra...),
		Operator:    "&",
	}
}

// Validate returns an error if there are fewer than two operands, or if the
// operator or any operand is invalid.
func (op MultiOp) Validate() error {
	if len(op.Expressions) < 2 {
		return fmt.Errorf("a logical operation needs at least 2 expressions, got %d", len(op.Expressions))
	}
	if op.Operator == "" {
		return fmt.Errorf("a logical operation needs an operator")
	}
	if err := errorIfAny("logical operator", op.Operator, invalidOperatorChars); err != nil {
		return err
	}
	for _, expr := range op.Expressions {
		if err := validate(expr); err != nil {
			return err
		}
	}
	return nil
}
