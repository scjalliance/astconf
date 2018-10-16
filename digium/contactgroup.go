package digium

import "encoding/xml"

// ContactGroup is a Digium contact group.
type ContactGroup struct {
	XMLName  xml.Name `xml:"contacts"`
	Name     string   `xml:"group_name,attr"`
	ID       string   `xml:"id,attr"`
	Contacts []Contact
}
