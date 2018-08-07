package cg

import (
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
	writer                 io.Writer
	magic                  uint32 //0xCAFEBABE
	MinorVersion           uint16
	MajorVersion           uint16
	ConstPool              []*ConstPool
	AccessFlag             uint16
	ThisClass              uint16
	SuperClass             uint16
	Interfaces             []uint16
	Fields                 []*FieldInfo
	Methods                []*MethodInfo
	Attributes             []*AttributeInfo
	AttributeCompilerAuto  *AttributeCompilerAuto
	AttributeGroupedByName AttributeGroupedByName
	TypeAlias              []*AttributeLucyTypeAlias
	AttributeLucyEnum      *AttributeLucyEnum
	//const caches
	Utf8Constants               map[string]*ConstPool
	IntConstants                map[int32]*ConstPool
	LongConstants               map[int64]*ConstPool
	FloatConstants              map[float32]*ConstPool
	DoubleConstants             map[float64]*ConstPool
	ClassConstants              map[string]*ConstPool
	StringConstants             map[string]*ConstPool
	FieldRefConstants           map[CONSTANT_Fieldref_info_high_level]*ConstPool
	NameAndTypeConstants        map[CONSTANT_NameAndType_info_high_level]*ConstPool
	MethodrefConstants          map[CONSTANT_Methodref_info_high_level]*ConstPool
	InterfaceMethodrefConstants map[CONSTANT_InterfaceMethodref_info_high_level]*ConstPool
	MethodTypeConstants         map[CONSTANT_MethodType_info_high_level]*ConstPool
}

