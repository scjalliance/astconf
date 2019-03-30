package dialplan

// String is a dialplan string value.
type String string

// Expr fulfills the expression interface.
func (v String) Expr() ExprDef {
	return ExprDef{Content: string(v), Kind: StringLit}
}
