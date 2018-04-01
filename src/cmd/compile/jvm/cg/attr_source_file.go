package cg

import (
	"encoding/binary"
)

type AttributeSourceFile struct {
	file string
}

func (a *AttributeSourceFile) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtfConst(ATTRIBUTE_NAME_SOURCE_FILE)
	ret.attributeLength = 2
	ret.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.Info, class.insertUtfConst(a.file))
	return ret
}
