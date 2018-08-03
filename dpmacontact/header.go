package dpmacontact

import "encoding/xml"

// Header is a header for a DPMA contact action.
type Header struct {
	XMLName xml.Name `xml:"header"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}
