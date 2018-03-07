package cg

import (
	"encoding/binary"
)

type AttributeConstantValue struct {
	Index uint16
}

func (a *AttributeConstantValue) ToAttributeInfo(class *Class) *AttributeInfo {
	info := &AttributeInfo{}
	info.NameIndex = class.insertUtfConst("ConstantValue")
	info.attributeLength = 2
	info.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(info.Info, a.Index)
	return info
}
