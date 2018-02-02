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

type ToConstPool interface {
	ToConstPool() *ConstPool
}

type ConstPool struct {
	tag  uint8
	info []byte
}

type CONSTANT_Class_info struct {
}

func (c *CONSTANT_Class_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Class
	p.info = make([]byte, 2)
	return p
}

type CONSTANT_String_info struct {
}

func (c *CONSTANT_String_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_String
	p.info = make([]byte, 2)
	return p
}

type CONSTANT_Integer_info struct {
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
	value int64
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
}

func (c *CONSTANT_NameAndType_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_NameAndType
	p.info = make([]byte, 4)
	return p
}

type CONSTANT_Methodref_info struct {
}

func (c *CONSTANT_Methodref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Methodref
	p.info = make([]byte, 4)
	return p
}

type CONSTANT_InterfaceMethodref_info struct {
	classIndex       uint16
	nameAndTypeIndex uint16
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
}

func (c *CONSTANT_Fieldref_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_Fieldref
	p.info = make([]byte, 4)
	return p
}

type CONSTANT_MethodHandle_info struct {
	referenceKind  uint8
	referenceIndex uint16
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
	length uint16
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
	descriptorIndex uint16
}

func (c *CONSTANT_MethodType_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_POOL_TAG_MethodType
	p.info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.info, uint16(c.descriptorIndex))
	return p
}

//type CONSTANT_InvokeDynamic_info struct {
//	//	bootstrapMethodAttrIndex uint16
//	//	nameAndTypeIndex         uint16
//}

//func (c *CONSTANT_InvokeDynamic_info) ToConstPool() *ConstPool {
//	p := &ConstPool{}
//	p.tag = CONSTANT_POOL_TAG_InvokeDynamic
//	p.info = make([]byte, 4)

//	return p
//}
