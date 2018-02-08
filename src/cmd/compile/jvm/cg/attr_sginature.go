package cg

import (
	"encoding/binary"
)

type AttributeSignature struct {
	index uint16
}

func (a *AttributeSignature) ToAttributeInfo() *AttributeInfo {
	ret := &AttributeInfo{}
	ret.attributeLength = 2
	ret.info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.info, uint16(a.index))
	return ret
}
