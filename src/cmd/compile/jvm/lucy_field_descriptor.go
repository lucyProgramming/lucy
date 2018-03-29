package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type LucyFieldSignatureParse struct {
}

func (l *LucyFieldSignatureParse) Need(variableType *ast.VariableType) bool {
	// map need descriptor
	if variableType.Typ == ast.VARIABLE_TYPE_MAP {
		return true
	}
	if variableType.Typ != ast.VARIABLE_TYPE_ARRAY {
		return false
	}
	if variableType.ArrayType.Typ == ast.VARIABLE_TYPE_BOOL ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_BYTE ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_SHORT ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_INT ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_LONG ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_FLOAT ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_DOUBLE ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_STRING {
		return false
	}
	return true
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
	}
	return Descriptor.typeDescriptor(variableType)
}
func (l *LucyFieldSignatureParse) Decode(bs []byte) ([]byte, *ast.VariableType, error) {
	var err error
	if bs[0] == 'm' {
		bs = bs[1:]
		var kt *ast.VariableType
		bs, kt, err = Descriptor.ParseType(bs)
		if err != nil {
			return bs, nil, err
		}
		var vt *ast.VariableType
		bs, vt, err = Descriptor.ParseType(bs)
		if err != nil {
			return bs, nil, err
		}
		m := &ast.VariableType{}
		m.Map = &ast.Map{}
		m.Map.K = kt
		m.Map.V = vt
		return bs, m, nil
	}
	if bs[0] == ']' {
		bs = bs[1:]
		a := &ast.VariableType{}
		a.Typ = ast.VARIABLE_TYPE_ARRAY
		bs, a.ArrayType, err = Descriptor.ParseType(bs)
		return bs, a, err
	}
	return Descriptor.ParseType(bs)
}
