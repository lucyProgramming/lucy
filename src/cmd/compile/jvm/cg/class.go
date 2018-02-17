package cg

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	ACC_CLASS_PUBLIC     uint16 = 0x0001 // 可以被包的类外访问。
	ACC_CLASS_FINAL      uint16 = 0x0010 //不允许有子类。
	ACC_CLASS_SUPER      uint16 = 0x0020 //当用到invokespecial指令时，需要特殊处理③的父类方法。
	ACC_CLASS_INTERFACE  uint16 = 0x0200 // 标识定义的是接口而不是类。
	ACC_CLASS_ABSTRACT   uint16 = 0x0400 //  不能被实例化。
	ACC_CLASS_SYNTHETIC  uint16 = 0x1000 //标识并非Java源码生成的代码。
	ACC_CLASS_ANNOTATION uint16 = 0x2000 // 标识注解类型
	ACC_CLASS_ENUM       uint16 = 0x4000 // 标识枚举类型
)

type Class struct {
	dest         io.Writer
	magic        uint32 //0xCAFEBABE
	minorVersion uint16
	majorVersion uint16
	constPool    []*ConstPool
	accessFlag   uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16
	fields       []*FieldInfo
	methods      []*MethodInfo
	attributes   []*AttributeInfo
	Utf8Consts   map[string]*ConstPool
}

func (c *Class) insertUtfConst(s string) uint16 {
	if c.Utf8Consts == nil {
		c.Utf8Consts = make(map[string]*ConstPool)
	}
	if con, ok := c.Utf8Consts[s]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Utf8_info{uint16(len(s)), []byte(s)}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	return info.selfindex
}

func FromHighLevel(high *ClassHighLevel) *Class {
	c := &Class{}
	c.fromHighLevel(high)
	return c
}

func (c *Class) fromHighLevel(high *ClassHighLevel) {
	c.minorVersion = 0
	c.majorVersion = 49
	c.constPool = []*ConstPool{nil} // jvm const pool index begin at 1
	//int const
	for i, locations := range high.IntConsts {
		info := CONSTANT_Integer_info{}
		info.value = i
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool())
	}
	//long const
	for l, locations := range high.LongConsts {
		info := CONSTANT_Long_info{}
		info.value = l
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool(), nil)
	}
	//float const
	for f, locations := range high.FloatConsts {
		info := CONSTANT_Float_info{}
		info.value = f
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool())
	}
	//double const
	for d, locations := range high.DoubleConsts {
		info := CONSTANT_Double_info{}
		info.value = d
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool(), nil)
	}
	//fieldref
	for f, locations := range high.FieldRefs {
		info := (&CONSTANT_Fieldref_info{}).ToConstPool()
		high.InsertClassConst(f.Class, info.info[0:2])
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
		high.InsertNameAndTypeConst(CONSTANT_NameAndType_info_high_level{
			Name: f.Name,
			Type: f.Descriptor,
		}, info.info[2:4])
	}
	//methodref
	for m, locations := range high.MethodRefs {
		info := (&CONSTANT_Methodref_info{}).ToConstPool()
		high.InsertClassConst(m.Class, info.info[0:2])
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
		high.InsertNameAndTypeConst(CONSTANT_NameAndType_info_high_level{
			Name: m.Name,
			Type: m.Descriptor,
		}, info.info[2:4])
	}
	//classess
	for cn, locations := range high.Classes {
		info := (&CONSTANT_Class_info{c.insertUtfConst(cn)}).ToConstPool()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
	}
	//name and type
	for nt, locations := range high.NameAndTypes {
		info := (&CONSTANT_NameAndType_info{c.insertUtfConst(nt.Name), c.insertUtfConst(nt.Type)}).ToConstPool()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)

	}
	//string
	for s, locations := range high.StringConsts {
		info := (&CONSTANT_String_info{c.insertUtfConst(s)}).ToConstPool()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
		info.selfindex = index
	}

	c.accessFlag = high.AccessFlags
	thisClassConst := (&CONSTANT_Class_info{}).ToConstPool()
	binary.BigEndian.PutUint16(thisClassConst.info[0:2], c.insertUtfConst(high.Name))
	c.thisClass = c.constPoolUint16Length()
	c.constPool = append(c.constPool, thisClassConst)
	superClassConst := (&CONSTANT_Class_info{}).ToConstPool()
	binary.BigEndian.PutUint16(superClassConst.info[0:2], c.insertUtfConst(high.SuperClass))
	c.superClass = c.constPoolUint16Length()
	c.constPool = append(c.constPool, superClassConst)
	for _, i := range high.Interfaces {
		inter := (&CONSTANT_Class_info{c.insertUtfConst(i)}).ToConstPool()
		index := c.constPoolUint16Length()
		c.interfaces = append(c.interfaces, index)
		c.constPool = append(c.constPool, inter)
	}
	for _, f := range high.Fields {
		field := &FieldInfo{}
		field.AccessFlags = f.AccessFlags
		field.NameIndex = c.insertUtfConst(f.Name)
		field.DescriptorIndex = c.insertUtfConst(f.Descriptor)
		c.fields = append(c.fields, field)
	}
	for _, ms := range high.Methods {
		for _, m := range ms {
			info := &MethodInfo{}
			info.AccessFlags = m.AccessFlags //accessflag
			info.nameIndex = c.insertUtfConst(m.Name)
			info.descriptorIndex = c.insertUtfConst(m.Descriptor)
			codeinfo := m.Code.ToAttributeInfo(c)
			binary.BigEndian.PutUint16(codeinfo.nameIndex[0:2], c.insertUtfConst("Code"))
			info.Attributes = append(info.Attributes, codeinfo)
			c.methods = append(c.methods, info)
		}
	}

	//source file
	c.attributes = append(c.attributes, (&AttributeSourceFile{high.getSourceFile()}).ToAttributeInfo(c))
	c.ifConstPoolOverMaxSize()
	return
}
func (c *Class) constPoolUint16Length() uint16 {
	return uint16(len(c.constPool))
}
func (c *Class) ifConstPoolOverMaxSize() {
	if len(c.constPool) > CONSTANT_POOL_MAX_SIZE {
		panic(fmt.Sprintf("const pool max size is:%d", CONSTANT_POOL_MAX_SIZE))
	}
}
