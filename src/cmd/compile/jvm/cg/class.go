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
	dest              io.Writer
	magic             uint32 //0xCAFEBABE
	minorVersion      uint16
	majorVersion      uint16
	constPool         []*ConstPool
	accessFlag        uint16
	thisClass         uint16
	superClass        uint16
	interfaces        []uint16
	fields            []*FieldInfo
	methods           []*MethodInfo
	attributes        []*AttributeInfo
	Utf8Consts        map[string]*ConstPool
	IntConsts         map[int32]*ConstPool
	LongConsts        map[int64]*ConstPool
	FloatConsts       map[float32]*ConstPool
	DoubleConsts      map[float64]*ConstPool
	ClassConsts       map[string]*ConstPool
	StringConsts      map[string]*ConstPool
	FieldRefConsts    map[CONSTANT_Fieldref_info_high_level]*ConstPool
	NameAndTypeConsts map[CONSTANT_NameAndType_info_high_level]*ConstPool
}

func (c *Class) InsertNameAndConst(n CONSTANT_NameAndType_info_high_level) uint16 {
	if c.NameAndTypeConsts == nil {
		c.NameAndTypeConsts = make(map[CONSTANT_NameAndType_info_high_level]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.NameAndTypeConsts[n]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_NameAndType_info{
		name:       c.insertUtfConst(n.Name),
		descriptor: c.insertUtfConst(n.Type),
	}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.NameAndTypeConsts[n] = info
	return info.selfindex
}
func (c *Class) InsertFieldRefConst(f CONSTANT_Fieldref_info_high_level) uint16 {
	if c.FieldRefConsts == nil {
		c.FieldRefConsts = make(map[CONSTANT_Fieldref_info_high_level]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.FieldRefConsts[f]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Fieldref_info{
		classIndex:       c.InsertClassConst(f.Class),
		nameAndTypeIndex: c.InsertNameAndConst(CONSTANT_NameAndType_info_high_level{f.Name, f.Descriptor}),
	}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.FieldRefConsts[f] = info
	return info.selfindex
}
func (c *Class) insertUtfConst(s string) uint16 {
	if c.Utf8Consts == nil {
		c.Utf8Consts = make(map[string]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.Utf8Consts[s]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Utf8_info{uint16(len(s)), []byte(s)}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.Utf8Consts[s] = info
	return info.selfindex
}

func (c *Class) InsertIntConst(i int32) uint16 {
	if c.IntConsts == nil {
		c.IntConsts = make(map[int32]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.IntConsts[i]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Integer_info{i}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.IntConsts[i] = info
	return info.selfindex
}
func (c *Class) InsertLongConst(i int64) uint16 {
	if c.LongConsts == nil {
		c.LongConsts = make(map[int64]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.LongConsts[i]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Long_info{i}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.LongConsts[i] = info
	return info.selfindex
}

func (c *Class) InsertFloatConst(f float32) uint16 {
	if c.FloatConsts == nil {
		c.FloatConsts = make(map[float32]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.FloatConsts[f]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Float_info{f}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.FloatConsts[f] = info
	return info.selfindex
}

func (c *Class) InsertDoubleConst(f float64) uint16 {
	if c.DoubleConsts == nil {
		c.DoubleConsts = make(map[float64]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.DoubleConsts[f]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Double_info{f}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.DoubleConsts[f] = info
	return info.selfindex
}

func (c *Class) InsertClassConst(name string) uint16 {
	if c.ClassConsts == nil {
		c.ClassConsts = make(map[string]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.ClassConsts[name]; ok {
		return con.selfindex
	}
	info := (&CONSTANT_Class_info{c.insertUtfConst(name)}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.ClassConsts[name] = info
	return info.selfindex
}

func (c *Class) InsertStringConst(s string) uint16 {
	if c.StringConsts == nil {
		c.StringConsts = make(map[string]*ConstPool)
	}
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil}
	}
	if con, ok := c.StringConsts[s]; ok {
		return con.selfindex
	}

	info := (&CONSTANT_String_info{c.insertUtfConst(s)}).ToConstPool()
	info.selfindex = c.constPoolUint16Length()
	c.constPool = append(c.constPool, info)
	c.StringConsts[s] = info
	return info.selfindex
}

func (high *ClassHighLevel) FromHighLevel() *Class {
	high.Class.fromHighLevel(high)
	return &high.Class
}

func (c *Class) fromHighLevel(high *ClassHighLevel) {
	c.minorVersion = 0
	c.majorVersion = 49
	if len(c.constPool) == 0 {
		c.constPool = []*ConstPool{nil} // jvm const pool index begin at 1
	}
	//fieldref
	//	for f, locations := range high.FieldRefs {
	//		info := (&CONSTANT_Fieldref_info{}).ToConstPool()
	//		high.InsertClassConst(f.Class, info.info[0:2])
	//		index := c.constPoolUint16Length()
	//		backPatchIndex(locations, index)
	//		c.constPool = append(c.constPool, info)
	//		high.InsertNameAndTypeConst(CONSTANT_NameAndType_info_high_level{
	//			Name: f.Name,
	//			Type: f.Descriptor,
	//		}, info.info[2:4])
	//	}
	//methodref
	for m, locations := range high.MethodRefs {
		info := (&CONSTANT_Methodref_info{}).ToConstPool()
		high.InsertClassConst(m.Class, info.info[0:2])
		index := c.constPoolUint16Length()
		backPatchIndex(locations, index)
		c.constPool = append(c.constPool, info)
		binary.BigEndian.PutUint16(info.info[2:4], c.InsertNameAndConst(CONSTANT_NameAndType_info_high_level{
			Name: m.Name,
			Type: m.Descriptor,
		}))

	}
	//name and type
	//	for nt, locations := range high.NameAndTypes {
	//		info := (&CONSTANT_NameAndType_info{c.insertUtfConst(nt.Name), c.insertUtfConst(nt.Type)}).ToConstPool()
	//		index := c.constPoolUint16Length()
	//		backPatchIndex(locations, index)
	//		c.constPool = append(c.constPool, info)

	//	}
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
