package cg

import "encoding/binary"

type AttributeLucyTypeAlias struct {
	Alias   string
	Comment string
}

func (this *AttributeLucyTypeAlias) FromBs(class *Class, bs []byte) {
	this.Alias = string(class.ConstPool[binary.BigEndian.Uint16(bs)].Info)
	this.Comment = string(class.ConstPool[binary.BigEndian.Uint16(bs[2:])].Info)
}

func (this *AttributeLucyTypeAlias) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyTypeAlias)
	ret.Info = make([]byte, 4)
	ret.attributeLength = uint32(len(ret.Info))
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(this.Alias))
	binary.BigEndian.PutUint16(ret.Info[2:], class.InsertUtf8Const(this.Comment))
	return ret
}
