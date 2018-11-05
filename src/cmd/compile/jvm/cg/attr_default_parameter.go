package cg

import (
	"encoding/binary"
)

type AttributeDefaultParameters struct {
	Start     uint16 // start
	Constants []uint16
}

func (this *AttributeDefaultParameters) FromBytes(bs []byte) {
	this.Start = binary.BigEndian.Uint16(bs)
	bs = bs[2:]
	for len(bs) > 0 {
		this.Constants = append(this.Constants, binary.BigEndian.Uint16(bs))
		bs = bs[2:]
	}
}

func (this *AttributeDefaultParameters) ToAttributeInfo(class *Class) *AttributeInfo {
	if this == nil && len(this.Constants) == 0 {
		return nil
	}
	info := &AttributeInfo{}
	info.NameIndex = class.InsertUtf8Const(AttributeNameLucyDefaultParameters)
	info.attributeLength = uint32(2 * (1 + len(this.Constants)))
	info.Info = make([]byte, info.attributeLength)
	binary.BigEndian.PutUint16(info.Info, this.Start)
	for i := 0; i < len(this.Constants); i++ {
		binary.BigEndian.PutUint16(info.Info[(i+1)*2:], this.Constants[i])
	}
	return info
}