func (c *Class) InsertMethodTypeConst(n CONSTANT_MethodType_info_high_level) uint16 {
	if c.MethodTypeConstants == nil {
		c.MethodTypeConstants = make(map[CONSTANT_MethodType_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.MethodTypeConstants[n]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_MethodType_info{
		descriptorIndex: c.InsertUtf8Const(n.Descriptor),
	}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.MethodTypeConstants[n] = info
	return info.selfIndex
}

func (c *Class) InsertInterfaceMethodrefConst(n CONSTANT_InterfaceMethodref_info_high_level) uint16 {
	if c.InterfaceMethodrefConstants == nil {
		c.InterfaceMethodrefConstants = make(map[CONSTANT_InterfaceMethodref_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.InterfaceMethodrefConstants[n]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_InterfaceMethodref_info{
		classIndex: c.InsertClassConst(n.Class),
		nameAndTypeIndex: c.InsertNameAndType(CONSTANT_NameAndType_info_high_level{
			Name:       n.Method,
			Descriptor: n.Descriptor,
		}),
	}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.InterfaceMethodrefConstants[n] = info
	return info.selfIndex
}

func (c *Class) InsertMethodrefConst(n CONSTANT_Methodref_info_high_level) uint16 {
	if c.MethodrefConstants == nil {
		c.MethodrefConstants = make(map[CONSTANT_Methodref_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.MethodrefConstants[n]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Methodref_info{
		classIndex: c.InsertClassConst(n.Class),
		nameAndTypeIndex: c.InsertNameAndType(CONSTANT_NameAndType_info_high_level{
			Name:       n.Method,
			Descriptor: n.Descriptor,
		}),
	}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.MethodrefConstants[n] = info
	return info.selfIndex
}

func (c *Class) InsertNameAndType(n CONSTANT_NameAndType_info_high_level) uint16 {
	if c.NameAndTypeConstants == nil {
		c.NameAndTypeConstants = make(map[CONSTANT_NameAndType_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.NameAndTypeConstants[n]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_NameAndType_info{
		name:       c.InsertUtf8Const(n.Name),
		descriptor: c.InsertUtf8Const(n.Descriptor),
	}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.NameAndTypeConstants[n] = info
	return info.selfIndex
}
func (c *Class) InsertFieldRefConst(f CONSTANT_Fieldref_info_high_level) uint16 {
	if c.FieldRefConstants == nil {
		c.FieldRefConstants = make(map[CONSTANT_Fieldref_info_high_level]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.FieldRefConstants[f]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Fieldref_info{
		classIndex: c.InsertClassConst(f.Class),
		nameAndTypeIndex: c.InsertNameAndType(CONSTANT_NameAndType_info_high_level{
			Name:       f.Field,
			Descriptor: f.Descriptor,
		}),
	}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.FieldRefConstants[f] = info
	return info.selfIndex
}
func (c *Class) InsertUtf8Const(s string) uint16 {
	if c.Utf8Constants == nil {
		c.Utf8Constants = make(map[string]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.Utf8Constants[s]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Utf8_info{uint16(len(s)), []byte(s)}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.Utf8Constants[s] = info
	return info.selfIndex
}

func (c *Class) InsertIntConst(i int32) uint16 {
	if c.IntConstants == nil {
		c.IntConstants = make(map[int32]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.IntConstants[i]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Integer_info{i}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.IntConstants[i] = info
	return info.selfIndex
}
func (c *Class) InsertLongConst(i int64) uint16 {
	if c.LongConstants == nil {
		c.LongConstants = make(map[int64]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.LongConstants[i]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Long_info{i}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info, nil)
	c.LongConstants[i] = info
	return info.selfIndex
}

func (c *Class) InsertFloatConst(f float32) uint16 {
	if c.FloatConstants == nil {
		c.FloatConstants = make(map[float32]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.FloatConstants[f]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Float_info{f}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.FloatConstants[f] = info
	return info.selfIndex
}

func (c *Class) InsertDoubleConst(f float64) uint16 {
	if c.DoubleConstants == nil {
		c.DoubleConstants = make(map[float64]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.DoubleConstants[f]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Double_info{f}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info, nil)
	c.DoubleConstants[f] = info
	return info.selfIndex
}

func (c *Class) InsertClassConst(name string) uint16 {
	if c.ClassConstants == nil {
		c.ClassConstants = make(map[string]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.ClassConstants[name]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_Class_info{c.InsertUtf8Const(name)}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.ClassConstants[name] = info
	return info.selfIndex
}

func (c *Class) InsertStringConst(s string) uint16 {
	if c.StringConstants == nil {
		c.StringConstants = make(map[string]*ConstPool)
	}
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil}
	}
	if con, ok := c.StringConstants[s]; ok {
		return con.selfIndex
	}
	info := (&CONSTANT_String_info{c.InsertUtf8Const(s)}).ToConstPool()
	info.selfIndex = c.constPoolUint16Length()
	c.ConstPool = append(c.ConstPool, info)
	c.StringConstants[s] = info
	return info.selfIndex
}

func (classHighLevel *ClassHighLevel) ToLow(jvmVersion int) *Class {
	classHighLevel.Class.fromHighLevel(classHighLevel, jvmVersion)
	return &classHighLevel.Class
}

func (c *Class) fromHighLevel(high *ClassHighLevel, jvmVersion int) {
	c.MinorVersion = 0
	c.MajorVersion = uint16(jvmVersion)
	if len(c.ConstPool) == 0 {
		c.ConstPool = []*ConstPool{nil} // jvm const pool index begin at 1
	}
	c.AccessFlag = high.AccessFlags
	c.ThisClass = c.InsertClassConst(high.Name)
	c.SuperClass = c.InsertClassConst(high.SuperClass)
	for _, i := range high.Interfaces {
		inter := (&CONSTANT_Class_info{c.InsertUtf8Const(i)}).ToConstPool()
		index := c.constPoolUint16Length()
		c.Interfaces = append(c.Interfaces, index)
		c.ConstPool = append(c.ConstPool, inter)
	}
	for _, f := range high.Fields {
		field := &FieldInfo{}
		field.AccessFlags = f.AccessFlags
		field.NameIndex = c.InsertUtf8Const(f.Name)
		if f.AttributeConstantValue != nil {
			field.Attributes = append(field.Attributes, f.AttributeConstantValue.ToAttributeInfo(c))
		}
		field.DescriptorIndex = c.InsertUtf8Const(f.Descriptor)
		if f.AttributeLucyFieldDescriptor != nil {
			field.Attributes = append(field.Attributes, f.AttributeLucyFieldDescriptor.ToAttributeInfo(c))
		}
		if f.AttributeLucyConst != nil {
			field.Attributes = append(field.Attributes, f.AttributeLucyConst.ToAttributeInfo(c))
		}
		c.Fields = append(c.Fields, field)
	}
	for _, ms := range high.Methods {
		for _, m := range ms {
			info := &MethodInfo{}
			info.AccessFlags = m.AccessFlags
			info.NameIndex = c.InsertUtf8Const(m.Name)
			info.DescriptorIndex = c.InsertUtf8Const(m.Descriptor)
			if m.Code != nil {
				info.Attributes = append(info.Attributes, m.Code.ToAttributeInfo(c))
			}

			if m.AttributeLucyMethodDescriptor != nil {
				info.Attributes = append(info.Attributes, m.AttributeLucyMethodDescriptor.ToAttributeInfo(c))
			}
			if m.AttributeLucyTriggerPackageInitMethod != nil {
				info.Attributes = append(info.Attributes, m.AttributeLucyTriggerPackageInitMethod.ToAttributeInfo(c))
			}
			if m.AttributeDefaultParameters != nil {
				info.Attributes = append(info.Attributes, m.AttributeDefaultParameters.ToAttributeInfo(c))
			}
			if m.AttributeMethodParameters != nil {
				t := m.AttributeMethodParameters.ToAttributeInfo(c)
				if t != nil {
					info.Attributes = append(info.Attributes, t)
				}
			}
			if m.AttributeCompilerAuto != nil {
				info.Attributes = append(info.Attributes, m.AttributeCompilerAuto.ToAttributeInfo(c))
			}
			if m.AttributeLucyReturnListNames != nil {
				t := m.AttributeLucyReturnListNames.ToAttributeInfo(c, AttributeNameLucyReturnListNames)
				if t != nil {
					info.Attributes = append(info.Attributes, t)
				}
			}
			c.Methods = append(c.Methods, info)
		}
	}
	//source file
	c.Attributes = append(c.Attributes, (&AttributeSourceFile{high.getSourceFile()}).ToAttributeInfo(c))
	if c.AttributeCompilerAuto != nil {
		c.Attributes = append(c.Attributes, c.AttributeCompilerAuto.ToAttributeInfo(c))
	}
	for _, v := range c.TypeAlias {
		c.Attributes = append(c.Attributes, v.ToAttributeInfo(c))
	}
	if c.AttributeLucyEnum != nil {
		c.Attributes = append(c.Attributes, c.AttributeLucyEnum.ToAttributeInfo(c))
	}
	for _, v := range high.TemplateFunctions {
		c.Attributes = append(c.Attributes, v.ToAttributeInfo(c))
	}
	c.ifConstPoolOverMaxSize()
	return
}

func (c *Class) constPoolUint16Length() uint16 {
	return uint16(len(c.ConstPool))
}
func (c *Class) ifConstPoolOverMaxSize() {
	if len(c.ConstPool) > ConstantPoolMaxSize {
		panic(fmt.Sprintf("const pool max size is:%d", ConstantPoolMaxSize))
	}
}
