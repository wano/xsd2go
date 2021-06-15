package xsd


import (
"encoding/xml"
"strings"

"github.com/iancoleman/strcase"
)

// Attribute defines single XML attribute
type Pattern struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema pattern"`
	Value   string   `xml:"value,attr"`
}

// Public Go Name of this struct item
func (e *Pattern) GoName() string {
	return strcase.ToCamel(strings.ToLower(e.Value))
}

func (e *Pattern) Modifiers() string {
	return "-"
}

func (e *Pattern) XmlName() string {
	return e.Value
}

func (e *Pattern) compile(s *Schema) {
}
