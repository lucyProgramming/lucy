package cg

type AttributeInfo struct {
	NameIndex       uint16
	attributeLength uint32
	Info            []byte
}

type ToAttributeInfo interface {
	ToAttributeInfo() *AttributeInfo
}
