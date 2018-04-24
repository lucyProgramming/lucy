package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StackMapState struct {
	Locals     []*cg.StackMap_verification_type_info
	LastLocals []*cg.StackMap_verification_type_info
	Stacks     []*cg.StackMap_verification_type_info
}

func (s *StackMapState) appendLocals(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableType) (varOffset uint16) {
	varOffset = code.MaxLocals
	s.Locals = append(s.Locals,
		s.newStackMapVerificationTypeInfo(class, v))
	code.MaxLocals += jvmSize(v)
	return
}

func (s *StackMapState) addTop(absent *StackMapState) {
	length := len(absent.Locals) - len(s.Locals)
	oldLength := len(s.Locals)
	t := &cg.StackMap_verification_type_info{}
	t.PayLoad = &cg.StackMap_Top_variable_info{}
	for i := 0; i < length; i++ {
		tt := absent.Locals[i+oldLength].PayLoad
		_, ok1 := tt.(*cg.StackMap_Double_variable_info)
		_, ok2 := tt.(*cg.StackMap_Long_variable_info)
		if ok1 || ok2 {
			s.Locals = append(s.Locals, t, t)
		} else {
			s.Locals = append(s.Locals, t)
		}

	}
}

func (s *StackMapState) newObjectVariableType(name string) *ast.VariableType {
	ret := &ast.VariableType{}
	ret.Typ = ast.VARIABLE_TYPE_OBJECT
	ret.Class = &ast.Class{}
	ret.Class.Name = name
	return ret
}

func (s *StackMapState) popStack(pop int) {
	if pop < 0 {
		panic("negative pop")
	}
	s.Stacks = s.Stacks[:len(s.Stacks)-pop]
}

//func (s *StackMapState) sliceOutLocals(pop int) {
//	if pop <= 0 {
//		panic("negative pop")
//	}
//	s.Locals = s.Locals[:len(s.Locals)-pop]
//}

func (s *StackMapState) FromLast(last *StackMapState) *StackMapState {
	s.Locals = make([]*cg.StackMap_verification_type_info, len(last.Locals))
	copy(s.Locals, last.Locals)
	s.Stacks = make([]*cg.StackMap_verification_type_info, len(last.Stacks))
	copy(s.Stacks, last.Stacks)
	return s
}

func (s *StackMapState) newStackMapVerificationTypeInfo(class *cg.ClassHighLevel, t *ast.VariableType) (ret *cg.StackMap_verification_type_info) {
	ret = &cg.StackMap_verification_type_info{}
	switch t.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		ret.PayLoad = &cg.StackMap_Integer_variable_info{}
	case ast.VARIABLE_TYPE_LONG:
		ret.PayLoad = &cg.StackMap_Long_variable_info{}
	case ast.VARIABLE_TYPE_FLOAT:
		ret.PayLoad = &cg.StackMap_Float_variable_info{}
	case ast.VARIABLE_TYPE_DOUBLE:
		ret.PayLoad = &cg.StackMap_Double_variable_info{}
	case ast.VARIABLE_TYPE_NULL:
		ret.PayLoad = &cg.StackMap_Null_variable_info{}
	case ast.VARIABLE_TYPE_STRING:
		ret.PayLoad = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_string_class),
		}
	case ast.VARIABLE_TYPE_OBJECT:
		ret.PayLoad = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(t.Class.Name),
		}
	case ast.VARIABLE_TYPE_MAP:
		ret.PayLoad = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_hashmap_class),
		}
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[t.ArrayType.Typ]
		ret.PayLoad = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(meta.classname),
		}
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		d := Descriptor.typeDescriptor(t)
		ret.PayLoad = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(d),
		}
	}
	return ret
}
