package cg

import (
	"encoding/binary"
)

type AttributeBootstrapMethods struct {
	numBootStrapMethods uint16
	methods             []*BootStrapMethod
}

func (a *AttributeBootstrapMethods) ToAttributeInfo() *AttributeInfo {
	ret := &AttributeInfo{}
	ret.info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.info, uint16(a.numBootStrapMethods))
	for _, v := range a.methods {
		bs := make([]byte, 4+len(v.bootStrapMethodArguments)*2)
		binary.BigEndian.PutUint16(bs, uint16(v.bootStrapMethodRef))
		binary.BigEndian.PutUint16(bs[2:], uint16(v.numBootStrapMethodArgument))
		for kk, vv := range v.bootStrapMethodArguments {
			binary.BigEndian.PutUint16(bs[4+kk*2:], uint16(vv))
		}
		ret.info = append(ret.info, bs...)
	}
	ret.attributeLength = uint32(len(ret.info))
	return ret
}

type BootStrapMethod struct {
	bootStrapMethodRef         uint16
	numBootStrapMethodArgument uint16
	bootStrapMethodArguments   []uint16
}
