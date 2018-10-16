package digium

import "encoding/xml"

// Contact is a Digium contact entry that can be serialized as XML.
type Contact struct {
	XMLName         xml.Name        `xml:"contact"`
	ServerUUID      string          `xml:"server_uuid,attr,omitempty"`
	ID              string          `xml:"id,attr,omitempty"`
	Prefix          string          `xml:"prefix,attr,omitempty"`
	FirstName       string          `xml:"first_name,attr"`
	SecondName      string          `xml:"second_name,attr,omitempty"`
	LastName        string          `xml:"last_name,attr,omitempty"`
	Suffix          string          `xml:"suffix,attr,omitempty"`
	Type            string          `xml:"contact_type,attr"`
	Organization    string          `xml:"organization,attr,omitempty"`
	JobTitle        string          `xml:"job_title,attr,omitempty"`
	Location        string          `xml:"location,attr,omitempty"`
	Notes           string          `xml:"notes,attr,omitempty"`
	AccountID       string          `xml:"account_id,attr,omitempty"`
	SubscriptionURI string          `xml:"subscribe_to,attr,omitempty"`
	PictureURL      string          `xml:"picture,attr,omitempty"`
	PickupAction    string          `xml:"pickup_action,attr,omitempty"`
	Emails          []ContactEmail  `xml:"emails"`
	Actions         []ContactAction `xml:"actions>action"`
}
