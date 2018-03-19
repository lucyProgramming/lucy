package cg

import "encoding/binary"

type AttributeLucyArrayDescriptor struct {
	Descriptor string
}

func (a *AttributeLucyArrayDescriptor) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtfConst(ATTRIBUTE_NAME_LUCY_ARRAY_DESCRIPTOR)
	ret.Info = make([]byte, 2)
	ret.attributeLength = 2
	binary.BigEndian.PutUint16(ret.Info, class.insertUtfConst(a.Descriptor))
	return ret
}
