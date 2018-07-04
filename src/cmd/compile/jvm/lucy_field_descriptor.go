package jvm

import (
	"bytes"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyFieldSignature struct {
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
		d += signature.Encode(variableType.Map.Key)
		d += signature.Encode(variableType.Map.Value)
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
		m.Map.Key = kt
		m.Map.Value = vt
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
