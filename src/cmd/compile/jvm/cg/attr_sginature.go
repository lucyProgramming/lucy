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
	ret.NameIndex = class.insertUtf8Const("Signature")
	ret.Info = make([]byte, 2)
	index := class.insertUtf8Const(a.Signature)
	binary.BigEndian.PutUint16(ret.Info, index)
	return ret
}
