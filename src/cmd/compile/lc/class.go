package lc

import (
	"encoding/binary"
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClassDecoder struct {
	bs  []byte
	ret *cg.Class
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
	var err error
	c.ret.AttributeGroupedByName, err = c.parseAttributes()
	return ret, err
}

func (c *ClassDecoder) parseConstPool() error {
	length := binary.BigEndian.Uint16(c.bs) - 1
	c.bs = c.bs[2:]
	for i := 0; i < int(length); i++ {
		switch c.bs[0] {
		case cg.CONSTANT_POOL_TAG_Utf8:
			p := &cg.ConstPool{}
			length := binary.BigEndian.Uint16(c.bs[1:])
			p.Tag = c.bs[0]
			c.bs = c.bs[3:]
			p.Info = c.bs[:length]
			c.bs = c.bs[length:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_Integer:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_Float:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_Fieldref:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_Methodref:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_InterfaceMethodref:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_NameAndType:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_InvokeDynamic:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:5]
			c.bs = c.bs[5:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_Long:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:9]
			c.bs = c.bs[9:]
			c.ret.ConstPool = append(c.ret.ConstPool, p, nil)
			i++ // increment twice
		case cg.CONSTANT_POOL_TAG_Double:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:9]
			c.bs = c.bs[9:]
			c.ret.ConstPool = append(c.ret.ConstPool, p, nil)
			i++ // increment twice
		case cg.CONSTANT_POOL_TAG_Class:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:3]
			c.bs = c.bs[3:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
		case cg.CONSTANT_POOL_TAG_String:
			p := &cg.ConstPool{}
			p.Tag = c.bs[0]
			p.Info = c.bs[1:3]
			c.bs = c.bs[3:]
			c.ret.ConstPool = append(c.ret.ConstPool, p)
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

func (c *ClassDecoder) parseFields() (err error) {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := uint16(0); i < length; i++ {
		f := &cg.FieldInfo{}
		f.AccessFlags = binary.BigEndian.Uint16(c.bs)
		f.NameIndex = binary.BigEndian.Uint16(c.bs[2:])
		f.DescriptorIndex = binary.BigEndian.Uint16(c.bs[4:])
		c.bs = c.bs[6:]
		f.AttributeGroupedByName, err = c.parseAttributes()
		if err != nil {
			return err
		}
		c.ret.Fields = append(c.ret.Fields, f)
	}
	return nil
}

func (c *ClassDecoder) parserMethods() (err error) {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := uint16(0); i < length; i++ {
		m := &cg.MethodInfo{}
		m.AccessFlags = binary.BigEndian.Uint16(c.bs)
		m.NameIndex = binary.BigEndian.Uint16(c.bs[2:])
		m.DescriptorIndex = binary.BigEndian.Uint16(c.bs[4:])
		c.bs = c.bs[6:]
		m.AttributeGroupedByName, err = c.parseAttributes()
		if err != nil {
			return err
		}
		c.ret.Methods = append(c.ret.Methods, m)
	}
	return nil
}

func (c *ClassDecoder) parseAttributes() (cg.AttributeGroupedByName, error) {
	ret := make(cg.AttributeGroupedByName)
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := uint16(0); i < length; i++ {
		a := &cg.AttributeInfo{}
		a.NameIndex = binary.BigEndian.Uint16(c.bs)
		if c.ret.ConstPool[a.NameIndex].Tag != cg.CONSTANT_POOL_TAG_Utf8 {
			return ret, fmt.Errorf("name index %d is not a utf8 const", a.NameIndex)
		}
		length := binary.BigEndian.Uint32(c.bs[2:])
		c.bs = c.bs[6:]
		a.Info = c.bs[:length]
		c.bs = c.bs[length:]
		name := string(c.ret.ConstPool[a.NameIndex].Info)
		if _, ok := ret[name]; ok {
			ret[name] = append(ret[name], a)
		} else {
			ret[name] = []*cg.AttributeInfo{a}
		}
	}
	return ret, nil
}
