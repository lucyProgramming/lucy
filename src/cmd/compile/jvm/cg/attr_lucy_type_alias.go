package cg

import "encoding/binary"

type AttributeLucyTypeAlias struct {
	Alias   string
	Comment string
}

func (a *AttributeLucyTypeAlias) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyTypeAlias)
	ret.Info = make([]byte, 4)
	ret.attributeLength = uint32(len(ret.Info))
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(a.Alias))
	binary.BigEndian.PutUint16(ret.Info[2:], class.InsertUtf8Const(a.Comment))
	return ret
}
