package cg

import (
	"encoding/binary"
)

type AttributeConstantValue struct {
	Index uint16
}

func (a *AttributeConstantValue) ToAttributeInfo(class *Class) *AttributeInfo {
	info := &AttributeInfo{}
	info.NameIndex = class.InsertUtf8Const(ATTRIBUTE_NAME_CONST_VALUE)
	info.attributeLength = 2
	info.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(info.Info, a.Index)
	return info
}
