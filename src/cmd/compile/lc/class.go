package lc

import (
	"encoding/binary"
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type ClassDecoder struct {
	bs  []byte
	ret *cg.Class
}

func (c *ClassDecoder) parseConstPool() error {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := 0; i < int(length)-1; i++ {
		switch c.bs[0] {
		case cg.CONSTANT_POOL_TAG_Utf8:
			p := &cg.ConstPool{}
			length := binary.BigEndian.Uint16(c.bs[1:])
			p.Tag = c.bs[0]
			c.bs = c.bs[3:]
			p.Info = c.bs[:length]
			c.bs = c.bs[length:]
			fmt.Println(len(c.ret.ConstPool))
			c.ret.ConstPool = append(c.ret.ConstPool, p)
			fmt.Println(string(p.Info))
		case cg.CONSTANT_POOL_TAG_Integer:
			fallthrough
		case cg.CONSTANT_POOL_TAG_Float:
			fallthrough
		case cg.CONSTANT_POOL_TAG_Fieldref:
			fallthrough
		case cg.CONSTANT_POOL_TAG_Methodref:
			fallthrough
		case cg.CONSTANT_POOL_TAG_InterfaceMethodref:
			fallthrough
		case cg.CONSTANT_POOL_TAG_InvokeDynamic:
			fallthrough
		case cg.CONSTANT_POOL_TAG_NameAndType:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_Long:
			fallthrough
		case cg.CONSTANT_POOL_TAG_Double:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:9]
			c.bs = c.bs[9:]
			c.ret.ConstPool = append(c.ret.ConstPool, p, nil)
		case cg.CONSTANT_POOL_TAG_Class:
			fallthrough
		case cg.CONSTANT_POOL_TAG_String:
			fallthrough
		case cg.CONSTANT_POOL_TAG_MethodType:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:3]
			c.bs = c.bs[3:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_MethodHandle:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:4]
			c.bs = c.bs[4:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		}
	}
	return nil
}
func (c *ClassDecoder) parseInterfaces() {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := uint16(0); i < length; i++ {
		c.ret.Interfaces = append(c.ret.Interfaces, binary.BigEndian.Uint16(c.bs))
		c.bs = c.bs[2:]
	}
}
func (c *ClassDecoder) parseFields() {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := uint16(0); i < length; i++ {
		f := &cg.FieldInfo{}
		f.AccessFlags = binary.BigEndian.Uint16(c.bs)
		f.NameIndex = binary.BigEndian.Uint16(c.bs[2:])
		f.DescriptorIndex = binary.BigEndian.Uint16(c.bs[4:])
		c.bs = c.bs[6:]
		fmt.Println("!!!!!!!!!!!!", string(c.ret.ConstPool[f.NameIndex].Info))
		f.Attributes = c.parseAttributes()
		c.ret.Fields = append(c.ret.Fields, f)
	}
}

func (c *ClassDecoder) parserMethods() {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := uint16(0); i < length; i++ {
		m := &cg.MethodInfo{}
		m.AccessFlags = binary.BigEndian.Uint16(c.bs)
		m.NameIndex = binary.BigEndian.Uint16(c.bs[2:])
		m.DescriptorIndex = binary.BigEndian.Uint16(c.bs[4:])
		c.bs = c.bs[6:]
		m.Attributes = c.parseAttributes()
		c.ret.Methods = append(c.ret.Methods, m)
	}
}

func (c *ClassDecoder) parseAttributes() []*cg.AttributeInfo {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	ret := []*cg.AttributeInfo{}
	for i := uint16(0); i < length; i++ {
		a := &cg.AttributeInfo{}
		a.NameIndex = binary.BigEndian.Uint16(c.bs)
		length := binary.BigEndian.Uint32(c.bs[2:])
		c.bs = c.bs[6:]
		a.Info = c.bs[:length]
		if string(c.ret.ConstPool[a.NameIndex].Info) == cg.CONSTANT_SOURCE_FILE { //sepcial case
			index := binary.BigEndian.Uint16(a.Info)
			c.ret.SourceFile = string(c.ret.ConstPool[index].Info)
		}
		c.bs = c.bs[length:]
		ret = append(ret, a)
	}
	return ret
}

func (c *ClassDecoder) decode(bs []byte) (*cg.Class, error) {
	c.bs = bs
	if binary.BigEndian.Uint32(bs) != cg.CLASS_MAGIC {
		return nil, fmt.Errorf("magic number is not right")
	}
	c.bs = c.bs[4:]
	ret := &cg.Class{}
	c.ret = ret
	//version
	ret.MinorVersion = binary.BigEndian.Uint16(c.bs)
	ret.MajorVersion = binary.BigEndian.Uint16(c.bs[2:])
	c.bs = c.bs[4:]
	ret.ConstPool = []*cg.ConstPool{nil} // pool start 1

	//const pool
	if err := c.parseConstPool(); err != nil {
		return ret, err
	}
	//access flag
	ret.AccessFlag = binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	// this class
	ret.ThisClass = binary.BigEndian.Uint16(c.bs)
	ret.SuperClass = binary.BigEndian.Uint16(c.bs[2:])
	c.bs = c.bs[4:]
	c.parseInterfaces()
	c.parseFields()
	c.parserMethods()

	c.ret.Attributes = c.parseAttributes()

	return ret, nil
}
