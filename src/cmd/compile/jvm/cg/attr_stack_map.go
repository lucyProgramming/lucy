package cg

import (
	"encoding/binary"
)

type StackMap interface {
	ToBytes() []byte
}
type AttributeStackMap struct {
	StackMaps []StackMap
}

func (a *AttributeStackMap) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil || len(a.StackMaps) == 0 {
		return nil
	}
	info := &AttributeInfo{}
	info.NameIndex = class.insertUtf8Const(ATTRIBUTE_NAME_STACK_MAP)
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.StackMaps)))
	for _, v := range a.StackMaps {
		bs = append(bs, v.ToBytes()...)
	}
	info.Info = bs
	info.attributeLength = uint32(len(info.Info))
	return info
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
	binary.BigEndian.PutUint16(bs2, uint16(len(s.Stacks)))
	bs = append(bs, bs2...)
	for _, v := range s.Stacks {
		bs = append(bs, v.ToBytes()...)
	}
	return bs
}

type StackMap_Top_variable_info struct{}
type StackMap_Integer_variable_info struct{}
type StackMap_Float_variable_info struct{}
type StackMap_Long_variable_info struct{}
type StackMap_Double_variable_info struct{}
type StackMap_Null_variable_info struct{}
type StackMap_UninitializedThis_variable_info struct{}
type StackMap_Object_variable_info struct {
	Index uint16
}
type StackMap_Uninitialized_variable_info struct {
	CodeOffset uint16
}

type StackMap_verification_type_info struct {
	Verify interface{}
}

//
//func (s *StackMap_verification_type_info) Equal(s2 *StackMap_verification_type_info) bool {
//	if s == s2 {
//		return true
//	}
//	if reflect.DeepEqual(s.Verify, s2.Verify) {
//		return true
//	}
//	// same as top
//	if t1, ok := s.Verify.(*StackMap_Top_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Top_variable_info); ok && t2 != nil {
//			return true
//		}
//	}
//	// same as int
//	if t1, ok := s.Verify.(*StackMap_Integer_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Integer_variable_info); ok && t2 != nil {
//			return true
//		}
//	}
//	// same as float
//	if t1, ok := s.Verify.(*StackMap_Float_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Float_variable_info); ok && t2 != nil {
//			return true
//		}
//	}
//	// same as double
//	if t1, ok := s.Verify.(*StackMap_Double_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Double_variable_info); ok && t2 != nil {
//			return true
//		}
//	}
//	// same as long
//	if t1, ok := s.Verify.(*StackMap_Long_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Long_variable_info); ok && t2 != nil {
//			return true
//		}
//	}
//	// same as null
//	if t1, ok := s.Verify.(*StackMap_Null_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Null_variable_info); ok && t2 != nil {
//			return true
//		}
//	}
//	// same as uninitialized this
//	if t1, ok := s.Verify.(*StackMap_UninitializedThis_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_UninitializedThis_variable_info); ok && t2 != nil {
//			return true
//		}
//	}
//	// same as object
//	if t1, ok := s.Verify.(*StackMap_Object_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Object_variable_info); ok && t2 != nil {
//			return t1.Index == t2.Index
//		}
//	}
//	// same as uninitialized variable
//	if t1, ok := s.Verify.(*StackMap_Uninitialized_variable_info); ok && t1 != nil {
//		if t2, ok := s2.Verify.(*StackMap_Uninitialized_variable_info); ok && t2 != nil {
//			return t1.CodeOffset == t2.CodeOffset
//		}
//	}
//	return false
//}

func (s *StackMap_verification_type_info) ToBytes() []byte {
	switch s.Verify.(type) {
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
		binary.BigEndian.PutUint16(bs[1:], s.Verify.(*StackMap_Object_variable_info).Index)
		return bs
	case *StackMap_Uninitialized_variable_info:
		bs := make([]byte, 3)
		bs[0] = 8
		binary.BigEndian.PutUint16(bs[1:], s.Verify.(*StackMap_Uninitialized_variable_info).CodeOffset)
		return bs
	default:
	}
	return nil
}
