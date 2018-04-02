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
	dest                     io.Writer
	magic                    uint32 //0xCAFEBABE
	MinorVersion             uint16
	MajorVersion             uint16
	ConstPool                []*ConstPool
	AccessFlag               uint16
	ThisClass                uint16
	SuperClass               uint16
	Interfaces               []uint16
	Fields                   []*FieldInfo
	Methods                  []*MethodInfo
	Attributes               []*AttributeInfo
	Utf8Consts               map[string]*ConstPool
	IntConsts                map[int32]*ConstPool
	LongConsts               map[int64]*ConstPool
	FloatConsts              map[float32]*ConstPool
	DoubleConsts             map[float64]*ConstPool
	ClassConsts              map[string]*ConstPool
	StringConsts             map[string]*ConstPool
	FieldRefConsts           map[CONSTANT_Fieldref_info_high_level]*ConstPool
	NameAndTypeConsts        map[CONSTANT_NameAndType_info_high_level]*ConstPool
	MethodrefConsts          map[CONSTANT_Methodref_info_high_level]*ConstPool
	InterfaceMethodrefConsts map[CONSTANT_InterfaceMethodref_info_high_level]*ConstPool
	AttributeClosureClass    *AttributeClosureFunctionClass
	AttributeGroupedByName   AttributeGroupedByName
	TypeAlias                []*AttributeLucyTypeAlias
}

