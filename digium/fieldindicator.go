package digium

import "encoding/xml"

// FieldIndicator is an indicator entry for a Digium busy lamp field that can be
// serialized as XML.
type FieldIndicator struct {
	XMLName         xml.Name `xml:"indicator"`
	TargetStatus    string   `xml:"target_status,attr,omitempty"`
	Ring            bool     `xml:"ring,attr,omitempty"`
	RingtoneID      string   `xml:"ringtone_id,attr,omitempty"`
	LEDColor        string   `xml:"led_color,attr,omitempty"`
	LEDState        string   `xml:"led_state,attr,omitempty"`
	ForegroundColor string   `xml:"line_label_fgcolor,attr,omitempty"` // D6x phones only
	BackgroundColor string   `xml:"line_label_bgcolor,attr,omitempty"` // D6x phones only
}
