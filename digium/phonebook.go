package digium

import "encoding/xml"

// Phonebook is a Digium phonebook containing contact groups.
type Phonebook struct {
	XMLName xml.Name `xml:"phonebooks"`
	Groups  []ContactGroup
}
