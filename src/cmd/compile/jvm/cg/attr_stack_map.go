package cg

import (
	"encoding/binary"
)

type AttributeStackMap struct {
	StackMaps []StackMap
}

func (a *AttributeStackMap) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil || len(a.StackMaps) == 0 {
		return nil
	}
	info := &AttributeInfo{}
	info.NameIndex = class.insertUtfConst(ATTRIBUTE_NAME_STACK_MAP)
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.StackMaps)))
	for _, v := range a.StackMaps {
		bs = append(bs, v.ToBytes()...)
	}
	info.Info = bs
	info.attributeLength = uint32(len(info.Info))
	return info
}

type StackMap interface {
	ToBytes() []byte
}
type StackMap_same_frame struct {
	FrameType byte
}

func (s *StackMap_same_frame) ToBytes() []byte {
	return []byte{s.FrameType}
}

type StackMap_same_locals_1_stack_item_frame struct {
	FrameType byte
	Stack     *StackMap_verification_type_info
}

func (s *StackMap_same_locals_1_stack_item_frame) ToBytes() []byte {
	bs := []byte{s.FrameType}
	bs = append(bs, s.Stack.ToBytes()...)
	return bs
}

type StackMap_same_locals_1_stack_item_frame_extended struct {
	FrameType byte
	Delta     uint16
	Stack     *StackMap_verification_type_info
}

func (s *StackMap_same_locals_1_stack_item_frame_extended) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	return append(bs, s.Stack.ToBytes()...)
}

type StackMap_chop_frame struct {
	FrameType byte
	Delta     uint16
}

func (s *StackMap_chop_frame) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	return bs
}

type StackMap_same_frame_extended struct {
	FrameType byte
	Delta     uint16
}

func (s *StackMap_same_frame_extended) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	return bs
}

type StackMap_append_frame struct {
	FrameType byte
	Delta     uint16
	Locals    []*StackMap_verification_type_info
}

func (s *StackMap_append_frame) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	for _, v := range s.Locals {
		bs = append(bs, v.ToBytes()...)
	}
	return bs
}

type StackMap_full_frame struct {
	FrameType byte
	Delta     uint16
	Locals    []*StackMap_verification_type_info
	Stacks    []*StackMap_verification_type_info
}

func (s *StackMap_full_frame) ToBytes() []byte {
	bs := make([]byte, 5)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	binary.BigEndian.PutUint16(bs[3:], uint16(len(s.Locals)))
	for _, v := range s.Locals {
		bs = append(bs, v.ToBytes()...)
	}
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs[3:], uint16(len(s.Stacks)))
	bs = append(bs, bs2...)
	for _, v := range s.Stacks {
		bs = append(bs, v.ToBytes()...)
	}
	return bs
}

type StackMap_Top_variable_info struct {
}
type StackMap_Integer_variable_info struct {
}
type StackMap_Float_variable_info struct {
}
type StackMap_Long_variable_info struct {
}
type StackMap_Double_variable_info struct {
}
type StackMap_Null_variable_info struct {
}
type StackMap_UninitializedThis_variable_info struct {
}
type StackMap_Object_variable_info struct {
	Index uint16
}
type StackMap_Uninitialized_variable_info struct {
	Index uint16
}

type StackMap_verification_type_info struct {
	T interface{}
}

func (s *StackMap_verification_type_info) ToBytes() []byte {
	switch s.T.(type) {
	case *StackMap_Top_variable_info:
		return []byte{0}
	case *StackMap_Integer_variable_info:

		return []byte{1}
	case *StackMap_Float_variable_info:
		return []byte{2}
	case *StackMap_Double_variable_info:
		return []byte{3}
	case *StackMap_Long_variable_info:
		return []byte{4}
	case *StackMap_Null_variable_info:
		return []byte{5}
	case *StackMap_UninitializedThis_variable_info:
		return []byte{6}
	case *StackMap_Object_variable_info:
		bs := make([]byte, 3)
		bs[0] = 7
		binary.BigEndian.PutUint16(bs[1:], s.T.(*StackMap_Object_variable_info).Index)
		return bs
	case *StackMap_Uninitialized_variable_info:
		bs := make([]byte, 3)
		bs[0] = 8
		binary.BigEndian.PutUint16(bs[1:], s.T.(*StackMap_Object_variable_info).Index)
		return bs
	default:
	}
	return nil
}
