package cg

import "encoding/binary"

type AttributeLucyMethodDescritor struct {
	Descriptor string
}

func (a *AttributeLucyMethodDescritor) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtfConst(ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR)
	ret.Info = make([]byte, 2)
	ret.attributeLength = 2
	binary.BigEndian.PutUint16(ret.Info, class.insertUtfConst(a.Descriptor))
	return ret
}
