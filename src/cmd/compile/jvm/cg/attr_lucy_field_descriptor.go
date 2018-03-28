package cg

import "encoding/binary"

type AttributeLucyFieldDescritor struct {
	Descriptor string
}

func (a *AttributeLucyFieldDescritor) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtfConst(ATTRIBUTE_NAME_LUCY_FIELD_DESCRIPTOR)
	ret.Info = make([]byte, 2)
	ret.attributeLength = 2
	binary.BigEndian.PutUint16(ret.Info, class.insertUtfConst(a.Descriptor))
	return ret
}
