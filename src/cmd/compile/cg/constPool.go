package cg

import (
	"encoding/binary"
	"math"
)

const (
	CONSTANT_POOL_TAG_Utf8               U1 = 1
	CONSTANT_POOL_TAG_Integer            U1 = 3
	CONSTANT_POOL_TAG_Float              U1 = 4
	CONSTANT_POOL_TAG_Long               U1 = 5
	CONSTANT_POOL_TAG_Double             U1 = 6
	CONSTANT_POOL_TAG_Class              U1 = 7
	CONSTANT_POOL_TAG_String             U1 = 8
	CONSTANT_POOL_TAG_Fieldref           U1 = 9
	CONSTANT_POOL_TAG_Methodref          U1 = 10
	CONSTANT_POOL_TAG_InterfaceMethodref U1 = 11
	CONSTANT_POOL_TAG_NameAndType        U1 = 12
	CONSTANT_POOL_TAG_MethodHandle       U1 = 15
	CONSTANT_POOL_TAG_MethodType         U1 = 16
	CONSTANT_POOL_TAG_InvokeDynamic      U1 = 18
)

type ToConstPool interface {
	ToConstPool() *ConstPool
}

type ConstPool struct {
	tag  U1
	info []byte
}

type CONSTANT_Class_info struct {
	ConstPool
	nameindex U2
}

func (c *CONSTANT_Class_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Class
	p.info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.info, uint16(c.nameindex))
	return p
}

type CONSTANT_String_info struct {
	ConstPool
	stringIndex U2
}

func (c *CONSTANT_String_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_String
	p.info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.info, uint16(c.stringIndex))
	return p
}

type CONSTANT_Integer_info struct {
	ConstPool
	value int32
}

func (c *CONSTANT_Integer_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Integer
	p.info = make([]byte, 4)
	p.info[0] = byte(c.value >> 24)
	p.info[1] = byte(c.value >> 16)
	p.info[2] = byte(c.value >> 8)
	p.info[3] = byte(c.value)
	return p

}

type CONSTANT_Float_info struct {
	ConstPool
	value float32
}

func (c *CONSTANT_Float_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Float
	p.info = make([]byte, 4)
	binary.BigEndian.PutUint32(p.info, math.Float32bits(c.value))
	return p
}

type CONSTANT_Long_info struct {
	ConstPool
	value uint64
}

func (c *CONSTANT_Long_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Long
	p.info = make([]byte, 8)
	p.info[0] = byte(c.value >> 56)
	p.info[1] = byte(c.value >> 48)
	p.info[2] = byte(c.value >> 40)
	p.info[3] = byte(c.value >> 32)
	p.info[4] = byte(c.value >> 24)
	p.info[5] = byte(c.value >> 16)
	p.info[6] = byte(c.value >> 8)
	p.info[7] = byte(c.value)
	return p
}

type CONSTANT_Double_info struct {
	ConstPool
	value float64
}

func (c *CONSTANT_Double_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Double
	p.info = make([]byte, 8)
	binary.BigEndian.PutUint64(p.info, math.Float64bits(c.value))
	return p
}

type CONSTANT_NameAndType_info struct {
	ConstPool
	nameIndex       U2
	descriptorIndex U2
}

func (c *CONSTANT_NameAndType_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_NameAndType
	p.info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.info, uint16(c.nameIndex))
	binary.BigEndian.PutUint16(p.info[2:], uint16(c.descriptorIndex))
	return p
}

type CONSTANT_Methodref_info struct {
	ConstPool
	classIndex       U2
	nameAndTypeIndex U2
}

func (c *CONSTANT_Methodref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Methodref
	p.info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.info, uint16(c.classIndex))
	binary.BigEndian.PutUint16(p.info[2:], uint16(c.nameAndTypeIndex))
	return p
}

type CONSTANT_InterfaceMethodref_info struct {
	ConstPool
	classIndex       U2
	nameAndTypeIndex U2
}

func (c *CONSTANT_InterfaceMethodref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_InterfaceMethodref
	p.info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.info, uint16(c.classIndex))
	binary.BigEndian.PutUint16(p.info[2:], uint16(c.nameAndTypeIndex))
	return p
}

type CONSTANT_Fieldref_info struct {
	ConstPool
	classIndex       U2
	nameAndTypeIndex U2
}

func (c *CONSTANT_Fieldref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Fieldref
	p.info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.info, uint16(c.classIndex))
	binary.BigEndian.PutUint16(p.info[2:], uint16(c.nameAndTypeIndex))
	return p
}

type CONSTANT_MethodHandle_info struct {
	ConstPool
	referenceKind  U1
	referenceIndex U2
}

func (c *CONSTANT_MethodHandle_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_MethodHandle
	p.info = make([]byte, 3)
	p.info[0] = byte(c.referenceKind)
	binary.BigEndian.PutUint16(p.info[1:], uint16(c.referenceIndex))
	return p
}

type CONSTANT_Utf8_info struct {
	ConstPool
	length U2
	bs     []byte
}

func (c *CONSTANT_Utf8_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Utf8
	p.info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.info, uint16(c.length))
	p.info = append(p.info, c.bs...)
	return p
}

type CONSTANT_MethodType_info struct {
	ConstPool
	descriptorIndex U2
}

func (c *CONSTANT_MethodType_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_MethodType
	p.info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.info, uint16(c.descriptorIndex))
	return p
}

type CONSTANT_InvokeDynamic_info struct {
	ConstPool
	bootstrapMethodAttrIndex U2
	nameAndTypeIndex         U2
}

func (c *CONSTANT_InvokeDynamic_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_InvokeDynamic
	p.info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.info, uint16(c.bootstrapMethodAttrIndex))
	binary.BigEndian.PutUint16(p.info[2:], uint16(c.nameAndTypeIndex))
	return p
}
