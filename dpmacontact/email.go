package dpmacontact

import "encoding/xml"

// Email is an email address within a DPMA contact entry.
type Email struct {
	XMLName xml.Name `xml:"email"`
	Address string   `xml:"address,attr"`
	Label   string   `xml:"label,attr"`
	Primary bool     `xml:"primary,attr"`
}
