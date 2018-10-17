package digium

import "encoding/xml"

// FieldBehavior is a behavior entry for a Digium busy lamp field that can be
// serialized as XML.
type FieldBehavior struct {
	XMLName           xml.Name `xml:"behavior"`
	PhoneState        string   `xml:"phone_state,attr,omitempty"`
	TargetStatus      string   `xml:"target_status,attr,omitempty"`
	PressAction       string   `xml:"press_action,attr,omitempty"`
	PressFunction     string   `xml:"press_function,attr,omitempty"`
	LongPressAction   string   `xml:"long_press_action,attr,omitempty"`
	LongPressFunction string   `xml:"long_press_function,attr,omitempty"`
}
