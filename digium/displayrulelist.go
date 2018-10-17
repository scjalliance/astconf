package digium

import "encoding/xml"

// DisplayRuleList is an ordered list of Digium display rules that can be
// serialized as XML.
type DisplayRuleList struct {
	XMLName xml.Name `xml:"display_rules"`
	Rules   []DisplayRule
}
