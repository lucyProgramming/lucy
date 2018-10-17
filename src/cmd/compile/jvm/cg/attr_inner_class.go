package cg

import (
	"encoding/binary"
)

type AttributeInnerClasses struct {
	Classes []*InnerClass
}

type InnerClass struct {
	InnerClassInfoIndex   uint16
	OuterClassInfoIndex   uint16
	InnerNameIndex        uint16
	InnerClassAccessFlags uint16
}

func (a *AttributeInnerClasses) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil || len(a.Classes) == 0 {
		return nil
	}
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameInnerClasses)
	ret.Info = make([]byte, 2)
	binary.BigEndian.PutUint16(ret.Info, uint16(len(a.Classes)))
	for _, v := range a.Classes {
		bs8 := make([]byte, 8)
		binary.BigEndian.PutUint16(bs8, v.InnerClassInfoIndex)
		binary.BigEndian.PutUint16(bs8[2:], v.OuterClassInfoIndex)
		binary.BigEndian.PutUint16(bs8[4:], v.InnerNameIndex)
		binary.BigEndian.PutUint16(bs8[6:], v.InnerClassAccessFlags)
		ret.Info = append(ret.Info, bs8...)
	}
	ret.attributeLength = uint32(len(ret.Info))
	return ret
}