func (c *Class) InsertInterfaceMethodrefConst(n CONSTANT_InterfaceMethodref_info_high_level) uint16 {
	if c.InterfaceMethodrefConsts == nil {
		c.InterfaceMethodrefConsts = make(map[CONSTANT_InterfaceMethodref_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.InterfaceMethodrefConsts[n]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_InterfaceMethodref_info{
		classIndex: c.InsertClassConst(n.Class),
		nameAndTypeIndex: c.InsertNameAndConst(CONSTANT_NameAndType_info_high_level{
			Name:       n.Method,
			Descriptor: n.Descriptor,
		}),
	}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.InterfaceMethodrefConsts[n] = info
	return info.selfindex
}

func (c *Class) InsertMethodrefConst(n CONSTANT_Methodref_info_high_level) uint16 {
	if c.MethodrefConsts == nil {
		c.MethodrefConsts = make(map[CONSTANT_Methodref_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.MethodrefConsts[n]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Methodref_info{
		classIndex: c.InsertClassConst(n.Class),
		nameAndTypeIndex: c.InsertNameAndConst(CONSTANT_NameAndType_info_high_level{
			Name:       n.Method,
			Descriptor: n.Descriptor,
		}),
	}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.MethodrefConsts[n] = info
	return info.selfindex
}

func (c *Class) InsertNameAndConst(n CONSTANT_NameAndType_info_high_level) uint16 {
	if c.NameAndTypeConsts == nil {
		c.NameAndTypeConsts = make(map[CONSTANT_NameAndType_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.NameAndTypeConsts[n]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_NameAndType_info{
		name:       c.insertUtfConst(n.Name),
		descriptor: c.insertUtfConst(n.Descriptor),
	}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.NameAndTypeConsts[n] = info
	return info.selfindex
}
func (c *Class) InsertFieldRefConst(f CONSTANT_Fieldref_info_high_level) uint16 {
	if c.FieldRefConsts == nil {
		c.FieldRefConsts = make(map[CONSTANT_Fieldref_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.FieldRefConsts[f]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Fieldref_info{
		classIndex:       c.InsertClassConst(f.Class),
		nameAndTypeIndex: c.InsertNameAndConst(CONSTANT_NameAndType_info_high_level{f.Field, f.Descriptor}),
	}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.FieldRefConsts[f] = info
	return info.selfindex
}
func (c *Class) insertUtfConst(s string) uint16 {
	if c.Utf8Consts == nil {
		c.Utf8Consts = make(map[string]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.Utf8Consts[s]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Utf8_info{uint16(len(s)), []byte(s)}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.Utf8Consts[s] = info
	return info.selfindex
}

func (c *Class) InsertIntConst(i int32) uint16 {
	if c.IntConsts == nil {
		c.IntConsts = make(map[int32]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.IntConsts[i]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Integer_info{i}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.IntConsts[i] = info
	return info.selfindex
}
func (c *Class) InsertLongConst(i int64) uint16 {
	if c.LongConsts == nil {
		c.LongConsts = make(map[int64]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.LongConsts[i]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Long_info{i}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info, nil)
	c.LongConsts[i] = info
	return info.selfindex
}

func (c *Class) InsertFloatConst(f float32) uint16 {
	if c.FloatConsts == nil {
		c.FloatConsts = make(map[float32]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.FloatConsts[f]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Float_info{f}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.FloatConsts[f] = info
	return info.selfindex
}

func (c *Class) InsertDoubleConst(f float64) uint16 {
	if c.DoubleConsts == nil {
		c.DoubleConsts = make(map[float64]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.DoubleConsts[f]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Double_info{f}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info, nil)
	c.DoubleConsts[f] = info
	return info.selfindex
}

func (c *Class) InsertClassConst(name string) uint16 {
	if c.ClassConsts == nil {
		c.ClassConsts = make(map[string]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.ClassConsts[name]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Class_info{c.insertUtfConst(name)}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.ClassConsts[name] = info
	return info.selfindex
}

func (c *Class) InsertStringConst(s string) uint16 {
	if c.StringConsts == nil {
		c.StringConsts = make(map[string]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.StringConsts[s]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_String_info{c.insertUtfConst(s)}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.StringConsts[s] = info
	return info.selfindex
}

func (high *ClassHighLevel) FromHighLevel() *Class {
	high.Class.fromHighLevel(high)
	return &high.Class
}

func (c *Class) fromHighLevel(high *ClassHighLevel) {
	c.MinorVersion = 0
	c.MajorVersion = 49
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil} // jvm const pool index begin at 1
	}
	c.AccessFlag = high.AccessFlags
	thisClassConst := (&CONSTANT_Class_info{}).ToConstPool()
	binary.BigEndian.PutUint16(thisClassConst.Info[0:2], c.insertUtfConst(high.Name))
	c.ThisClass = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, thisClassConst)
	superClassConst := (&CONSTANT_Class_info{}).ToConstPool()
	binary.BigEndian.PutUint16(superClassConst.Info[0:2], c.insertUtfConst(high.SuperClass))
	c.SuperClass = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, superClassConst)
	for _, i := range high.Interfaces {
		inter := (&CONSTANT_Class_info{c.insertUtfConst(i)}).ToConstPool()
		index := c.constPoolUint16Length()
		c.Interfaces = append(c.Interfaces, index)
		c.ConstPool = append(c.ConstPool, inter)
	}
	for _, f := range high.Fields {
		field := &FieldInfo{}
		field.AccessFlags = f.AccessFlags
		field.NameIndex = c.insertUtfConst(f.Name)
		if f.ConstantValue != nil {
			field.Attributes = append(field.Attributes, f.ConstantValue.ToAttributeInfo(c))
		}
		field.DescriptorIndex = c.insertUtfConst(f.Descriptor)
		if f.AttributeLucyFieldDescritor != nil {
			field.Attributes = append(field.Attributes, f.AttributeLucyFieldDescritor.ToAttributeInfo(c))
		}
		c.Fields = append(c.Fields, field)
	}
	for _, ms := range high.Methods {
		for _, m := range ms {
			info := &MethodInfo{}
			info.AccessFlags = m.AccessFlags //accessflag
			info.NameIndex = c.insertUtfConst(m.Name)
			info.DescriptorIndex = c.insertUtfConst(m.Descriptor)
			codeinfo := m.Code.ToAttributeInfo(c)
			info.Attributes = append(info.Attributes, codeinfo)
			if m.AttributeLucyInnerStaticMethod != nil {
				info.Attributes = append(info.Attributes, m.AttributeLucyInnerStaticMethod.ToAttributeInfo(c))
			}
			if m.AttributeLucyMethodDescritor != nil {
				info.Attributes = append(info.Attributes, m.AttributeLucyMethodDescritor.ToAttributeInfo(c))
			}
			if m.AttributeLucyPackageEntranceMethod != nil {
				info.Attributes = append(info.Attributes, m.AttributeLucyPackageEntranceMethod.ToAttributeInfo(c))
			}
			c.Methods = append(c.Methods, info)
		}
	}
	//source file
	c.Attributes = append(c.Attributes, (&AttributeSourceFile{high.getSourceFile()}).ToAttributeInfo(c))
	if c.AttributeClosureClass != nil {
		c.Attributes = append(c.Attributes, c.AttributeClosureClass.ToAttributeInfo(c))
	}
	for _, v := range c.TypeAlias {
		c.Attributes = append(c.Attributes, v.ToAttributeInfo(c))
	}
	c.ifConstPoolOverMaxSize()
	return
}

func (c *Class) constPoolUint16Length() uint16 {
	return uint16(len(c.ConstPool))
}
func (c *Class) ifConstPoolOverMaxSize() {
	if len(c.ConstPool) > CONSTANT_POOL_MAX_SIZE {
		panic(fmt.Sprintf("const pool max size is:%d", CONSTANT_POOL_MAX_SIZE))
	}
}
