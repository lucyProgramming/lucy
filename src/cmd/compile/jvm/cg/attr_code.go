package cg

import (
	"encoding/binary"
)

type AttributeCode struct {
	MaxStack          uint16
	MaxLocals         uint16
	CodeLength        int
	Codes             []byte
	LineNumbers       AttributeLineNumber
	Exceptions        []*ExceptionTable
	attributes        []*AttributeInfo
	AttributeStackMap AttributeStackMap
}

type ExceptionTable struct {
	StartPc   uint16
	EndPc     uint16
	HandlerPc uint16
	CatchType uint16
}

func (a *AttributeCode) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtf8Const("Code")
	ret.Info = make([]byte, 8)
	binary.BigEndian.PutUint16(ret.Info[0:2], a.MaxStack)
	binary.BigEndian.PutUint16(ret.Info[2:4], a.MaxLocals)
	binary.BigEndian.PutUint32(ret.Info[4:8], uint32(a.CodeLength))
	ret.Info = append(ret.Info, a.Codes...)
	ret.Info = append(ret.Info, a.mkExceptions()...)
	if info := a.LineNumbers.ToAttributeInfo(class); info != nil {
		a.attributes = append(a.attributes, info)
	}
	if info := a.AttributeStackMap.ToAttributeInfo(class); info != nil {
		a.attributes = append(a.attributes, info)
	}
	ret.Info = append(ret.Info, a.mkAttributes(class)...)
	ret.attributeLength = uint32(len(ret.Info))
	return ret
}

/*
	mk line number attribute
*/
func (a *AttributeCode) MKLineNumber(lineno int) {
	line := &AttributeLinePc{}
	line.startPc = uint16(a.CodeLength)
	line.lineNumber = uint16(lineno)
	a.LineNumbers.lineNumbers = append(a.LineNumbers.lineNumbers, line)
}

func (a *AttributeCode) mkAttributes(class *Class) []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.attributes)))

	if len(a.attributes) > 0 {
		b := make([]byte, 0)
		for _, v := range a.attributes {
			bb := make([]byte, 2)
			binary.BigEndian.PutUint16(bb, v.NameIndex)
			bs4 := make([]byte, 4)
			binary.BigEndian.PutUint32(bs4, uint32(v.attributeLength))
			bb = append(bb, bs4...)
			bb = append(bb, v.Info...)
			b = append(b, bb...)
		}
		bs = append(bs, b...)
	}
	return bs
}

func (a *AttributeCode) mkExceptions() []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.Exceptions)))
	if len(a.Exceptions) > 0 {
		b := make([]byte, 8*len(a.Exceptions))
		for k, v := range a.Exceptions {
			binary.BigEndian.PutUint16(b[k*8:], v.StartPc)
			binary.BigEndian.PutUint16(b[k*8+2:], v.EndPc)
			binary.BigEndian.PutUint16(b[k*8+4:], v.HandlerPc)
			binary.BigEndian.PutUint16(b[k*8+6:], v.CatchType)
		}
		bs = append(bs, b...)
	}
	return bs
}
