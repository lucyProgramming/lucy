package cg

import "encoding/binary"

type AttributeLucyFieldDescriptor struct {
	Descriptor       string
	MethodAccessFlag uint16
}

func (a *AttributeLucyFieldDescriptor) FromBs(class *Class, bs []byte) {
	if len(bs) != 4 {
		panic("length is not 4")
	}
	a.Descriptor = string(class.ConstPool[binary.BigEndian.Uint16(bs[0:2])].Info)
	a.MethodAccessFlag = binary.BigEndian.Uint16(bs[2:4])
}

func (a *AttributeLucyFieldDescriptor) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyFieldDescriptor)
	ret.Info = make([]byte, 4)
	ret.attributeLength = 4
	binary.BigEndian.PutUint16(ret.Info, class.InsertUtf8Const(a.Descriptor))
	binary.BigEndian.PutUint16(ret.Info[2:], a.MethodAccessFlag)
	return ret
}
