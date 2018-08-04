package cg

import (
	"encoding/binary"
)

type AttributeTemplateFunction struct {
	Name string
	/*
		reGenerate pos
	*/
	Filename    string
	StartLine   uint16
	StartColumn uint16
	Code        string
	AccessFlag  uint16
}

func (a *AttributeTemplateFunction) FromBytes(class *Class, bs []byte) {
	a.Name = string(class.ConstPool[binary.BigEndian.Uint16(bs)].Info)
	a.Filename = string(class.ConstPool[binary.BigEndian.Uint16(bs[2:])].Info)
	a.StartLine = binary.BigEndian.Uint16(bs[4:])
	a.StartColumn = binary.BigEndian.Uint16(bs[6:])
	a.Code = string(class.ConstPool[binary.BigEndian.Uint16(bs[8:])].Info)
	a.AccessFlag = binary.BigEndian.Uint16(bs[10:])
}

func (a *AttributeTemplateFunction) ToAttributeInfo(class *Class) *AttributeInfo {
	info := &AttributeInfo{}
	info.NameIndex = class.InsertUtf8Const(AttributeNameLucyTemplateFunction)
	info.attributeLength = 12
	info.Info = make([]byte, info.attributeLength)
	binary.BigEndian.PutUint16(info.Info, class.InsertUtf8Const(a.Name))
	binary.BigEndian.PutUint16(info.Info[2:], class.InsertUtf8Const(a.Filename))
	binary.BigEndian.PutUint16(info.Info[4:], a.StartLine)
	binary.BigEndian.PutUint16(info.Info[6:], a.StartColumn)
	binary.BigEndian.PutUint16(info.Info[8:], class.InsertUtf8Const(a.Code))
	binary.BigEndian.PutUint16(info.Info[10:], a.AccessFlag)
	return info
}
