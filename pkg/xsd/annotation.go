package xsd

import (
	"encoding/xml"
)

// Attribute defines single XML attribute
type Annotation struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema annotation"`
	Documentation Documentation `xml:"documentation"`
}

func (self Annotation) Doc() string {
	return self.Documentation.Text
}

type Documentation struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema documentation"`
	Source string `xml:"source,attr"`
	Text string `xml:",chardata"`
}