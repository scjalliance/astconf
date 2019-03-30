package dialplan

import "strconv"

// Int is a dialplan number value that is restricted to integers.
type Int int

// Expr fulfills the expression interface.
func (v Int) Expr() ExprDef {
	return ExprDef{Content: strconv.Itoa(int(v)), Kind: NumLit}
}
