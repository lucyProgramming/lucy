package cg

import "encoding/binary"

type AttributeCode struct {
	AttributeInfo
	maxStack             uint16
	maxLocals            uint16
	codeLength           uint32
	codes                []byte
	exceptionTableLength uint16
	exceptions           []*ExceptionTable
	attributeCounts      uint16
	attributes           []*AttributeInfo
}

type ExceptionTable struct {
	startPc   uint16
	endpc     uint16
	handlerPc uint16
	catchType uint16
}

func (a *AttributeCode) ToAttributeInfo() *AttributeInfo {
	ret := &AttributeInfo{}
	ret.info = make([]byte, 8)
	binary.BigEndian.PutUint16(ret.info, uint16(a.maxStack))
	binary.BigEndian.PutUint16(ret.info[2:], uint16(a.maxLocals))
	binary.BigEndian.PutUint32(ret.info[4:], uint32(len(a.codes)))
	ret.info = append(ret.info, a.codes...)
	ret.info = append(ret.info, a.mkExceptions()...)
	ret.info = append(ret.info, a.mkAttributes()...)
	ret.attributeLength = uint32(len(ret.info))
	return ret
}
func (a *AttributeCode) mkAttributes() []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.attributes)))
	if len(a.attributes) > 0 {
		b := make([]byte, 0)
		for _, v := range a.attributes {
			bb := make([]byte, 6)
			binary.BigEndian.PutUint16(bb, uint16(v.attributeIndex))
			binary.BigEndian.PutUint32(bb[2:], uint32(v.attributeLength))
			bb = append(bb, v.info...)
			b = append(b, bb...)
		}
		bs = append(bs, b...)
	}
	return bs
}

func (a *AttributeCode) mkExceptions() []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.exceptions)))
	if len(a.exceptions) > 0 {
		b := make([]byte, 8*len(a.exceptions))
		for k, v := range a.exceptions {
			binary.BigEndian.PutUint16(b[k*8:], uint16(v.startPc))
			binary.BigEndian.PutUint16(b[k*8+2:], uint16(v.endpc))
			binary.BigEndian.PutUint16(b[k*8+4:], uint16(v.handlerPc))
			binary.BigEndian.PutUint16(b[k*8+6:], uint16(v.catchType))
		}
		bs = append(bs, b...)
	}

	return bs
}
