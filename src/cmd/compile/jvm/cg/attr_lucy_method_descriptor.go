package cg

import "encoding/binary"

type AttributeLucyMethodDescriptor struct {
	Descriptor string
}

func (this *AttributeLucyMethodDescriptor) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyMethodDescriptor)
	ret.Info = make([]byte, 2)
	ret.attributeLength = 2
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(this.Descriptor))
	return ret
}
