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

func (this *AttributeCode) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const("Code")
	ret.Info = make([]byte, 8)
	binary.BigEndian.PutUint16(ret.Info[0:2], this.MaxStack)
	binary.BigEndian.PutUint16(ret.Info[2:4], this.MaxLocals)
	binary.BigEndian.PutUint32(ret.Info[4:8], uint32(this.CodeLength))
	ret.Info = append(ret.Info, this.Codes...)
	ret.Info = append(ret.Info, this.mkExceptions()...)
	if info := this.LineNumbers.ToAttributeInfo(class); info != nil {
		this.attributes = append(this.attributes, info)
	}
	if info := this.AttributeStackMap.ToAttributeInfo(class); info != nil {
		this.attributes = append(this.attributes, info)
	}
	ret.Info = append(ret.Info, this.mkAttributes(class)...)
	ret.attributeLength = uint32(len(ret.Info))
	return ret
}

/*
	mk line number attribute
*/
func (this *AttributeCode) MKLineNumber(lineNumber int) {
	line := &AttributeLinePc{}
	line.startPc = uint16(this.CodeLength)
	line.lineNumber = uint16(lineNumber)
	this.LineNumbers.lineNumbers = append(this.LineNumbers.lineNumbers, line)
}

func (this *AttributeCode) mkAttributes(class *Class) []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(this.attributes)))

	if len(this.attributes) > 0 {
		b := make([]byte, 0)
		for _, v := range this.attributes {
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

func (this *AttributeCode) mkExceptions() []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(this.Exceptions)))
	if len(this.Exceptions) > 0 {
		b := make([]byte, 8*len(this.Exceptions))
		for k, v := range this.Exceptions {
			binary.BigEndian.PutUint16(b[k*8:], v.StartPc)
			binary.BigEndian.PutUint16(b[k*8+2:], v.EndPc)
			binary.BigEndian.PutUint16(b[k*8+4:], v.HandlerPc)
			binary.BigEndian.PutUint16(b[k*8+6:], v.CatchType)
		}
		bs = append(bs, b...)
	}
	return bs
}
