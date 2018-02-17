package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
)

type Descriptor struct {
}

func (m *Descriptor) methodDescriptor(f *ast.Function) string {
	if f.Name == ast.MAIN_FUNCTION_NAME {
		return "([Ljava/lang/String;)V"
	}
	s := "("
	for _, v := range f.Typ.ParameterList {
		s += m.typeDescriptor(v.Typ)
	}
	s += ")"
	if len(f.Typ.ReturnList) == 0 {
		s += "V"
	} else if len(f.Typ.ReturnList) == 1 {
		s += m.typeDescriptor(f.Typ.ReturnList[0].Typ)
	} else {
		s += "Ljava/util/ArrayList;" //always this type
	}
	return s
}

func (m *Descriptor) typeDescriptor(v *ast.VariableType) string {
	switch v.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		return "Z"
	case ast.VARIABLE_TYPE_BYTE:
		return "B"
	case ast.VARIABLE_TYPE_SHORT:
		return "S"
	case ast.VARIABLE_TYPE_INT:
		return "I"
	case ast.VARIABLE_TYPE_LONG:
		return "J"
	case ast.VARIABLE_TYPE_FLOAT:
		return "F"
	case ast.VARIABLE_TYPE_DOUBLE:
		return "D"
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		meta := ArrayMetas[v.CombinationType.Typ] // combination type
		return "L" + meta.classname + ";"
	case ast.VARIABLE_TYPE_STRING:
		return "Ljava/lang/String;"
	case ast.VARIABLE_TYPE_VOID:
		return "V"
	case ast.VARIABLE_TYPE_OBJECT:
		return "L" + v.Class.Name + ";"
	}
	panic("unhandle type signature")
}
