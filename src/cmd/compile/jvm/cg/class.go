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
	thisClass    [2]byte
	superClass   [2]byte
	interfaces   []uint16
	fields       []*FieldInfo
	methods      []*MethodInfo
	attributes   []*AttributeInfo
}

func FromHighLevel(high *ClassHighLevel) *Class {
	c := &Class{}
	c.fromHighLevel(high)
	return c
}
func (c *Class) fromHighLevel(high *ClassHighLevel) {
	c.minorVersion = 0
	c.majorVersion = 52
	//int const
	for i, locations := range high.IntConsts {
		info := CONSTANT_Integer_info{}
		info.value = i
		c.ifConstPoolOverMaxSize()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool())
	}
	//long const
	for l, locations := range high.LongConsts {
		info := CONSTANT_Long_info{}
		info.value = l
		c.ifConstPoolOverMaxSize()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool(), nil)
	}
	//float const
	for f, locations := range high.FloatConsts {
		info := CONSTANT_Float_info{}
		info.value = f
		c.ifConstPoolOverMaxSize()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool())
	}
	//double const
	for d, locations := range high.DoubleConsts {
		info := CONSTANT_Double_info{}
		info.value = d
		c.ifConstPoolOverMaxSize()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info.ToConstPool(), nil)
	}
	//fieldref
	for f, locations := range high.FieldRefs {
		info := (&CONSTANT_Fieldref_info{}).ToConstPool()
		high.InsertClasses(f.Class, info.info[0:2])
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
		high.InsertNameAndType(CONSTANT_NameAndType_info_high_level{
			Name: f.Name,
			Type: f.Descriptor,
		}, info.info[2:4])
	}
	//methodref
	for m, locations := range high.MethodRefs {
		info := (&CONSTANT_Methodref_info{}).ToConstPool()
		high.InsertClasses(m.Class, info.info[0:2])
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
		high.InsertNameAndType(CONSTANT_NameAndType_info_high_level{
			Name: m.Name,
			Type: m.Descriptor,
		}, info.info[2:4])
	}
	//classess
	for cn, locations := range high.Classes {
		info := (&CONSTANT_Class_info{}).ToConstPool()
		high.InsertStringConst(cn, info.info[0:2])
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
	}
	//name and type
	for nt, locations := range high.NameAndTypes {
		info := (&CONSTANT_NameAndType_info{}).ToConstPool()
		high.InsertStringConst(nt.Name, info.info[0:2])
		high.InsertStringConst(nt.Type, info.info[2:4])
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
	}
	c.accessFlag = high.AccessFlags
	thisClassConst := (&CONSTANT_Class_info{}).ToConstPool()
	high.InsertStringConst(high.Name, thisClassConst.info[0:2])
	c.constPool = append(c.constPool, thisClassConst)
	superClassConst := (&CONSTANT_Class_info{}).ToConstPool()
	high.InsertStringConst(high.SuperClass, superClassConst.info[0:2])
	c.constPool = append(c.constPool, superClassConst)
	for _, i := range high.Interfaces {
		inter := (&CONSTANT_Class_info{}).ToConstPool()
		high.InsertStringConst(i, inter.info[0:2])
		index := c.constPoolUint16Length()
		c.interfaces = append(c.interfaces, index)
		c.constPool = append(c.constPool, inter)
	}
	for _, f := range high.Fields {
		field := &FieldInfo{}
		field.AccessFlags = f.AccessFlags
		high.InsertStringConst(f.Name, field.NameIndex[0:2])
		high.InsertStringConst(f.Descriptor, field.DescriptorIndex[0:2])
		c.fields = append(c.fields, field)
	}
	for _, ms := range high.Methods {
		for _, m := range ms {
			info := &MethodInfo{}
			info.AccessFlags = m.AccessFlags
			high.InsertStringConst(m.Name, info.nameIndex[0:2])
			high.InsertStringConst(m.Descriptor, info.descriptorIndex[0:2])
			codeinfo := m.Code.ToAttributeInfo()
			high.InsertStringConst("Code", codeinfo.nameIndex[0:2])
			m.Attributes = append(m.Attributes, codeinfo)
		}
	}
	//string consts
	for s, locations := range high.StringConsts {
		if len(s) > 65536 {
			panic("string length over 65536")
		}
		info := (&CONSTANT_Utf8_info{uint16(len(s)), []byte(s)}).ToConstPool()
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
	}
	c.ifConstPoolOverMaxSize()
	return
}
func (c *Class) constPoolUint16Length() uint16 {
	return uint16(len(c.constPool))
}
func (c *Class) ifConstPoolOverMaxSize() {
	if len(c.constPool) > 65535 {
		panic(fmt.Sprintf("const pool max size is:%d", 65535))
	}
}

