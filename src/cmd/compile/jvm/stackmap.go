package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StackMapState struct {
	Locals []*cg.StackMapVerificationTypeInfo
	Stacks []*cg.StackMapVerificationTypeInfo
}

// same as last
//func (s *StackMapState) isSame(locals []*cg.StackMap_verification_type_info, stacks []*cg.StackMap_verification_type_info) bool {
//	if len(s.Locals) != len(locals) || len(s.Stacks) != len(stacks) {
//		return false
//	}
//	for k, v := range s.Locals {
//		if v.Equal(locals[k]) == false {
//			return false
//		}
//	}
//	for k, v := range s.Stacks {
//		if v.Equal(stacks[k]) == false {
//			return false
//		}
//	}
//	return true
//}

func (s *StackMapState) appendLocals(class *cg.ClassHighLevel, v *ast.Type) {
	s.Locals = append(s.Locals,
		s.newStackMapVerificationTypeInfo(class, v))
}

func (s *StackMapState) addTop(absent *StackMapState) {
	if s == absent {
		return
	}
	length := len(absent.Locals) - len(s.Locals)
	oldLength := len(s.Locals)
	t := &cg.StackMapVerificationTypeInfo{}
	t.Verify = &cg.StackMapTopVariableInfo{}
	for i := 0; i < length; i++ {
		tt := absent.Locals[i+oldLength].Verify
		_, ok1 := tt.(*cg.StackMapDoubleVariableInfo)
		_, ok2 := tt.(*cg.StackMapLongVariableInfo)
		if ok1 || ok2 {
			s.Locals = append(s.Locals, t, t)
		} else {
			s.Locals = append(s.Locals, t)
		}
	}
}

func (s *StackMapState) newObjectVariableType(name string) *ast.Type {
	ret := &ast.Type{}
	ret.Type = ast.VARIABLE_TYPE_OBJECT
	ret.Class = &ast.Class{}
	ret.Class.Name = name
	return ret
}

func (s *StackMapState) popStack(pop int) {
	if pop < 0 {
		panic("negative pop")
	}
	if pop == 0 {
		return
	}
	if len(s.Stacks) == 0 {
		panic("already 0")
	}
	s.Stacks = s.Stacks[:len(s.Stacks)-pop]
}
func (s *StackMapState) pushStack(class *cg.ClassHighLevel, v *ast.Type) {
	if s == nil {
		panic("s is nil")
	}
	s.Stacks = append(s.Stacks, s.newStackMapVerificationTypeInfo(class, v))
}
func (s *StackMapState) FromLast(last *StackMapState) *StackMapState {
	s.Locals = make([]*cg.StackMapVerificationTypeInfo, len(last.Locals))
	copy(s.Locals, last.Locals)
	return s
}

func (s *StackMapState) newStackMapVerificationTypeInfo(class *cg.ClassHighLevel,
	t *ast.Type) (ret *cg.StackMapVerificationTypeInfo) {
	ret = &cg.StackMapVerificationTypeInfo{}
	switch t.Type {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		ret.Verify = &cg.StackMapIntegerVariableInfo{}
	case ast.VARIABLE_TYPE_LONG:
		ret.Verify = &cg.StackMapLongVariableInfo{}
	case ast.VARIABLE_TYPE_FLOAT:
		ret.Verify = &cg.StackMapFloatVariableInfo{}
	case ast.VARIABLE_TYPE_DOUBLE:
		ret.Verify = &cg.StackMapDoubleVariableInfo{}
	case ast.VARIABLE_TYPE_NULL:
		ret.Verify = &cg.StackMapNullVariableInfo{}
	case ast.VARIABLE_TYPE_STRING:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(java_string_class),
		}
	case ast.VARIABLE_TYPE_OBJECT:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(t.Class.Name),
		}
	case ast.VARIABLE_TYPE_MAP:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(java_hashmap_class),
		}
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[t.ArrayType.Type]
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(meta.className),
		}
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		d := Descriptor.typeDescriptor(t)
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(d),
		}
	}
	return ret
}
