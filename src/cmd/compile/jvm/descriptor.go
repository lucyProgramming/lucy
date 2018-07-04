package jvm

import (
	"bytes"
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type Description struct {
}

func (description *Description) methodDescriptor(f *ast.FunctionType) string {
	s := "("
	for _, v := range f.ParameterList {
		s += description.typeDescriptor(v.Type)
	}
	s += ")"
	if f.NoReturnValue() {
		s += "V"
	} else if len(f.ReturnList) == 1 {
		s += description.typeDescriptor(f.ReturnList[0].Type)
	} else {
		s += "[Ljava/lang/Object;" //always this type
	}
	return s
}

func (description *Description) typeDescriptor(v *ast.Type) string {
	switch v.Type {
	case ast.VariableTypeBool:
		return "Z"
	case ast.VariableTypeByte:
		return "B"
	case ast.VariableTypeShort:
		return "S"
	case ast.VariableTypeInt, ast.VariableTypeEnum:
		return "I"
	case ast.VariableTypeLong:
		return "J"
	case ast.VariableTypeFloat:
		return "F"
	case ast.VariableTypeDouble:
		return "D"
	case ast.VariableTypeArray:
		meta := ArrayMetas[v.Array.Type] // combination type
		return "L" + meta.className + ";"
	case ast.VariableTypeString:
		return "Ljava/lang/String;"
	case ast.VariableTypeVoid:
		return "V"
	case ast.VariableTypeObject:
		return "L" + v.Class.Name + ";"
	case ast.VariableTypeMap:
		return "L" + javaMapClass + ";"
	case ast.VariableTypeFunction:
		return "L" + javaMethodHandleClass + ";"
	case ast.VariableTypeJavaArray:
		return "[" + description.typeDescriptor(v.Array)
	}
	panic("unHandle type signature")
}

func (description *Description) ParseType(bs []byte) ([]byte, *ast.Type, error) {
	switch bs[0] {
	case 'V':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeVoid,
		}, nil
	case 'B':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeByte,
		}, nil
	case 'D':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeDouble,
		}, nil
	case 'F':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeFloat,
		}, nil
	case 'I':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeInt,
		}, nil
	case 'J':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeLong,
		}, nil
	case 'S', 'C':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeShort,
		}, nil
	case 'Z':
		bs = bs[1:]
		return bs, &ast.Type{
			Type: ast.VariableTypeBool,
		}, nil
	case 'L':
		bs = bs[1:]
		index := bytes.Index(bs, []byte{';'}) // end token
		t := &ast.Type{}
		t.Type = ast.VariableTypeObject
		t.Class = &ast.Class{}
		t.Class.Name = string(bs[:index])
		bs = bs[index+1:] // skip ;
		t.Class.NotImportedYet = true
		if t.Class.Name == javaStringClass {
			t.Type = ast.VariableTypeString
		}
		return bs, t, nil
	case '[':
		bs = bs[1:]
		var t *ast.Type
		var err error
		bs, t, err = description.ParseType(bs)
		ret := &ast.Type{}
		if err == nil {
			ret.Type = ast.VariableTypeJavaArray
			ret.Array = t
		}
		return bs, ret, err
	}
	return bs, nil, fmt.Errorf("unkown type:%v", string(bs))
}

func (description *Description) ParseFunctionType(bs []byte) (ast.FunctionType, error) {
	t := ast.FunctionType{}
	if bs[0] != '(' {
		return t, fmt.Errorf("function descriptor does not start with '('")
	}
	bs = bs[1:]
	i := 1
	var err error
	for bs[0] != ')' {
		vd := &ast.Variable{}
		vd.Name = fmt.Sprintf("var_%d", i)
		bs, vd.Type, err = description.ParseType(bs)
		if err != nil {
			return t, err
		}
		t.ParameterList = append(t.ParameterList, vd)
		i++
	}
	bs = bs[1:] // skip )
	vd := &ast.Variable{}
	vd.Name = "returnValue"
	_, vd.Type, err = description.ParseType(bs)
	if err != nil {
		return t, err
	}
	t.ReturnList = append(t.ReturnList, vd)
	return t, nil
}
