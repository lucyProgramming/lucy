package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StackMapState struct {
	Locals             []*cg.StackMap_verification_type_info
	LastStackMapLocals []*cg.StackMap_verification_type_info
	Stacks             []*cg.StackMap_verification_type_info
}

func (s *StackMapState) appendLocals(class *cg.ClassHighLevel, v *ast.VariableType) {
	s.Locals = append(s.Locals,
		s.newStackMapVerificationTypeInfo(class, v))
}

func (s *StackMapState) addTop(absent *StackMapState) {
	length := len(absent.Locals) - len(s.Locals)
	oldLength := len(s.Locals)
	t := &cg.StackMap_verification_type_info{}
	t.Verify = &cg.StackMap_Top_variable_info{}
	for i := 0; i < length; i++ {
		tt := absent.Locals[i+oldLength].Verify
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
	if len(s.Stacks) == 0 {
		panic("already 0")
	}
	//fmt.Println(0, s.Stacks, len(s.Stacks), len(s.Stacks)-pop, s.Stacks)
	s.Stacks = s.Stacks[:len(s.Stacks)-pop]
}
func (s *StackMapState) pushStack(class *cg.ClassHighLevel, v *ast.VariableType) {
	s.Stacks = append(s.Stacks, s.newStackMapVerificationTypeInfo(class, v))
}
func (s *StackMapState) FromLast(last *StackMapState) *StackMapState {
	s.Locals = make([]*cg.StackMap_verification_type_info, len(last.Locals))
	copy(s.Locals, last.Locals)
	return s
}

func (s *StackMapState) newStackMapVerificationTypeInfo(class *cg.ClassHighLevel,
	t *ast.VariableType) (ret *cg.StackMap_verification_type_info) {
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
		ret.Verify = &cg.StackMap_Integer_variable_info{}
	case ast.VARIABLE_TYPE_LONG:
		ret.Verify = &cg.StackMap_Long_variable_info{}
	case ast.VARIABLE_TYPE_FLOAT:
		ret.Verify = &cg.StackMap_Float_variable_info{}
	case ast.VARIABLE_TYPE_DOUBLE:
		ret.Verify = &cg.StackMap_Double_variable_info{}
	case ast.VARIABLE_TYPE_NULL:
		ret.Verify = &cg.StackMap_Null_variable_info{}
	case ast.VARIABLE_TYPE_STRING:
		ret.Verify = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_string_class),
		}
	case ast.VARIABLE_TYPE_OBJECT:
		ret.Verify = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(t.Class.Name),
		}
	case ast.VARIABLE_TYPE_MAP:
		ret.Verify = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_hashmap_class),
		}
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[t.ArrayType.Typ]
		ret.Verify = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(meta.classname),
		}
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		d := Descriptor.typeDescriptor(t)
		ret.Verify = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(d),
		}
	default:
		panic(11)
	}
	return ret
}
