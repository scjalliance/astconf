package dpmaapp

// Queue holds configuraiton for the DPMA queue application.
type Queue struct {
	Queue           string `astconf:"queue"`
	MemberName      string `astconf:"membername"`
	Location        string `astconf:"location"`
	Member          bool   `astconf:"member"`
	Permission      string `astconf:"permission"`
	LoginExtension  bool   `astconf:"login_exten,omitempty"`
	LogoutExtension bool   `astconf:"logout_exten,omitempty"`
}
