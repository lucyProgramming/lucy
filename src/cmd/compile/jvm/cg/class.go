package cg

import (
	"encoding/binary"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"io"
)

const (
	AccClassPublic     uint16 = 0x0001 // 可以被包的类外访问。
	AccClassFinal      uint16 = 0x0010 //不允许有子类。
	AccClassSuper      uint16 = 0x0020 //当用到invokespecial指令时，需要特殊处理③的父类方法。
	AccClassInterface  uint16 = 0x0200 // 标识定义的是接口而不是类。
	AccClassAbstract   uint16 = 0x0400 //  不能被实例化。
	AccClassSynthetic  uint16 = 0x1000 //标识并非Java源码生成的代码。
	AccClassAnnotation uint16 = 0x2000 // 标识注解类型
	AccClassEnum       uint16 = 0x4000 // 标识枚举类型
)

type Class struct {
	writer                  io.Writer
	magic                   uint32 //0xCAFEBABE
	MinorVersion            uint16
	MajorVersion            uint16
	ConstPool               []*ConstPool
	AccessFlag              uint16
	ThisClass               uint16
	SuperClass              uint16
	Interfaces              []uint16
	Fields                  []*FieldInfo
	Methods                 []*MethodInfo
	Attributes              []*AttributeInfo
	AttributeGroupedByName  AttributeGroupedByName
	TypeAlias               []*AttributeLucyTypeAlias
	AttributeLucyEnum       *AttributeLucyEnum
	AttributeLucyComment    *AttributeLucyComment
	AttributeLucyClassConst *AttributeLucyClassConst
	AttributeInnerClasses   AttributeInnerClasses
	//const caches
	Utf8Constants               map[string]*ConstPool
	IntConstants                map[int32]*ConstPool
	LongConstants               map[int64]*ConstPool
	FloatConstants              map[float32]*ConstPool
	DoubleConstants             map[float64]*ConstPool
	ClassConstants              map[string]*ConstPool
	StringConstants             map[string]*ConstPool
	FieldRefConstants           map[ConstantInfoFieldrefHighLevel]*ConstPool
	NameAndTypeConstants        map[ConstantInfoNameAndTypeHighLevel]*ConstPool
	MethodrefConstants          map[ConstantInfoMethodrefHighLevel]*ConstPool
	InterfaceMethodrefConstants map[ConstantInfoInterfaceMethodrefHighLevel]*ConstPool
	MethodTypeConstants         map[ConstantInfoMethodTypeHighLevel]*ConstPool
}

func (this *Class) IsSynthetic() bool {
	return (this.AccessFlag & AccClassSynthetic) != 0
}

