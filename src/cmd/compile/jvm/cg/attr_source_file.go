package cg

import (
	"encoding/binary"
)

type AttributeSourceFile struct {
	file string
}

func (this *AttributeSourceFile) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameSourceFile)
	ret.attributeLength = 2
	ret.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(this.file))
	return ret
}
