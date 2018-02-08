package cg

type AttributeInfo struct {
	nameIndex       [2]byte
	attributeLength uint32
	info            []byte
}

type ToAttributeInfo interface {
	ToAttributeInfo() *AttributeInfo
}
