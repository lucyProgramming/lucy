package cg

import (
	"encoding/binary"
)

type AttributeLineNumber struct {
	lineNumbers []*AttributeLinePc
}

func (this *AttributeLineNumber) ToAttributeInfo(class *Class) *AttributeInfo {
	if this == nil || len(this.lineNumbers) == 0 {
		return nil
	}
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const("LineNumberTable")
	ret.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.Info, uint16(len(this.lineNumbers)))
	for _, v := range this.lineNumbers {
		bs4 := make([]byte, 4)
		binary.BigEndian.PutUint16(bs4[0:2], v.startPc)
		binary.BigEndian.PutUint16(bs4[2:4], v.lineNumber)
		ret.Info = append(ret.Info, bs4...)
	}
	ret.attributeLength = uint32(len(ret.Info))
	return ret
}

type AttributeLinePc struct {
	startPc    uint16
	lineNumber uint16
}
