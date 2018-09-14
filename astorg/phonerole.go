package astorg

// PhoneRole is a role that is assigned to a particular phone.
type PhoneRole struct {
	Extension    string
	Username     string
	DisplayName  string
	Location     string
	AreaCode     string
	Firmware     string
	CallerID     string
	Apps         []string
	PagingGroups []string
	Phones       []string // MACs of phones to be assigned this role
}
