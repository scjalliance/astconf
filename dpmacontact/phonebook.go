package dpmacontact

import "encoding/xml"

// Phonebook is a DPMA phonebook containing contact groups.
type Phonebook struct {
	XMLName xml.Name `xml:"phonebooks"`
	Groups  []Group
}
