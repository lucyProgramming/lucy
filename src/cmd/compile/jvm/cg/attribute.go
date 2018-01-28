package cg

import "encoding/binary"

type AttributeInfo struct {
	nameIndex       [2]byte
	attributeLength uint32
	info            []byte
}

type ToAttributeInfo interface {
	ToAttributeInfo() *AttributeInfo
}

type AttributeSignature struct {
	index uint16
}

func (a *AttributeSignature) ToAttributeInfo() *AttributeInfo {
	ret := &AttributeInfo{}
	ret.attributeLength = 2
	ret.info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.info, uint16(a.index))
	return ret
}

type AttributeSourceFile struct {
	index uint16
}

func (a *AttributeSourceFile) ToAttributeInfo() *AttributeInfo {
	ret := &AttributeInfo{}
	ret.attributeLength = 2
	ret.info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.info, uint16(a.index))
	return ret
}

type AttributeLineNumber struct {
	linenumbers []*AttributeLinePc
}

func (a *AttributeLineNumber) ToAttributeInfo() *AttributeInfo {
	ret := &AttributeInfo{}
	ret.attributeLength = uint32(len(a.linenumbers)) * 4
	ret.info = make([]byte, ret.attributeLength)
	for k, v := range a.linenumbers {
		binary.BigEndian.PutUint16(ret.info[k*4:], uint16(v.startPc))
		binary.BigEndian.PutUint16(ret.info[k*4+2:], uint16(v.lineNumber))
	}
	return ret
}

type AttributeLinePc struct {
	startPc    uint16
	lineNumber uint16
}

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
