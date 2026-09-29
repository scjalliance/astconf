package dialplan

// String is a dialplan string value.
type String string

// Expr fulfills the expression interface.
func (v String) Expr() ExprDef {
	return ExprDef{Content: string(v), Kind: StringLit}
}

// Validate returns an error if the string contains characters that would
// change the meaning of the expression or argument it appears in.
func (v String) Validate() error {
	return errorIfAny("string", string(v), invalidArgChars)
}
