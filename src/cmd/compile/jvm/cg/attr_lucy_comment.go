package cg

import "encoding/binary"

type AttributeLucyComment struct {
	Comment string
}

func (a *AttributeLucyComment) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.attributeLength = 2
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyComment)
	ret.Info = make([]byte, 2)
	index := class.InsertUtf8Const(a.Comment)
	binary.BigEndian.PutUint16(ret.Info, index)
	return ret
}
