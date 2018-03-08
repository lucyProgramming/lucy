package cg

import (
	"encoding/binary"
	"fmt"
)

type AttributeLineNumber struct {
	linenumbers []*AttributeLinePc
}

func (a *AttributeLineNumber) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil || len(a.linenumbers) == 0 {
		return nil
	}
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtfConst("LineNumberTable")
	fmt.Println("!!!!!!!!!!!!", ret.NameIndex)
	ret.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.Info, uint16(len(a.linenumbers)))
	for _, v := range a.linenumbers {
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
