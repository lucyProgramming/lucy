package jvm

import (
	"bytes"

	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyFieldSignature struct {
}

func (signature *LucyFieldSignature) typeOf(variableType *ast.Type) string {
	s := ""
	signature.typeOf_(variableType, &s)
	fmt.Println()
	return s
}

func (signature *LucyFieldSignature) typeOf_(variableType *ast.Type, s *string) {
	switch variableType.Type {
	case ast.VariableTypeBool:
		*s += "bool"
	case ast.VariableTypeByte:
		*s += "byte"
	case ast.VariableTypeShort:
		*s += "short"
	case ast.VariableTypeInt:
		*s += "int"
	case ast.VariableTypeLong:
		*s += "long"
	case ast.VariableTypeFloat:
		*s += "float"
	case ast.VariableTypeDouble:
		*s += "double"
	case ast.VariableTypeString:
		*s += "string"
	case ast.VariableTypeObject:
		*s += variableType.Class.Name
	case ast.VariableTypeMap:
		*s += "map{"
		*s = signature.typeOf(variableType.Map.K)
		*s += "->"
		*s += signature.typeOf(variableType.Map.V)
		*s += "}"
	case ast.VariableTypeArray:
		*s += "[]"
		*s += signature.typeOf(variableType.Array)
	case ast.VariableTypeJavaArray:
		*s += signature.typeOf(variableType.Array)
		if variableType.IsVArgs {
			*s += "..."
		} else {
			*s += "[]"
		}
	case ast.VariableTypeFunction:
		*s += "fn("
		for k, v := range variableType.FunctionType.ParameterList {
			*s += signature.typeOf(v.Type)
			if k != len(variableType.FunctionType.ParameterList)-1 {
				*s += ","
			}
		}
		if variableType.FunctionType.VArgs != nil {
			if len(variableType.FunctionType.ParameterList) > 0 {
				*s += ","
			}
			*s += signature.typeOf(variableType.FunctionType.VArgs.Type)
		}
		*s += ")"
		if len(variableType.FunctionType.ReturnList) > 0 {
			*s += "->("
			for k, v := range variableType.FunctionType.ReturnList {
				*s += signature.typeOf(v.Type)
				if k != len(variableType.FunctionType.ParameterList)-1 {
					*s += ","
				}
			}
			*s += ")"
		}
	case ast.VariableTypeEnum:
		*s += variableType.Enum.Name

	}

}

func (signature *LucyFieldSignature) Need(variableType *ast.Type) bool {
	return variableType.Type == ast.VariableTypeMap ||
		variableType.Type == ast.VariableTypeArray ||
		variableType.Type == ast.VariableTypeEnum ||
		variableType.Type == ast.VariableTypeFunction
}

func (signature *LucyFieldSignature) Encode(variableType *ast.Type) (d string) {
	if variableType.Type == ast.VariableTypeMap {
		d = "M" // start token of map
		d += signature.Encode(variableType.Map.K)
		d += signature.Encode(variableType.Map.V)
		return d
	}
	if variableType.Type == ast.VariableTypeEnum {
		d = "E"
		d += variableType.Enum.Name + ";"
		return d
	}
	if variableType.Type == ast.VariableTypeArray {
		d = "]"
		d += signature.Encode(variableType.Array)
		return d
	}
	if variableType.Type == ast.VariableTypeFunction {
		d = LucyMethodSignatureParser.Encode(variableType.FunctionType)
		return d
	}
	return Descriptor.typeDescriptor(variableType)
}

func (signature *LucyFieldSignature) Decode(bs []byte) ([]byte, *ast.Type, error) {
	var err error
	if bs[0] == 'M' {
		bs = bs[1:]
		var kt *ast.Type
		bs, kt, err = signature.Decode(bs)
		if err != nil {
			return bs, nil, err
		}
		var vt *ast.Type
		bs, vt, err = signature.Decode(bs)
		if err != nil {
			return bs, nil, err
		}
		m := &ast.Type{}
		m.Type = ast.VariableTypeMap
		m.Map = &ast.Map{}
		m.Map.K = kt
		m.Map.V = vt
		return bs, m, nil
	}
	if bs[0] == 'E' {
		bs = bs[1:]
		a := &ast.Type{}
		a.Type = ast.VariableTypeEnum
		index := bytes.Index(bs, []byte{';'})
		a.Enum = &ast.Enum{}
		a.Enum.Name = string(bs[:index])
		bs = bs[index+1:]
		return bs, a, nil
	}
	if bs[0] == '(' {
		a := &ast.Type{}
		a.Type = ast.VariableTypeFunction
		a.FunctionType = &ast.FunctionType{}
		bs, err = LucyMethodSignatureParser.Decode(a.FunctionType, bs)
		if err != nil {
			return bs, nil, err
		}
		return bs, a, nil
	}
	if bs[0] == ']' {
		bs = bs[1:]
		a := &ast.Type{}
		a.Type = ast.VariableTypeArray
		bs, a.Array, err = signature.Decode(bs)
		return bs, a, err
	}
	return Descriptor.ParseType(bs)
}
