package cg

import (
	"encoding/binary"
)

type AttributeDefaultParameters struct {
	Start     uint16 // start
	Constants []uint16
}

func (a *AttributeDefaultParameters) FromBytes(bs []byte) {
	a.Start = binary.BigEndian.Uint16(bs)
	bs = bs[2:]
	for len(bs) > 0 {
		a.Constants = append(a.Constants, binary.BigEndian.Uint16(bs))
		bs = bs[2:]
	}
}

func (a *AttributeDefaultParameters) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil && len(a.Constants) == 0 {
		return nil
	}
	info := &AttributeInfo{}
	info.NameIndex = class.InsertUtf8Const(AttributeNameLucyDefaultParameters)
	info.attributeLength = uint32(2 * (1 + len(a.Constants)))
	info.Info = make([]byte, info.attributeLength)
	binary.BigEndian.PutUint16(info.Info, a.Start)
	for i := 0; i < len(a.Constants); i++ {
		binary.BigEndian.PutUint16(info.Info[(i+1)*2:], a.Constants[i])
	}
	return info
}
