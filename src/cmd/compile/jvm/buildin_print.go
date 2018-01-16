package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

/*
	function print
*/
func (m *MakeClass) mkBuildinPrint(class *cg.ClassHighLevel, call *ast.ExpressionFunctionCall, code cg.AttributeCode) {
	if len(call.Args) == 0 {
		return
	}
	argsize := 0
	for _, v := range call.Args {
		if v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL {
			panic(111)
		} else {
			argsize++
		}
	}
	if argsize > 1 {

	} else {

	}
}
