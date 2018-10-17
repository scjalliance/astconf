package digium

import "encoding/xml"

// DisplayRule is a Digium display rule that can be serialized as XML.
type DisplayRule struct {
	XMLName      xml.Name `xml:"display_rule"`
	ID           string   `xml:"id,attr,omitempty"`
	ActionID     string   `xml:"action_id,attr,omitempty"`
	PhoneState   string   `xml:"phone_state,attr,omitempty"`
	TargetStatus string   `xml:"target_status,attr,omitempty"`
	Show         bool     `xml:"show,attr"`
}
