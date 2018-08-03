package dpmacontact

import "encoding/xml"

// Group is a DPMA contact group.
type Group struct {
	XMLName  xml.Name `xml:"contacts"`
	Name     string   `xml:"group_name,attr"`
	ID       string   `xml:"id,attr"`
	Contacts []Entry
}