func (c *Class) OutPut(dest io.Writer) error {
	c.dest = dest
	//magic number
	bs4 := make([]byte, 4)
	binary.BigEndian.PutUint32(bs4, 0xCAFEBABE)
	_, err := dest.Write(bs4)
	if err != nil {
		return err
	}
	// minorversion
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(c.minorVersion))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	// major version
	binary.BigEndian.PutUint16(bs2, uint16(c.majorVersion))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	//const pool
	binary.BigEndian.PutUint16(bs2, c.constPoolUint16Length())
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.constPool {
		_, err = dest.Write([]byte{byte(v.tag)})
		if err != nil {
			return err
		}
		_, err = dest.Write(v.info)
		if err != nil {
			return err
		}
	}
	//access flag
	binary.BigEndian.PutUint16(bs2, uint16(c.accessFlag))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	//this class
	_, err = dest.Write(c.thisClass[0:2])
	if err != nil {
		return err
	}
	//super Class
	_, err = dest.Write(c.superClass[0:2])
	if err != nil {
		return err
	}
	// interface
	binary.BigEndian.PutUint16(bs2, uint16(len(c.interfaces)))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.interfaces {
		binary.BigEndian.PutUint16(bs2, uint16(v))
		_, err = dest.Write(bs2)
		if err != nil {
			return err
		}
	}

	err = c.writeFields()
	if err != nil {
		return err
	}
	//methods
	err = c.writeMethods()
	if err != nil {
		return err
	}
	// attribute
	binary.BigEndian.PutUint16(bs2, uint16(len(c.attributes)))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	if len(c.attributes) > 0 {
		return c.writeAttributeInfo(c.attributes)
	}
	return nil
}

func (c *Class) writeMethods() error {
	var err error
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(c.methods)))
	_, err = c.dest.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.methods {
		binary.BigEndian.PutUint16(bs2, uint16(v.AccessFlags))
		_, err = c.dest.Write(bs2)
		if err != nil {
			return err
		}
		_, err = c.dest.Write(v.nameIndex[0:2])
		if err != nil {
			return err
		}
		_, err = c.dest.Write(v.descriptorIndex[0:2])
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(len(v.Attributes)))
		_, err = c.dest.Write(bs2)
		if err != nil {
			return err
		}
		if len(v.Attributes) > 0 {
			err = c.writeAttributeInfo(v.Attributes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Class) writeFields() error {
	var err error
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(c.fields)))
	_, err = c.dest.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.fields {
		binary.BigEndian.PutUint16(bs2, uint16(v.AccessFlags))
		_, err = c.dest.Write(bs2)
		if err != nil {
			return err
		}
		_, err = c.dest.Write(v.NameIndex[0:2])
		if err != nil {
			return err
		}
		_, err = c.dest.Write(v.DescriptorIndex[0:2])
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(len(v.Attributes)))
		_, err = c.dest.Write(bs2)
		if err != nil {
			return err
		}
		if len(v.Attributes) > 0 {
			err = c.writeAttributeInfo(v.Attributes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Class) writeAttributeInfo(as []*AttributeInfo) error {
	var err error
	bs2 := make([]byte, 2)
	bs4 := make([]byte, 4)
	for _, v := range as {
		_, err = c.dest.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint32(bs4, uint32(v.attributeLength))
		_, err = c.dest.Write(bs4)
		if err != nil {
			return err
		}
		_, err = c.dest.Write(v.info)
		if err != nil {
			return err
		}
	}
	return nil
}
