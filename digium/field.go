package digium

import (
	"encoding/xml"

	"github.com/scjalliance/astconf/astval"
)

// Field is a Digium busy lamp field entry that can be serialized as XML.
type Field struct {
	XMLName    xml.Name           `xml:"blf_item"`
	Location   string             `xml:"location,attr"`
	Index      int                `xml:"index,attr"`
	Paging     astval.OneZeroNone `xml:"paging,attr,omitempty"`
	ContactID  string             `xml:"contact_id,attr,omitempty"`
	AppID      string             `xml:"app_id,attr,omitempty"`
	Blank      bool               `xml:"blank,attr,omitempty"`
	Behaviors  []FieldBehavior    `xml:"behaviors,omitempty"`
	Indicators []FieldIndicator   `xml:"indicators,omitempty"`
}
