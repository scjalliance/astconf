package dialplan

// PJSIPHeader returns a dialplan function that reads or changes a SIP
// header on a PJSIP channel. Action is "read", "add", "update" or
// "remove".
func PJSIPHeader(action, name string) PJSIPHeaderFunc {
	return PJSIPHeaderFunc{
		Action: action,
		Name:   name,
	}
}

// PJSIPHeaderFunc is a dialplan function that reads or changes a SIP header
// on a PJSIP channel.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Functions/PJSIP_HEADER/
type PJSIPHeaderFunc struct {
	Action string
	Name   string
}

// Func returns the assembled function call.
func (f PJSIPHeaderFunc) Func() FuncCall {
	return FuncCall{
		Name: "PJSIP_HEADER",
		Args: []string{f.Action, f.Name},
	}
}

// Expr returns the function as an expression.
func (f PJSIPHeaderFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}

// Ref returns a reference to the function, for use as the name in Set.
func (f PJSIPHeaderFunc) Ref() string {
	return f.Func().String()
}

// Validate returns an error if the action or header name contains invalid
// characters.
func (f PJSIPHeaderFunc) Validate() error {
	if err := errorIfAny("pjsip header action", f.Action, invalidArgChars); err != nil {
		return err
	}
	return errorIfAny("pjsip header name", f.Name, invalidArgChars)
}
