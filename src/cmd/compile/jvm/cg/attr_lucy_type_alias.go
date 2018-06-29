package cg

import "encoding/binary"

type AttributeLucyTypeAlias struct {
	Alias string
}

func (a *AttributeLucyTypeAlias) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyTypeAlias)
	ret.Info = make([]byte, 2)
	ret.attributeLength = 2
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(a.Alias))
	return ret
}
