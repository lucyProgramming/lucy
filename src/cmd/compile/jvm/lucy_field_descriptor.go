package jvm

import (
	"bytes"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyFieldSignatureParse struct {
}

func (l *LucyFieldSignatureParse) Need(variableType *ast.VariableType) bool {
	return variableType.Typ == ast.VARIABLE_TYPE_MAP ||
		variableType.Typ == ast.VARIABLE_TYPE_ARRAY
}
func (l *LucyFieldSignatureParse) Encode(variableType *ast.VariableType) (d string) {
	if variableType.Typ == ast.VARIABLE_TYPE_MAP {
		d = "m" // start token of map
		d += l.Encode(variableType.Map.K)
		d += l.Encode(variableType.Map.V)
		return d
	}
	if variableType.Typ == ast.VARIABLE_TYPE_ARRAY {
		d = "]"
		d += l.Encode(variableType.ArrayType)
		return d
	}
	if variableType.Typ == ast.VARIABLE_TYPE_T {
		return fmt.Sprintf("T%s;", variableType.Name)
	}
	return Descriptor.typeDescriptor(variableType)
}
func (l *LucyFieldSignatureParse) Decode(bs []byte) ([]byte, *ast.VariableType, error) {
	var err error
	if bs[0] == 'm' {
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
	if bs[0] == ']' {
		bs = bs[1:]
		a := &ast.VariableType{}
		a.Typ = ast.VARIABLE_TYPE_ARRAY
		bs, a.ArrayType, err = l.Decode(bs)
		return bs, a, err
	}
	if bs[0] == 'T' {
		bs = bs[1:]
		a := &ast.VariableType{}
		a.Typ = ast.VARIABLE_TYPE_T
		index := bytes.Index(bs, []byte{';'})
		a.Name = string(bs[0:index])
		bs = bs[index+1:]
		fmt.Println(a.Name)
		return bs, a, nil
	}
	return Descriptor.ParseType(bs)
}
