package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyMethodSignatureParse struct {
}

func (LucyMethodSignatureParse) Parse(bs []byte) (*ast.FunctionType, error) {
	return nil, nil
	//if bs[0] != '(' {
	//	return nil, fmt.Errorf("signature dose not beging with '('")
	//}
	//bs = bs[1:]
	//ret := &ast.FunctionType{}
	//i := 1
	//for bs[0] != ')' {
	//	switch bs[0] {
	//	case 'B':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_BYTE
	//		bs = bs[1:]
	//	case 'C':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_SHORT
	//		bs = bs[1:]
	//	case 'D':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_DOUBLE
	//		bs = bs[1:]
	//	case 'F':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_FLOAT
	//		bs = bs[1:]
	//	case 'I':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_INT
	//		bs = bs[1:]
	//	case 'J':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_LONG
	//		bs = bs[1:]
	//	case 'S':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_SHORT
	//		bs = bs[1:]
	//	case 'Z':
	//		vd := &ast.VariableDefinition{}
	//		vd.Name = fmt.Sprintf("var_%d", i)
	//		vd.Typ = &ast.VariableType{}
	//		vd.Typ.Typ = ast.VARIABLE_TYPE_BOOL
	//		bs = bs[1:]
	//	case '[':
	//	case 'L':
	//	default:
	//	}
	//	i++
	//}
	//bs = bs[1:] // skip )

}

func (LucyMethodSignatureParse) Need(functionType *ast.FunctionType) bool {
	for _, v := range functionType.ParameterList {
		if LucyFieldSignatureParser.Need(v.Typ) {
			return true
		}
	}
	for _, v := range functionType.ReturnList {
		if LucyFieldSignatureParser.Need(v.Typ) {
			return true
		}
	}
	return false
}
