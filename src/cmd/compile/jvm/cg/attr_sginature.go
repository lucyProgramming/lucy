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
	binary.BigEndian.PutUint16(ret.nameIndex[0:2], class.insertUtfConst("Signature"))
	ret.info = make([]byte, 2)
	index := class.insertUtfConst(a.Signature)
	binary.BigEndian.PutUint16(ret.info, index)
	return ret
}
