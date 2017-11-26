package cg

import "encoding/binary"

type AttributeCode struct {
	AttributeInfo
	maxStack             U2
	maxLocals            U2
	codeLength           U4
	codes                []byte
	exceptionTableLength U2
	exceptions           []*ExceptionTable
	attributeCounts      U2
	attributes           []*AttributeInfo
}

type ExceptionTable struct {
	startPc   U2
	endpc     U2
	handlerPc U2
	catchType U2
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
	ret.attributeLength = U4(len(ret.info))
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
