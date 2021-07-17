package dialplan

// PresenceState returns a dialplan function that retrieves presence state.
func PresenceState(provider, field string) PresenceStateFunc {
	return PresenceStateFunc{
		Provider: provider,
		Field:    field,
	}
}

// PresenceStateFunc is a dialplan function that retrieves presence state.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Function_PRESENCE_STATE
type PresenceStateFunc struct {
	Provider string
	Field    string
}

// Func returns the assembled function call.
func (f PresenceStateFunc) Func() FuncCall {
	return FuncCall{
		Name: "PRESENCE_STATE",
		Args: []string{f.Provider, f.Field},
	}
}

// Expr returns the presence state function as an expression.
func (f PresenceStateFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}

// Ref returns a reference to the presence state function.
func (f PresenceStateFunc) Ref() string {
	return f.Func().String()
}

// QuotedContent returns true, which indicates that the function call shoud be
// wrapped in double quotes when compared.
func (f PresenceStateFunc) QuotedContent() bool {
	return true
}
