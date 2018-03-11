package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinRecover(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	//if e.IsStatementExpression { // first level statement
	//	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, b.InheritedAttribute.Function.AutoVarForException.Offset)...)
	//	return
	//}
	//maxstack = 2
	////load to stack
	//copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, b.InheritedAttribute.Function.AutoVarForException.Offset)...) // load
	////set 2 null
	//code.Codes[code.CodeLength] = cg.OP_aconst_null
	//code.CodeLength++
	//copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, b.InheritedAttribute.Function.AutoVarForException.Offset)...) // load
	return
}
