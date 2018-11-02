package cg

import (
	"encoding/binary"
)

const (
	MethodParameterTypeAccFinal     = 0x0010
	MethodParameterTypeAccSynthetic = 0x1000
	MethodParameterTypeAccMandated  = 0x8000
)

type AttributeMethodParameters struct {
	Parameters []*MethodParameter
}

type MethodParameter struct {
	Name        string
	AccessFlags uint16
}

func (a *AttributeMethodParameters) FromBs(class *Class, bs []byte) {
	if len(bs) != int(bs[0])*4+1 {
		panic("impossible")
	}
	bs = bs[1:]
	for len(bs) > 0 {
		p := &MethodParameter{}
		p.Name = string(class.ConstPool[binary.BigEndian.Uint16(bs)].Info)
		p.AccessFlags = binary.BigEndian.Uint16(bs[2:])
		a.Parameters = append(a.Parameters, p)
		bs = bs[4:]
	}
}

func (a *AttributeMethodParameters) ToAttributeInfo(class *Class, attrName ...string) *AttributeInfo {
	if a == nil || len(a.Parameters) == 0 {
		return nil
	}
	ret := &AttributeInfo{}
	if len(attrName) > 0 {
		ret.NameIndex = class.InsertUtf8Const(attrName[0])
	} else {
		ret.NameIndex = class.InsertUtf8Const(AttributeNameMethodParameters)
	}
	ret.attributeLength = uint32(len(a.Parameters)*4 + 1)
	ret.Info = make([]byte, ret.attributeLength)
	ret.Info[0] = byte(len(a.Parameters))
	for k, v := range a.Parameters {
		binary.BigEndian.PutUint16(ret.Info[4*k+1:], class.InsertUtf8Const(v.Name))
		binary.BigEndian.PutUint16(ret.Info[4*k+3:], v.AccessFlags)
	}
	return ret
}
