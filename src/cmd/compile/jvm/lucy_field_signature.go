package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type LucyFieldSignatureParse struct {
}

func (l *LucyFieldSignatureParse) Need(variableType *ast.VariableType) (n bool) {
	if variableType.Typ == ast.VARIABLE_TYPE_MAP {
		return true
	}
	if variableType.Typ != ast.VARIABLE_TYPE_ARRAY {
		return
	}
	if variableType.ArrayType.Typ == ast.VARIABLE_TYPE_BOOL ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_BYTE ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_SHORT ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_INT ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_LONG ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_FLOAT ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_DOUBLE ||
		variableType.ArrayType.Typ == ast.VARIABLE_TYPE_STRING {
		return
	}

	return true
}
