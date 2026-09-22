package dialplan

// https://wiki.asterisk.org/wiki/display/AST/Macros (this page is wrong)

// Macro returns a macro application.
//
// Macros are deprecated and should be avoided if possible. The preferred
// subroutine calling mechanism is gosub.
func Macro(name string, args ...Expression) MacroApp {
	return MacroApp{
		Name: name,
		Args: args,
	}
}

// MacroApp is a deprecated dialplan application that invokes subroutines.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Macro
type MacroApp struct {
	Name string
	Args []Expression
}

// App returnes the assembled application call.
func (macro MacroApp) App() AppCall {
	args := make([]string, 0, len(macro.Args)+1)
	args = append(args, macro.Name)
	for _, arg := range macro.Args {
		args = append(args, arg.Expr().String())
	}
	return AppCall{
		Name: "Macro",
		Args: args,
	}
}

// Validate returns an error if the macro name or any argument is invalid.
func (macro MacroApp) Validate() error {
	if err := errorIfAny("macro name", macro.Name, invalidArgChars); err != nil {
		return err
	}
	for _, arg := range macro.Args {
		if err := validate(arg); err != nil {
			return err
		}
	}
	return nil
}
