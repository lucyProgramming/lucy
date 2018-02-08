package cg

import (
	"encoding/binary"
)

type AttributeSourceFile struct {
	s string
}

func (a *AttributeSourceFile) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	binary.BigEndian.PutUint16(ret.nameIndex[0:2], class.insertUtfConst("SourceFile"))
	ret.attributeLength = 2
	ret.info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.info, class.insertUtfConst(a.s))
	return ret
}
