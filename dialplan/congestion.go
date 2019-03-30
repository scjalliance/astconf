package dialplan

import "strconv"

// Congestion returns a dialplan application that signals congestion to the
// caller.
func Congestion() CongestionApp {
	return CongestionApp{}
}

// CongestionApp is a dialplan application that signals congestion to the
// caller.
//
// https://wiki.asterisk.org/wiki/display/AST/Asterisk+16+Application_Congestion
type CongestionApp struct {
	Timeout int // Number of seconds to play congestion before hanging up
}

// App returns the assembled application call.
func (app CongestionApp) App() AppCall {
	if app.Timeout > 0 {
		return AppCall{Name: "Congestion", Args: []string{strconv.Itoa(app.Timeout)}}
	}
	return AppCall{Name: "Congestion"}
}
