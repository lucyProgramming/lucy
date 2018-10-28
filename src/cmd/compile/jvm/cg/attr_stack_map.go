package cg

import (
	"encoding/binary"
	"fmt"
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
	info.NameIndex = class.InsertUtf8Const(AttributeNameStackMap)
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(len(a.StackMaps)))
	for _, v := range a.StackMaps {
		bs = append(bs, v.ToBytes()...)
	}
	info.Info = bs
	info.attributeLength = uint32(len(info.Info))
	return info
}

type StackMapSameFrame struct {
	FrameType byte
}

func (s *StackMapSameFrame) ToBytes() []byte {
	return []byte{s.FrameType}
}

type StackMapSameLocals1StackItemFrame struct {
	FrameType byte
	Stack     *StackMapVerificationTypeInfo
}

func (s *StackMapSameLocals1StackItemFrame) ToBytes() []byte {
	bs := []byte{s.FrameType}
	bs = append(bs, s.Stack.ToBytes()...)
	return bs
}

type StackMapSameLocals1StackItemFrameExtended struct {
	FrameType byte
	Delta     uint16
	Stack     *StackMapVerificationTypeInfo
}

func (s *StackMapSameLocals1StackItemFrameExtended) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	return append(bs, s.Stack.ToBytes()...)
}

type StackMapChopFrame struct {
	FrameType byte
	Delta     uint16
}

func (s *StackMapChopFrame) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	return bs
}

type StackMapSameFrameExtended struct {
	FrameType byte
	Delta     uint16
}

func (s *StackMapSameFrameExtended) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	return bs
}

type StackMapAppendFrame struct {
	FrameType byte
	Delta     uint16
	Locals    []*StackMapVerificationTypeInfo
}

func (s *StackMapAppendFrame) ToBytes() []byte {
	bs := make([]byte, 3)
	bs[0] = s.FrameType
	binary.BigEndian.PutUint16(bs[1:], s.Delta)
	for _, v := range s.Locals {
		bs = append(bs, v.ToBytes()...)
	}
	return bs
}

type StackMapFullFrame struct {
	FrameType byte
	Delta     uint16
	Locals    []*StackMapVerificationTypeInfo
	Stacks    []*StackMapVerificationTypeInfo
}

func (s *StackMapFullFrame) ToBytes() []byte {
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

type StackMapTopVariableInfo struct{}
type StackMapIntegerVariableInfo struct{}
type StackMapFloatVariableInfo struct{}
type StackMapLongVariableInfo struct{}
type StackMapDoubleVariableInfo struct{}
type StackMapNullVariableInfo struct{}
type StackMapUninitializedThisVariableInfo struct{}
type StackMapObjectVariableInfo struct {
	Index    uint16
	Readable string
}
type StackMapUninitializedVariableInfo struct {
	CodeOffset uint16
}

type StackMapVerificationTypeInfo struct {
	Verify interface{}
}

func (s *StackMapVerificationTypeInfo) ToString() string {
	switch s.Verify.(type) {
	case *StackMapTopVariableInfo:
		return "top"
	case *StackMapIntegerVariableInfo:
		return "int"
	case *StackMapFloatVariableInfo:
		return "float"
	case *StackMapDoubleVariableInfo:
		return "double"
	case *StackMapLongVariableInfo:
		return "long"
	case *StackMapNullVariableInfo:
		return "null"
	case *StackMapUninitializedThisVariableInfo:
		return "unInitThis"
	case *StackMapObjectVariableInfo:
		return s.Verify.(*StackMapObjectVariableInfo).Readable
	case *StackMapUninitializedVariableInfo:
		return fmt.Sprintf("unInitVariable@%d", s.Verify.(*StackMapUninitializedVariableInfo).CodeOffset)
	}
	return ""
}

func (s *StackMapVerificationTypeInfo) ToBytes() []byte {
	switch s.Verify.(type) {
	case *StackMapTopVariableInfo:
		return []byte{0}
	case *StackMapIntegerVariableInfo:
		return []byte{1}
	case *StackMapFloatVariableInfo:
		return []byte{2}
	case *StackMapDoubleVariableInfo:
		return []byte{3}
	case *StackMapLongVariableInfo:
		return []byte{4}
	case *StackMapNullVariableInfo:
		return []byte{5}
	case *StackMapUninitializedThisVariableInfo:
		return []byte{6}
	case *StackMapObjectVariableInfo:
		bs := make([]byte, 3)
		bs[0] = 7
		binary.BigEndian.PutUint16(bs[1:], s.Verify.(*StackMapObjectVariableInfo).Index)
		return bs
	case *StackMapUninitializedVariableInfo:
		bs := make([]byte, 3)
		bs[0] = 8
		binary.BigEndian.PutUint16(bs[1:], s.Verify.(*StackMapUninitializedVariableInfo).CodeOffset)
		return bs
	default:
	}
	return nil
}
