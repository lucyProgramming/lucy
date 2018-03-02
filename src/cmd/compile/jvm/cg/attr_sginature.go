package cg

import (
	"encoding/binary"
)

type AttributeSignature struct {
	Signature string
}

func (a *AttributeSignature) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.attributeLength = 2
	ret.info = make([]byte, 2)
	index := class.insertUtfConst(a.Signature)
	binary.BigEndian.PutUint16(ret.info, index)
	return ret
}
