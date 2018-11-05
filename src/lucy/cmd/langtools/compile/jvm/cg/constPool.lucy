package cg

import (
	"encoding/binary"
	"math"
)

const (
	ConstantPoolTagUtf8               uint8 = 1
	ConstantPoolTagInteger            uint8 = 3
	ConstantPoolTagFloat              uint8 = 4
	ConstantPoolTagLong               uint8 = 5
	ConstantPoolTagDouble             uint8 = 6
	ConstantPoolTagClass              uint8 = 7
	ConstantPoolTagString             uint8 = 8
	ConstantPoolTagFieldref           uint8 = 9
	ConstantPoolTagMethodref          uint8 = 10
	ConstantPoolTagInterfaceMethodref uint8 = 11
	ConstantPoolTagNameAndType        uint8 = 12
	ConstantPoolTagMethodHandle       uint8 = 15
	ConstantPoolTagMethodType         uint8 = 16
	ConstantPoolTagInvokeDynamic      uint8 = 18
)

type ConstPool struct {
	selfIndex uint16 // using when it`s self
	Tag       uint8
	Info      []byte
}

type ConstantInfoClass struct {
	nameIndex uint16
}

func (c *ConstantInfoClass) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagClass
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, c.nameIndex)
	return p
}

type ConstantInfoString struct {
	index uint16
}

func (this *ConstantInfoString) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagString
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, this.index)
	return p
}

type ConstantInfoInteger struct {
	value int32
}

func (this *ConstantInfoInteger) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagInteger
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint32(p.Info, uint32(this.value))
	return p

}

type ConstantInfoFloat struct {
	value float32
}

func (this *ConstantInfoFloat) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagFloat
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint32(p.Info, math.Float32bits(this.value))
	return p
}

type ConstantInfoLong struct {
	value int64
}

/*
	big endian
*/
func (this *ConstantInfoLong) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagLong
	p.Info = make([]byte, 8)
	binary.BigEndian.PutUint64(p.Info, uint64(this.value))
	return p
}

type ConstantInfoDouble struct {
	value float64
}

func (this *ConstantInfoDouble) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagDouble
	p.Info = make([]byte, 8)
	binary.BigEndian.PutUint64(p.Info, math.Float64bits(this.value))
	return p
}

type ConstantInfoNameAndType struct {
	name, descriptor uint16
}

func (this *ConstantInfoNameAndType) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagNameAndType
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, this.name)
	binary.BigEndian.PutUint16(p.Info[2:], this.descriptor)
	return p
}

type ConstantInfoMethodref struct {
	classIndex, nameAndTypeIndex uint16
}

func (this *ConstantInfoMethodref) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagMethodref
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, this.classIndex)
	binary.BigEndian.PutUint16(p.Info[2:], this.nameAndTypeIndex)
	return p
}

type ConstantInfoInterfaceMethodref struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (this *ConstantInfoInterfaceMethodref) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagInterfaceMethodref
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, uint16(this.classIndex))
	binary.BigEndian.PutUint16(p.Info[2:], uint16(this.nameAndTypeIndex))
	return p
}

type ConstantInfoFieldref struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (this *ConstantInfoFieldref) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagFieldref
	p.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(p.Info, this.classIndex)
	binary.BigEndian.PutUint16(p.Info[2:], this.nameAndTypeIndex)
	return p
}

type ConstantInfoMethodHandle struct {
	referenceKind  uint8
	referenceIndex uint16
}

func (this *ConstantInfoMethodHandle) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagMethodHandle
	p.Info = make([]byte, 3)
	p.Info[0] = byte(this.referenceKind)
	binary.BigEndian.PutUint16(p.Info[1:], uint16(this.referenceIndex))
	return p
}

type ConstantInfoUtf8 struct {
	length uint16
	bs     []byte
}

func (this *ConstantInfoUtf8) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagUtf8
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, uint16(this.length))
	p.Info = append(p.Info, this.bs...)
	return p
}

type ConstantInfoMethodType struct {
	descriptorIndex uint16
}

func (this *ConstantInfoMethodType) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.Tag = ConstantPoolTagMethodType
	p.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Info, uint16(this.descriptorIndex))
	return p
}

type ConstantInfoInvokeDynamic struct {
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

func (this *ConstantInfoInvokeDynamic) ToConstPool() *ConstPool {
	info := &ConstPool{}
	info.Tag = ConstantPoolTagInvokeDynamic
	info.Info = make([]byte, 4)
	binary.BigEndian.PutUint16(info.Info, this.bootstrapMethodAttrIndex)
	binary.BigEndian.PutUint16(info.Info[2:], this.nameAndTypeIndex)
	return info
}
