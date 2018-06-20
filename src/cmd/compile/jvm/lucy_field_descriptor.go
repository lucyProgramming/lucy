package jvm

import (
	"bytes"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyFieldSignature struct {
}

func (signature *LucyFieldSignature) Need(variableType *ast.Type) bool {
	return variableType.Type == ast.VARIABLE_TYPE_MAP ||
		variableType.Type == ast.VARIABLE_TYPE_ARRAY ||
		variableType.Type == ast.VARIABLE_TYPE_ENUM
}
func (signature *LucyFieldSignature) Encode(variableType *ast.Type) (d string) {
	if variableType.Type == ast.VARIABLE_TYPE_MAP {
		d = "M" // start token of map
		d += signature.Encode(variableType.Map.K)
		d += signature.Encode(variableType.Map.V)
		return d
	}
	if variableType.Type == ast.VARIABLE_TYPE_ENUM {
		d = "E"
		d += variableType.Enum.Name + ";"
		return d
	}
	if variableType.Type == ast.VARIABLE_TYPE_ARRAY {
		d = "]"
		d += signature.Encode(variableType.ArrayType)
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
		m.Type = ast.VARIABLE_TYPE_MAP
		m.Map = &ast.Map{}
		m.Map.K = kt
		m.Map.V = vt
		return bs, m, nil
	}
	if bs[0] == 'E' {
		bs = bs[1:]
		a := &ast.Type{}
		a.Type = ast.VARIABLE_TYPE_ENUM
		index := bytes.Index(bs, []byte{';'})
		a.Enum = &ast.Enum{}
		a.Enum.Name = string(bs[:index])
		bs = bs[index+1:]
		return bs, a, nil
	}
	if bs[0] == ']' {
		bs = bs[1:]
		a := &ast.Type{}
		a.Type = ast.VARIABLE_TYPE_ARRAY
		bs, a.ArrayType, err = signature.Decode(bs)
		return bs, a, err
	}
	return Descriptor.ParseType(bs)
}
