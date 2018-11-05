package cg

import (
	"encoding/binary"
)

type AttributeInnerClasses struct {
	Classes []*InnerClass
}

func (this *AttributeInnerClasses) FromBs(class *Class, bs []byte) {
	length := binary.BigEndian.Uint16(bs)
	bs = bs[2:]
	if int(length*8) != len(bs) {
		panic("length not match")
	}
	this.Classes = nil
	for len(bs) > 0 {
		inner := &InnerClass{}
		inner.FromBs(class, bs[:8])
		this.Classes = append(this.Classes, inner)
		bs = bs[8:]
	}
}

func (this *AttributeInnerClasses) ToAttributeInfo(class *Class) *AttributeInfo {
	if this == nil || len(this.Classes) == 0 {
		return nil
	}
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameInnerClasses)
	ret.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.Info, uint16(len(this.Classes)))
	for _, v := range this.Classes {
		bs8 := make([]byte, 8)
		binary.BigEndian.PutUint16(bs8, class.InsertClassConst(v.InnerClass))
		binary.BigEndian.PutUint16(bs8[2:], class.InsertClassConst(v.OuterClass))
		binary.BigEndian.PutUint16(bs8[4:], class.InsertUtf8Const(v.Name))
		binary.BigEndian.PutUint16(bs8[6:], v.AccessFlags)
		ret.Info = append(ret.Info, bs8...)
	}
	ret.attributeLength = uint32(len(ret.Info))
	return ret
}

type InnerClass struct {
	InnerClass  string
	OuterClass  string
	Name        string
	AccessFlags uint16
}

func (inner *InnerClass) FromBs(class *Class, bs []byte) {
	nameIndex := binary.BigEndian.Uint16(class.ConstPool[binary.BigEndian.Uint16(bs)].Info)
	inner.InnerClass = string(class.ConstPool[nameIndex].Info)
	if 0 == binary.BigEndian.Uint16(bs[2:]) {
		//TODO:: what zero means???
	} else {
		nameIndex = binary.BigEndian.Uint16(class.ConstPool[binary.BigEndian.Uint16(bs[2:])].Info)
		inner.OuterClass = string(class.ConstPool[nameIndex].Info)
	}
	if 0 == binary.BigEndian.Uint16(bs[4:]) {
		//TODO:: what zero means???
	} else {
		inner.Name = string(class.ConstPool[binary.BigEndian.Uint16(bs[4:])].Info)
	}
	inner.AccessFlags = binary.BigEndian.Uint16(bs[6:])
}
