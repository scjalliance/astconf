package digium

import "encoding/xml"

// Config is a Digium configuration that can be serialized as XML.
type Config struct {
	XMLName  xml.Name `xml:"config"`
	SmartBLF SmartBLF `xml:"smart_blf,omitempty"`
}
