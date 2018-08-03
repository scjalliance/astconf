package dpmaqueue

import "github.com/scjalliance/astconf/astval"

// Config holds configuration for the DPMA queue application.
type Config struct {
	Queue           string       `astconf:"queue"`
	MemberName      string       `astconf:"membername"`
	Location        string       `astconf:"location"`
	Member          astval.YesNo `astconf:"member"`
	Permission      string       `astconf:"permission"`
	LoginExtension  string       `astconf:"login_exten,omitempty"`
	LogoutExtension string       `astconf:"logout_exten,omitempty"`
}
