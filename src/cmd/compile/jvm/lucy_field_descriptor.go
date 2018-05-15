package jvm

import (
	"bytes"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyFieldSignatureParse struct {
}

func (l *LucyFieldSignatureParse) Need(variableType *ast.VariableType) bool {
	return variableType.Typ == ast.VARIABLE_TYPE_MAP ||
		variableType.Typ == ast.VARIABLE_TYPE_ARRAY ||
		variableType.Typ == ast.VARIABLE_TYPE_ENUM
}
func (l *LucyFieldSignatureParse) Encode(variableType *ast.VariableType) (d string) {
	if variableType.Typ == ast.VARIABLE_TYPE_MAP {
		d = "M" // start token of map
		d += l.Encode(variableType.Map.K)
		d += l.Encode(variableType.Map.V)
		return d
	}
	if variableType.Typ == ast.VARIABLE_TYPE_ENUM {
		d = "E"
		d += variableType.Enum.Name + ";"
		return d
	}
	if variableType.Typ == ast.VARIABLE_TYPE_ARRAY {
		d = "]"
		d += l.Encode(variableType.ArrayType)
		return d
	}
	return Descriptor.typeDescriptor(variableType)
}
func (l *LucyFieldSignatureParse) Decode(bs []byte) ([]byte, *ast.VariableType, error) {
	var err error
	if bs[0] == 'M' {
		bs = bs[1:]
		var kt *ast.VariableType
		bs, kt, err = l.Decode(bs)
		if err != nil {
			return bs, nil, err
		}
		var vt *ast.VariableType
		bs, vt, err = l.Decode(bs)
		if err != nil {
			return bs, nil, err
		}
		m := &ast.VariableType{}
		m.Typ = ast.VARIABLE_TYPE_MAP
		m.Map = &ast.Map{}
		m.Map.K = kt
		m.Map.V = vt
		return bs, m, nil
	}
	if bs[0] == 'E' {
		bs = bs[1:]
		a := &ast.VariableType{}
		a.Typ = ast.VARIABLE_TYPE_ENUM
		index := bytes.Index(bs, []byte{';'})
		a.Enum = &ast.Enum{}
		a.Enum.Name = string(bs[:index])
		bs = bs[index+1:]
		return bs, a, nil
	}
	if bs[0] == ']' {
		bs = bs[1:]
		a := &ast.VariableType{}
		a.Typ = ast.VARIABLE_TYPE_ARRAY
		bs, a.ArrayType, err = l.Decode(bs)
		return bs, a, err
	}
	return Descriptor.ParseType(bs)
}
