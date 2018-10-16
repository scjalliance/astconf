package digium

import "encoding/xml"

// ContactEmail is an email address within a Digium contact entry.
type ContactEmail struct {
	XMLName xml.Name `xml:"email"`
	Address string   `xml:"address,attr"`
	Label   string   `xml:"label,attr"`
	Primary bool     `xml:"primary,attr"`
}
