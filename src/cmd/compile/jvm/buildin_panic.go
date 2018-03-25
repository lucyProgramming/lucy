package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinPanic(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context) (maxstack uint16) {
	maxstack, _ = m.build(class, code, call.Args[0], context)
	code.Codes[code.CodeLength] = cg.OP_athrow
	code.CodeLength++
	return
}
