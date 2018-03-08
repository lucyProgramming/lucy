package jvm

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"

	"bytes"
	"strings"
)

type Descript struct {
}

func (m *Descript) methodDescriptor(f *ast.Function) string {
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

func (m *Descript) typeDescriptor(v *ast.VariableType) string {
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
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[v.ArrayType.Typ] // combination type
		return "L" + meta.classname + ";"
	case ast.VARIABLE_TYPE_STRING:
		return "Ljava/lang/String;"
	case ast.VARIABLE_TYPE_VOID:
		return "V"
	case ast.VARIABLE_TYPE_OBJECT:
		return "L" + v.Class.ClassNameDefinition.Name + ";"
	}
	panic("unhandle type signature")
}

func (m *Descript) ParseType(descritpor string) (*ast.VariableType, error) {
	switch descritpor {
	case "V":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_VOID,
		}, nil
	case "B":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BYTE,
		}, nil
	case "D":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_DOUBLE,
		}, nil
	case "F":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_FLOAT,
		}, nil
	case "I":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_INT,
		}, nil
	case "J":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_LONG,
		}, nil
	case "S", "C":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_SHORT,
		}, nil
	case "Z":
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_SHORT,
		}, nil
	default:
		if strings.HasPrefix(descritpor, "L") {
			return &ast.VariableType{
				Typ:  ast.VARIABLE_TYPE_OBJECT,
				Name: descritpor[1:],
			}, nil
		} else if strings.HasPrefix(descritpor, "[") {
			t, err := m.ParseType(descritpor[1:])
			if err != nil {
				return nil, err
			}
			return &ast.VariableType{
				Typ:       ast.VARIABLE_TYPE_ARRAY,
				ArrayType: t,
			}, nil
		}
	}
	return nil, fmt.Errorf("unkown type:%v", descritpor)
}

func (m *Descript) ParseFunctionType(bs []byte) (*ast.FunctionType, error) {
	t := &ast.FunctionType{}
	if bs[0] != '(' {
		return nil, fmt.Errorf("function descriptor does not start with '('")
	}
	bs = bs[1:]
	i := 1
	for bs[0] != ')' {
		switch bs[0] {
		case 'B':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_BYTE
			t.ParameterList = append(t.ParameterList, vd)
		case 'D':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_DOUBLE
			t.ParameterList = append(t.ParameterList, vd)
		case 'F':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_FLOAT
			t.ParameterList = append(t.ParameterList, vd)
		case 'I':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_INT
			t.ParameterList = append(t.ParameterList, vd)
		case 'J':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_LONG
			t.ParameterList = append(t.ParameterList, vd)
		case 'S', 'C':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_SHORT
			t.ParameterList = append(t.ParameterList, vd)
		case 'Z':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_BOOL
			t.ParameterList = append(t.ParameterList, vd)
		case 'L':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_OBJECT

			bytes.Index(bs, []byte{';'})
		case '[':
			panic(22)
		}
		i++
	}
	bs = bs[1:] // skip )
	vd := &ast.VariableDefinition{}
	vd.Name = "return"
	var err error
	vd.Typ, err = m.ParseType(string(bs))
	if err != nil {
		return t, err
	}
	t.ReturnList = append(t.ReturnList, vd)
	return t, nil
}
