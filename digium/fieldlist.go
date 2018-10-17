package digium

import (
	"encoding/xml"
)

type rawFieldList struct {
	XMLName xml.Name `xml:"blf_items"`
	Fields  []Field
}

// FieldList is an ordered list of Digium busy lamp field entries that
// can be serialized as XML.
type FieldList []Field

// MarshalXML marshals the field list as XML.
func (fields FieldList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "blf_items"
	return e.EncodeElement(&rawFieldList{Fields: fields}, start)
}

// UnmarshalXML marshals the field list as XML.
func (fields *FieldList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw rawFieldList
	err := d.DecodeElement(&raw, &start)
	*fields = raw.Fields
	return err
}
