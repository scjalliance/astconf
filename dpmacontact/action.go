package dpmacontact

import "encoding/xml"

// Action is an action that can be taken for a DPMA contact.
type Action struct {
	XMLName    xml.Name `xml:"action"`
	ID         string   `xml:"id,attr"`
	Dial       string   `xml:"dial,attr,omitempty"`
	DialPrefix string   `xml:"dial_prefix,attr,omitempty"`
	Label      string   `xml:"label,attr"`
	Name       string   `xml:"name,attr"`
	Headers    []Header `xml:"headers"`
}
