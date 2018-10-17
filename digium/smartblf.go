package digium

import "encoding/xml"

// SmartBLF is a Digium Smart BLF configuration that can be
// serialized as XML.
type SmartBLF struct {
	XMLName xml.Name  `xml:"smart_blf"`
	Fields  FieldList `xml:"blf_items,omitempty"`
}
