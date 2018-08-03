package dpmacontact

// Entry is a DPMA contact entry.
type Entry struct {
	Type            string   `xml:"contact_type,attr"`
	PictureURL      string   `xml:"picture,attr"`
	FirstName       string   `xml:"first_name,attr"`
	PickupAction    string   `xml:"pickup_action,attr"`
	SecondName      string   `xml:"second_name,attr,omitempty"`
	JobTitle        string   `xml:"job_title,attr"`
	LastName        string   `xml:"last_name,attr"`
	Prefix          string   `xml:"prefix,attr"`
	Organization    string   `xml:"organization,attr"`
	ID              string   `xml:"id,attr"`
	Suffix          string   `xml:"suffix,attr"`
	Location        string   `xml:"location,attr"`
	ServerUUID      string   `xml:"server_uuid,attr"`
	SubscriptionURI string   `xml:"subscribe_to,attr"`
	AccountID       string   `xml:"account_id,attr"`
	Notes           string   `xml:"notes,attr"`
	Emails          []Email  `xml:"emails"`
	Actions         []Action `xml:"actions"`
}
