package cg

import (
	"encoding/binary"
	"math"
)

const (
	CONSTANT_POOL_TAG_Utf8               uint8 = 1
	CONSTANT_POOL_TAG_Integer            uint8 = 3
	CONSTANT_POOL_TAG_Float              uint8 = 4
	CONSTANT_POOL_TAG_Long               uint8 = 5
	CONSTANT_POOL_TAG_Double             uint8 = 6
	CONSTANT_POOL_TAG_Class              uint8 = 7
	CONSTANT_POOL_TAG_String             uint8 = 8
	CONSTANT_POOL_TAG_Fieldref           uint8 = 9
	CONSTANT_POOL_TAG_Methodref          uint8 = 10
	CONSTANT_POOL_TAG_InterfaceMethodref uint8 = 11
	CONSTANT_POOL_TAG_NameAndType        uint8 = 12
	CONSTANT_POOL_TAG_MethodHandle       uint8 = 15
	CONSTANT_POOL_TAG_MethodType         uint8 = 16
	CONSTANT_POOL_TAG_InvokeDynamic      uint8 = 18
)

type ConstPool struct {
	selfIndex uint16 // using when it`s self
	Tag       uint8
	Info      []byte
}

type CONSTANT_Class_info struct {
	nameIndex uint16
}

func (c *CONSTANT_Class_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Class
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, c.nameIndex)
	return p
}

type CONSTANT_String_info struct {
	index uint16
}

func (c *CONSTANT_String_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_String
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, c.index)
	return p
}

type CONSTANT_Integer_info struct {
	value int32
}

func (c *CONSTANT_Integer_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Integer
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint32(p.Info, uint32(c.value))
	return p

}

type CONSTANT_Float_info struct {
	value float32
}

func (c *CONSTANT_Float_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Float
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint32(p.Info, math.Float32bits(c.value))
	return p
}

type CONSTANT_Long_info struct {
	value int64
}

/*
	big endian
*/
func (c *CONSTANT_Long_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Long
	p.Info = make([]byte, 8)
	binary.BigEndian.PutUint64(p.Info, uint64(c.value))
	return p
}

type CONSTANT_Double_info struct {
	value float64
}

func (c *CONSTANT_Double_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Double
	p.Info = make([]byte, 8)
	binary.BigEndian.PutUint64(p.Info, math.Float64bits(c.value))
	return p
}

type CONSTANT_NameAndType_info struct {
	name, descriptor uint16
}

func (c *CONSTANT_NameAndType_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_NameAndType
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, c.name)
	binary.BigEndian.PutUint16(p.Info[2:], c.descriptor)
	return p
}

type CONSTANT_Methodref_info struct {
	classIndex, nameAndTypeIndex uint16
}

func (c *CONSTANT_Methodref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Methodref
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, c.classIndex)
	binary.BigEndian.PutUint16(p.Info[2:], c.nameAndTypeIndex)
	return p
}

type CONSTANT_InterfaceMethodref_info struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (c *CONSTANT_InterfaceMethodref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_InterfaceMethodref
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, uint16(c.classIndex))
	binary.BigEndian.PutUint16(p.Info[2:], uint16(c.nameAndTypeIndex))
	return p
}

type CONSTANT_Fieldref_info struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (c *CONSTANT_Fieldref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Fieldref
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, c.classIndex)
	binary.BigEndian.PutUint16(p.Info[2:], c.nameAndTypeIndex)
	return p
}

type CONSTANT_MethodHandle_info struct {
	referenceKind  uint8
	referenceIndex uint16
}

func (c *CONSTANT_MethodHandle_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_MethodHandle
	p.Info = make([]byte, 3)
	p.Info[0] = byte(c.referenceKind)
	binary.BigEndian.PutUint16(p.Info[1:], uint16(c.referenceIndex))
	return p
}

type CONSTANT_Utf8_info struct {
	length uint16
	bs     []byte
}

func (c *CONSTANT_Utf8_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_Utf8
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, uint16(c.length))
	p.Info = append(p.Info, c.bs...)
	return p
}

type CONSTANT_MethodType_info struct {
	descriptorIndex uint16
}

func (c *CONSTANT_MethodType_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = CONSTANT_POOL_TAG_MethodType
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, uint16(c.descriptorIndex))
	return p
}

type CONSTANT_InvokeDynamic_info struct {
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

func (c *CONSTANT_InvokeDynamic_info) ToConstPool() *ConstPool {
	info := &ConstPool{}
	info.Tag = CONSTANT_POOL_TAG_InvokeDynamic
	info.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(info.Info, c.bootstrapMethodAttrIndex)
	binary.BigEndian.PutUint16(info.Info[2:], c.nameAndTypeIndex)
	return info
}
