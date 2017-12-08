package cg

import (
	"encoding/binary"
	"os"
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
	ACC_VOLATILE         uint16 = 0x0040 //volatile，表示字段是易变的。
	ACC_TRANSIENT        uint16 = 0x0080 //transient，表示字段不会被序列化。
	ACC_SYNTHETIC        uint16 = 0x1000 //表示字段由编译器自动产生。
	ACC_ENUM             uint16 = 0x4000 //enum，表示字段为枚举类型。
)

type Class struct {
	f              *os.File
	magic          uint32 //0xCAFEBABE
	minorVersion   uint16
	majorVersion   uint16
	constPoolCount uint16
	constPool      []*ConstPool
	accessFlag     uint16
	thisClass      uint16
	superClass     uint16
	interfaceCount uint16
	interfaces     []uint16
	fieldCount     uint16
	fields         []*FieldInfo
	methodCount    uint16
	methods        []*MethodInfo
	attributeCount uint16
	attributes     []*AttributeInfo
}

func (c *Class) OutPut(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	//magic number
	bs4 := make([]byte, 4)
	binary.BigEndian.PutUint32(bs4, c.magic)
	_, err = f.Write(bs4)
	if err != nil {
		return err
	}
	// minorversion
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(c.minorVersion))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	// major version
	binary.BigEndian.PutUint16(bs2, uint16(c.majorVersion))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	//const pool
	binary.BigEndian.PutUint16(bs2, uint16(c.constPoolCount))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.constPool {
		_, err = f.Write([]byte{byte(v.tag)})
		if err != nil {
			return err
		}
		_, err = f.Write(v.info)
		if err != nil {
			return err
		}
	}
	//access flag
	binary.BigEndian.PutUint16(bs2, uint16(c.accessFlag))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	//this class
	binary.BigEndian.PutUint16(bs2, uint16(c.thisClass))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	//super class
	binary.BigEndian.PutUint16(bs2, uint16(c.superClass))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	// interface
	binary.BigEndian.PutUint16(bs2, uint16(c.interfaceCount))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.interfaces {
		binary.BigEndian.PutUint16(bs2, uint16(v))
		_, err = f.Write(bs2)
		if err != nil {
			return err
		}
	}
	//fields
	c.f = f
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

	binary.BigEndian.PutUint16(bs2, uint16(c.attributeCount))
	_, err = f.Write(bs2)
	if err != nil {
		return err
	}
	if c.attributeCount > 0 {
		return c.writeAttributeInfo(c.attributes)
	}

	return nil
}

/*

type MethodInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributeCount  uint16
	attributes      []*AttributeInfo
}
*/

func (c *Class) writeMethods() error {
	var err error
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(c.methodCount))
	_, err = c.f.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.methods {
		binary.BigEndian.PutUint16(bs2, uint16(v.accessFlags))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(v.nameIndex))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(v.descriptorIndex))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(v.attributeCount))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		if v.attributeCount > 0 {
			err = c.writeAttributeInfo(v.attributes)
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
	binary.BigEndian.PutUint16(bs2, uint16(c.fieldCount))
	_, err = c.f.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.fields {
		binary.BigEndian.PutUint16(bs2, uint16(v.accessFlags))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(v.nameIndex))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(v.descriptorIndex))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, uint16(v.attributeCount))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		if v.attributeCount > 0 {
			err = c.writeAttributeInfo(v.attributes)
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
		binary.BigEndian.PutUint16(bs2, uint16(v.attributeIndex))
		_, err = c.f.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint32(bs4, uint32(v.attributeIndex))
		_, err = c.f.Write(bs4)
		if err != nil {
			return err
		}
		_, err = c.f.Write(v.info)
		if err != nil {
			return err
		}
	}
	return nil
}
