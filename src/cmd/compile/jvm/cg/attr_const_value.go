package cg

import (
	"encoding/binary"
)

type AttributeConstantValue struct {
	Index uint16
}

func (a *AttributeConstantValue) ToAttributeInfo(class *Class) *AttributeInfo {
	info := &AttributeInfo{}
	binary.BigEndian.PutUint16(info.nameIndex[0:2], class.insertUtfConst("ConstantValue"))
	info.attributeLength = 2
	info.info = make([]byte, 2)
	binary.BigEndian.PutUint16(info.info, a.Index)
	return info
}
