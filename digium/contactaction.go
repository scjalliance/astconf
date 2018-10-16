package digium

import "encoding/xml"

// ContactAction is an action that can be taken for a Digium contact.
type ContactAction struct {
	XMLName    xml.Name        `xml:"action"`
	ID         string          `xml:"id,attr"`
	Dial       string          `xml:"dial,attr,omitempty"`
	DialPrefix string          `xml:"dial_prefix,attr,omitempty"`
	Label      string          `xml:"label,attr"`
	Name       string          `xml:"name,attr"`
	Headers    []ContactHeader `xml:"headers"`
}
