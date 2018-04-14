package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type StackMapState struct {
	Locals []*cg.StackMap_verification_type_info
	Stacks []*cg.StackMap_verification_type_info
}

func (s *StackMapState) popStack(pop int) {
	if pop <= 0 {
		panic("negative pop")
	}
	s.Stacks = s.Stacks[:len(s.Stacks)-pop]
}

func (s *StackMapState) FromLast(last *StackMapState) *StackMapState {
	s.Locals = make([]*cg.StackMap_verification_type_info, len(last.Locals))
	copy(s.Locals, last.Locals)
	s.Stacks = make([]*cg.StackMap_verification_type_info, len(last.Stacks))
	copy(s.Stacks, last.Stacks)
	return s
}

func (s *StackMapState) newStackMapVerificationTypeInfo(class *cg.ClassHighLevel, t *ast.VariableType) (ret []*cg.StackMap_verification_type_info) {
	ret = []*cg.StackMap_verification_type_info{}
	switch t.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Integer_variable_info{}
	case ast.VARIABLE_TYPE_LONG:
		ret = make([]*cg.StackMap_verification_type_info, 2)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[1] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Long_variable_info{}
		ret[1].T = &cg.StackMap_Top_variable_info{}
	case ast.VARIABLE_TYPE_FLOAT:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Float_variable_info{}
	case ast.VARIABLE_TYPE_DOUBLE:
		ret = make([]*cg.StackMap_verification_type_info, 2)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[1] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Double_variable_info{}
		ret[1].T = &cg.StackMap_Top_variable_info{}
	case ast.VARIABLE_TYPE_NULL:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Null_variable_info{}
	case ast.VARIABLE_TYPE_STRING:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_string_class),
		}
	case ast.VARIABLE_TYPE_OBJECT:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(t.Class.Name),
		}
	case ast.VARIABLE_TYPE_MAP:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_hashmap_class),
		}
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[t.ArrayType.Typ]
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(meta.classname),
		}
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		d := Descriptor.typeDescriptor(t)
		if d == "" {
			panic(11)
		}
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(d),
		}
	default:
		panic(11)
	}
	return ret
}
