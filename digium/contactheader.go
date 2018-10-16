package digium

import "encoding/xml"

// ContactHeader is a header for a Digium contact action.
type ContactHeader struct {
	XMLName xml.Name `xml:"header"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}
