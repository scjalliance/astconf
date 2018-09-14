package astorg

// Person represents a person in an organization.
type Person struct {
	//Server             string // Identifies which server manages presence
	Extension          string // Primary Extension
	Username           string
	FirstName          string
	LastName           string
	FullName           string
	Organization       string
	Title              string
	Location           string
	Hidden             bool // Hidden from organization contact groups
	CallerID           string
	AreaCode           string
	Firmware           string
	CalendarURL        string
	VoicemailExtension string
	Ringtones          []string
	Apps               []string
	PagingGroups       []string
	ContactGroups      []string
	ContactNumbers     []Number // Additional phone numbers
	Phones             []string // MACs of phones to be assigned this role
	EmailAddresses     []Email
}
