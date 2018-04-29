package cg

import (
	"encoding/binary"
)

type AttributeDefaultParameters struct {
	Start  uint16 // first
	Consts []uint16
}

func (a *AttributeDefaultParameters) FromBytes(bs []byte) {
	a.Start = binary.BigEndian.Uint16(bs)
	bs = bs[2:]
	for len(bs) > 0 {
		a.Consts = append(a.Consts, binary.BigEndian.Uint16(bs))
		bs = bs[2:]
	}
}

func (a *AttributeDefaultParameters) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil && len(a.Consts) == 0 {
		return nil
	}
	info := &AttributeInfo{}
	info.NameIndex = class.insertUtf8Const(ATTRIBUTE_NAME_LUCY_DEFAULT_PARAMETERS)
	info.attributeLength = uint32(2 * (1 + len(a.Consts)))
	info.Info = make([]byte, info.attributeLength)
	binary.BigEndian.PutUint16(info.Info, a.Start)
	for i := 0; i < len(a.Consts); i++ {
		binary.BigEndian.PutUint16(info.Info[(i+1)*2:], a.Consts[i])
	}
	return info
}
