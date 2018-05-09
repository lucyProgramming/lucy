package cg

import (
	"encoding/binary"
)

type AttributeTemplateFunction struct {
	Name        string
	Filename    string
	StartLine   uint16
	StartColumn uint16
	Code        string
}

func (a *AttributeTemplateFunction) ToAttributeInfo(class *Class) *AttributeInfo {
	info := &AttributeInfo{}
	info.NameIndex = class.insertUtf8Const(ATTRIBUTE_NAME_LUCY_TEMPLATE_FUNCTION)
	info.attributeLength = 10
	info.Info = make([]byte, info.attributeLength)
	binary.BigEndian.PutUint16(info.Info, class.insertUtf8Const(a.Name))
	binary.BigEndian.PutUint16(info.Info[2:], class.insertUtf8Const(a.Filename))
	binary.BigEndian.PutUint16(info.Info[4:], a.StartLine)
	binary.BigEndian.PutUint16(info.Info[6:], a.StartColumn)
	binary.BigEndian.PutUint16(info.Info[8:], class.insertUtf8Const(a.Code))
	return info
}
