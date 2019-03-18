package il

import "encoding/xml"

type TypeData int

const (
	_ = iota
	BoolType TypeData = iota * -1
	IntType
	CharType
	StringType
)

//type FieldData struct {
//	XMLName xml.Name `xml:"Field"`
//	Type int `xml: "Type"`
//}

type ContractData struct {
	XMLName xml.Name `xml:"Contract"`
	Identifier string `xml:"id, attr"`
	Fields []TypeData `xml:"Field"`
}
