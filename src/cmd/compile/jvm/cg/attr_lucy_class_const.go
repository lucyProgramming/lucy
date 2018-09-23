package cg

import "encoding/binary"

type AttributeLucyClassConst struct {
	Constants []*LucyClassConst
}

type LucyClassConst struct {
	Name       string
	Descriptor string
	ValueIndex uint16
}

func (a *AttributeLucyClassConst) FromBs(class *Class, bs []byte) {
	for len(bs) > 0 {
		constant := &LucyClassConst{}
		constant.Name = string(class.ConstPool[binary.BigEndian.Uint16(bs)].Info)
		constant.Descriptor = string(class.ConstPool[binary.BigEndian.Uint16(bs[2:])].Info)
		constant.ValueIndex = binary.BigEndian.Uint16(bs[4:])
		bs = bs[6:]
		a.Constants = append(a.Constants, constant)
	}
}

func (a *AttributeLucyClassConst) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil || len(a.Constants) == 0 {
		return nil
	}
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyClassConst)
	ret.Info = make([]byte, len(a.Constants)*6)
	for k, v := range a.Constants {
		b := ret.Info[k*6:]
		binary.BigEndian.PutUint16(b, class.InsertUtf8Const(v.Name))
		binary.BigEndian.PutUint16(b[2:], class.InsertUtf8Const(v.Descriptor))
		binary.BigEndian.PutUint16(b[4:], v.ValueIndex)
	}
	ret.attributeLength = uint32(len(ret.Info))
	return ret
}
