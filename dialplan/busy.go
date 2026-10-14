package dialplan

// Busy returns a dialplan application that signals busy to the caller.
func Busy() BusyApp {
	return BusyApp{}
}

// BusyApp is a dialplan application that signals busy to the caller.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Applications/Busy/
type BusyApp struct{}

// App returns the assembled application call.
func (app BusyApp) App() AppCall {
	return AppCall{Name: "Busy"}
}
