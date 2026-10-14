package dialplan

import (
	"strconv"
	"strings"
)

// Cat returns an expression that joins the rendered form of each part, such
// as a variable followed by literal text.
func Cat(parts ...Expression) CatExpr {
	return CatExpr{Parts: parts}
}

// CatExpr is an expression that joins the rendered form of its parts.
type CatExpr struct {
	Parts []Expression
}

// Expr returns the joined parts as a literal expression.
func (cat CatExpr) Expr() ExprDef {
	var builder strings.Builder
	for _, part := range cat.Parts {
		builder.WriteString(part.Expr().String())
	}
	return ExprDef{Content: builder.String(), Kind: StringLit}
}

// Validate returns an error if any part is invalid.
func (cat CatExpr) Validate() error {
	for _, part := range cat.Parts {
		if err := validate(part); err != nil {
			return err
		}
	}
	return nil
}

// Substring returns an expression that reads a variable from the given
// offset onward, as ${name:offset}.
func Substring(name Var, offset int) SubstringExpr {
	return SubstringExpr{Name: name, Offset: offset}
}

// SubstringExpr is an expression that reads part of a variable.
type SubstringExpr struct {
	Name   Var
	Offset int
}

// Expr returns the substring as a variable reference.
func (sub SubstringExpr) Expr() ExprDef {
	return ExprDef{Content: string(sub.Name) + ":" + strconv.Itoa(sub.Offset)}
}

// Validate returns an error if the variable name contains invalid
// characters.
func (sub SubstringExpr) Validate() error {
	return errorIfAny("variable name", string(sub.Name), invalidArgChars+":")
}