func (this *Class) InsertMethodTypeConst(n ConstantInfoMethodTypeHighLevel) uint16 {
	if this.MethodTypeConstants == nil {
		this.MethodTypeConstants = make(map[ConstantInfoMethodTypeHighLevel]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.MethodTypeConstants[n]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoMethodType{
		descriptorIndex: this.InsertUtf8Const(n.Descriptor),
	}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.MethodTypeConstants[n] = info
	return info.selfIndex
}

func (this *Class) InsertInterfaceMethodrefConst(n ConstantInfoInterfaceMethodrefHighLevel) uint16 {
	if this.InterfaceMethodrefConstants == nil {
		this.InterfaceMethodrefConstants = make(map[ConstantInfoInterfaceMethodrefHighLevel]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.InterfaceMethodrefConstants[n]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoInterfaceMethodref{
		classIndex: this.InsertClassConst(n.Class),
		nameAndTypeIndex: this.InsertNameAndType(ConstantInfoNameAndTypeHighLevel{
			Name:       n.Method,
			Descriptor: n.Descriptor,
		}),
	}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.InterfaceMethodrefConstants[n] = info
	return info.selfIndex
}

func (this *Class) InsertMethodrefConst(n ConstantInfoMethodrefHighLevel) uint16 {
	if this.MethodrefConstants == nil {
		this.MethodrefConstants = make(map[ConstantInfoMethodrefHighLevel]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.MethodrefConstants[n]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoMethodref{
		classIndex: this.InsertClassConst(n.Class),
		nameAndTypeIndex: this.InsertNameAndType(ConstantInfoNameAndTypeHighLevel{
			Name:       n.Method,
			Descriptor: n.Descriptor,
		}),
	}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.MethodrefConstants[n] = info
	return info.selfIndex
}

func (this *Class) InsertNameAndType(n ConstantInfoNameAndTypeHighLevel) uint16 {
	if this.NameAndTypeConstants == nil {
		this.NameAndTypeConstants = make(map[ConstantInfoNameAndTypeHighLevel]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.NameAndTypeConstants[n]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoNameAndType{
		name:       this.InsertUtf8Const(n.Name),
		descriptor: this.InsertUtf8Const(n.Descriptor),
	}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.NameAndTypeConstants[n] = info
	return info.selfIndex
}
func (this *Class) InsertFieldRefConst(f ConstantInfoFieldrefHighLevel) uint16 {
	if this.FieldRefConstants == nil {
		this.FieldRefConstants = make(map[ConstantInfoFieldrefHighLevel]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.FieldRefConstants[f]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoFieldref{
		classIndex: this.InsertClassConst(f.Class),
		nameAndTypeIndex: this.InsertNameAndType(ConstantInfoNameAndTypeHighLevel{
			Name:       f.Field,
			Descriptor: f.Descriptor,
		}),
	}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.FieldRefConstants[f] = info
	return info.selfIndex
}
func (this *Class) InsertUtf8Const(s string) uint16 {
	if this.Utf8Constants == nil {
		this.Utf8Constants = make(map[string]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.Utf8Constants[s]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoUtf8{uint16(len(s)), []byte(s)}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.Utf8Constants[s] = info
	return info.selfIndex
}

func (this *Class) InsertIntConst(i int32) uint16 {
	if this.IntConstants == nil {
		this.IntConstants = make(map[int32]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.IntConstants[i]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoInteger{i}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.IntConstants[i] = info
	return info.selfIndex
}
func (this *Class) InsertLongConst(i int64) uint16 {
	if this.LongConstants == nil {
		this.LongConstants = make(map[int64]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.LongConstants[i]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoLong{i}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info, nil)
	this.LongConstants[i] = info
	return info.selfIndex
}

func (this *Class) InsertFloatConst(f float32) uint16 {
	if this.FloatConstants == nil {
		this.FloatConstants = make(map[float32]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.FloatConstants[f]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoFloat{f}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.FloatConstants[f] = info
	return info.selfIndex
}

func (this *Class) InsertDoubleConst(f float64) uint16 {
	if this.DoubleConstants == nil {
		this.DoubleConstants = make(map[float64]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.DoubleConstants[f]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoDouble{f}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info, nil)
	this.DoubleConstants[f] = info
	return info.selfIndex
}

func (this *Class) InsertClassConst(name string) uint16 {
	if this.ClassConstants == nil {
		this.ClassConstants = make(map[string]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.ClassConstants[name]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoClass{this.InsertUtf8Const(name)}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.ClassConstants[name] = info
	return info.selfIndex
}

func (this *Class) InsertStringConst(s string) uint16 {
	if this.StringConstants == nil {
		this.StringConstants = make(map[string]*ConstPool)
	}
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil}
	}
	if con, ok := this.StringConstants[s]; ok {
		return con.selfIndex
	}
	info := (&ConstantInfoString{this.InsertUtf8Const(s)}).ToConstPool()
	info.selfIndex = this.constPoolUint16Length()
	this.ConstPool = append(this.ConstPool, info)
	this.StringConstants[s] = info
	return info.selfIndex
}

func (this *ClassHighLevel) ToLow() *Class {
	this.Class.fromHighLevel(this)
	return &this.Class
}

func (this *Class) fromHighLevel(high *ClassHighLevel) {
	this.MajorVersion = uint16(common.CompileFlags.JvmMajorVersion)
	this.MinorVersion = uint16(common.CompileFlags.JvmMinorVersion)
	if len(this.ConstPool) == 0 {
		this.ConstPool = []*ConstPool{nil} // jvm const pool index begin at 1
	}
	this.AccessFlag = high.AccessFlags
	this.ThisClass = this.InsertClassConst(high.Name)
	this.SuperClass = this.InsertClassConst(high.SuperClass)
	for _, i := range high.Interfaces {
		inter := (&ConstantInfoClass{this.InsertUtf8Const(i)}).ToConstPool()
		index := this.constPoolUint16Length()
		this.Interfaces = append(this.Interfaces, index)
		this.ConstPool = append(this.ConstPool, inter)
	}
	for _, f := range high.Fields {
		field := &FieldInfo{}
		field.AccessFlags = f.AccessFlags
		field.NameIndex = this.InsertUtf8Const(f.Name)
		if f.AttributeConstantValue != nil {
			field.Attributes = append(field.Attributes, f.AttributeConstantValue.ToAttributeInfo(this))
		}
		field.DescriptorIndex = this.InsertUtf8Const(f.Descriptor)
		if f.AttributeLucyFieldDescriptor != nil {
			field.Attributes = append(field.Attributes, f.AttributeLucyFieldDescriptor.ToAttributeInfo(this))
		}
		if f.AttributeLucyConst != nil {
			field.Attributes = append(field.Attributes, f.AttributeLucyConst.ToAttributeInfo(this))
		}
		if f.AttributeLucyComment != nil {
			field.Attributes = append(field.Attributes, f.AttributeLucyComment.ToAttributeInfo(this))
		}
		this.Fields = append(this.Fields, field)
	}
	for _, ms := range high.Methods {
		for _, m := range ms {
			info := &MethodInfo{}
			info.AccessFlags = m.AccessFlags
			info.NameIndex = this.InsertUtf8Const(m.Name)
			info.DescriptorIndex = this.InsertUtf8Const(m.Descriptor)
			if m.Code != nil {
				info.Attributes = append(info.Attributes, m.Code.ToAttributeInfo(this))
			}

			if m.AttributeLucyMethodDescriptor != nil {
				info.Attributes = append(info.Attributes, m.AttributeLucyMethodDescriptor.ToAttributeInfo(this))
			}
			//if m.AttributeLucyTriggerPackageInitMethod != nil {
			//	info.Attributes = append(info.Attributes, m.AttributeLucyTriggerPackageInitMethod.ToAttributeInfo(this))
			//}
			if m.AttributeDefaultParameters != nil {
				info.Attributes = append(info.Attributes, m.AttributeDefaultParameters.ToAttributeInfo(this))
			}
			if m.AttributeLucyComment != nil {
				info.Attributes = append(info.Attributes, m.AttributeLucyComment.ToAttributeInfo(this))
			}
			if m.AttributeMethodParameters != nil {
				t := m.AttributeMethodParameters.ToAttributeInfo(this)
				if t != nil {
					info.Attributes = append(info.Attributes, t)
				}
			}
			if m.AttributeLucyReturnListNames != nil {
				t := m.AttributeLucyReturnListNames.ToAttributeInfo(this, AttributeNameLucyReturnListNames)
				if t != nil {
					info.Attributes = append(info.Attributes, t)
				}
			}
			this.Methods = append(this.Methods, info)
		}
	}
	//source file
	this.Attributes = append(this.Attributes,
		(&AttributeSourceFile{high.getSourceFile()}).ToAttributeInfo(this))
	for _, v := range this.TypeAlias {
		this.Attributes = append(this.Attributes, v.ToAttributeInfo(this))
	}
	if this.AttributeLucyEnum != nil {
		this.Attributes = append(this.Attributes, this.AttributeLucyEnum.ToAttributeInfo(this))
	}
	if this.AttributeLucyComment != nil {
		this.Attributes = append(this.Attributes, this.AttributeLucyComment.ToAttributeInfo(this))
	}
	if this.AttributeLucyClassConst != nil {
		this.Attributes = append(this.Attributes, this.AttributeLucyClassConst.ToAttributeInfo(this))
	}
	if a := this.AttributeInnerClasses.ToAttributeInfo(this); a != nil {
		this.Attributes = append(this.Attributes, a)
	}
	for _, v := range high.TemplateFunctions {
		this.Attributes = append(this.Attributes, v.ToAttributeInfo(this))
	}
	this.ifConstPoolOverMaxSize()
	return
}

func (this *Class) constPoolUint16Length() uint16 {
	return uint16(len(this.ConstPool))
}
func (this *Class) ifConstPoolOverMaxSize() {
	if len(this.ConstPool) > ConstantPoolMaxSize {
		panic(fmt.Sprintf("const pool max size is:%d", ConstantPoolMaxSize))
	}
}

func (this *Class) IsInnerClass() (is bool) {
	if len(this.AttributeGroupedByName[AttributeNameInnerClasses]) == 0 {
		return
	}
	innerClass := this.AttributeGroupedByName[AttributeNameInnerClasses][0]
	var attr AttributeInnerClasses
	thisClassName := string(this.ConstPool[binary.BigEndian.Uint16(this.ConstPool[this.ThisClass].Info)].Info)
	attr.FromBs(this, innerClass.Info)
	for _, v := range attr.Classes {
		if thisClassName == v.InnerClass {
			return true
		}
	}
	return
}
