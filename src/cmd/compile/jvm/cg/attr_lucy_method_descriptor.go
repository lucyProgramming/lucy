package cg

import "encoding/binary"

type AttributeLucyMethodDescriptor struct {
	Descriptor string
}

func (a *AttributeLucyMethodDescriptor) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR)
	ret.Info = make([]byte, 2)
	ret.attributeLength = 2
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(a.Descriptor))
	return ret
}
