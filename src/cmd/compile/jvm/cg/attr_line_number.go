package cg

import (
	"encoding/binary"
)

type AttributeLineNumber struct {
	linenumbers []*AttributeLinePc
}

func (a *AttributeLineNumber) ToAttributeInfo(class *Class) *AttributeInfo {
	if len(a.linenumbers) == 0 {
		return nil
	}
	ret := &AttributeInfo{}
	binary.BigEndian.PutUint16(ret.nameIndex[0:2], class.insertUtfConst("LineNumberTable"))
	ret.info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.info, uint16(len(a.linenumbers)))
	for _, v := range a.linenumbers {
		bs4 := make([]byte, 4)
		binary.BigEndian.PutUint16(bs4[0:2], v.startPc)
		binary.BigEndian.PutUint16(bs4[2:4], v.lineNumber)
		ret.info = append(ret.info, bs4...)
	}
	ret.attributeLength = uint32(len(ret.info))
	return ret
}

type AttributeLinePc struct {
	startPc    uint16
	lineNumber uint16
}
