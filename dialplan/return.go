package dialplan

// Return returns a dialplan application that returns from a subroutine.
func Return() ReturnApp {
	return ReturnApp{}
}

// ReturnApp is a dialplan application that returns from a subroutine
// entered with Gosub or a predial option.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Return/
type ReturnApp struct{}

// App returns the assembled application call.
func (app ReturnApp) App() AppCall {
	return AppCall{Name: "Return"}
}
