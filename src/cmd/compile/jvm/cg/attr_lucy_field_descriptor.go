package cg

import "encoding/binary"

type AttributeLucyFieldDescriptor struct {
	Descriptor string
}

func (a *AttributeLucyFieldDescriptor) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyFieldDescriptor)
	ret.Info = make([]byte, 2)
	ret.attributeLength = 2
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(a.Descriptor))
	return ret
}
