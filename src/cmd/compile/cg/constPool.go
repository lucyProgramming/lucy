package cg

import "encoding/binary"

const (
	CONSTANT_Utf8               U1 = 1
	CONSTANT_Integer            U1 = 3
	CONSTANT_Float              U1 = 4
	CONSTANT_Long               U1 = 5
	CONSTANT_Double             U1 = 6
	CONSTANT_Class              U1 = 7
	CONSTANT_String             U1 = 8
	CONSTANT_Fieldref           U1 = 9
	CONSTANT_Methodref          U1 = 10
	CONSTANT_InterfaceMethodref U1 = 11
	CONSTANT_NameAndType        U1 = 12
	CONSTANT_MethodHandle       U1 = 15
	CONSTANT_MethodType         U1 = 16
	CONSTANT_InvokeDynamic      U1 = 18
)

type ToConstPool interface {
	ToConstPool() *ConstPool
}

type ConstPool struct {
	tag  U1
	info []byte
}

type CONSTANT_String_info struct {
	ConstPool
	stringIndex U2
}

type CONSTANT_Integer_info struct {
	ConstPool
	bytes [4]byte
}

type CONSTANT_Float_info struct {
	ConstPool
	bytes [4]byte
}

type CONSTANT_Long_info struct {
	ConstPool
	highBytes [4]byte
	logBytes  [4]byte
}

type CONSTANT_Double_info struct {
	ConstPool
	highBytes [4]byte
	logBytes  [4]byte
}

type CONSTANT_NameAndType_info struct {
	ConstPool
	nameIndex       U2
	descriptorIndex U2
}

type CONSTANT_Class_info struct {
	ConstPool
	nameindex U2
}

type CONSTANT_Methodref_info struct {
	ConstPool
	classIndex       U2
	nameAndTypeIndex U2
}

type CONSTANT_InterfaceMethodref_info struct {
	ConstPool
	classIndex       U2
	namdAndTypeIndex U2
}

type CONSTANT_Fieldref_info struct {
	ConstPool
	classIndex       U2
	nameAndTypeIndex U2
}

type CONSTANT_MethodHandle_info struct {
	ConstPool
	referenceKind  uint8
	referenceIndex U2
}

type CONSTANT_Utf8_info struct {
	ConstPool
	length U2
	bs     []byte
}

type CONSTANT_MethodType_info struct {
	ConstPool
	descriptorIndex U2
}

type CONSTANT_InvokeDynamic_info struct {
	ConstPool
	bootstrapMethodAttrIndex U2
	nameAndTypeIndex         U2
}

func (c *CONSTANT_Class_info) ToConstPool() *ConstPool {
	p := &ConstPool{}
	p.tag = CONSTANT_Class
	p.info = make([]byte, 2)
	binary.BigEndian.PutUint16(p.info, uint16(c.nameindex))
	return p
}
