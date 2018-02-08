package cg

import (
	"encoding/binary"
)

type AttributeCode struct {
	MaxStack    uint16
	MaxLocals   uint16
	CodeLength  uint16
	Codes       []byte
	LineNumbers AttributeLineNumber
	Exceptions  []*ExceptionTable
	attributes  []*AttributeInfo
}

type ExceptionTable struct {
	StartPc   uint16
	Endpc     uint16
	HandlerPc uint16
	CatchType uint16
}

func (a *AttributeCode) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.info = make([]byte, 8)
	binary.BigEndian.PutUint16(ret.info[0:2], a.MaxStack)
	binary.BigEndian.PutUint16(ret.info[2:4], a.MaxLocals)
	binary.BigEndian.PutUint32(ret.info[4:8], uint32(a.CodeLength))
	ret.info = append(ret.info, a.Codes...)
	ret.info = append(ret.info, a.mkExceptions()...)
	a.attributes = append(a.attributes, a.LineNumbers.ToAttributeInfo(class))
	ret.info = append(ret.info, a.mkAttributes(class)...)
	ret.attributeLength = uint32(len(ret.info))
	return ret
}

/*
	mk line number attribute
*/
func (a *AttributeCode) MKLineNumber(lineno int) {
	line := &AttributeLinePc{}
	line.startPc = a.CodeLength
	line.lineNumber = uint16(lineno)
	a.LineNumbers.linenumbers = append(a.LineNumbers.linenumbers, line)
}

func (a *AttributeCode) mkAttributes(class *Class) []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.attributes)))
	if len(a.attributes) > 0 {
		b := make([]byte, 0)
		for _, v := range a.attributes {
			bb := []byte{}
			bb = append(bb, v.nameIndex[0:2]...)
			bs4 := make([]byte, 4)
			binary.BigEndian.PutUint32(bs4, uint32(v.attributeLength))
			bb = append(bb, bs4...)
			bb = append(bb, v.info...)
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
			binary.BigEndian.PutUint16(b[k*8+2:], v.Endpc)
			binary.BigEndian.PutUint16(b[k*8+4:], v.HandlerPc)
			binary.BigEndian.PutUint16(b[k*8+6:], v.CatchType)
		}
		bs = append(bs, b...)
	}
	return bs
}
