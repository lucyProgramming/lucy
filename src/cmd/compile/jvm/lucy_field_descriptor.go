package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

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
