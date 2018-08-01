package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) mkBuildInTypeOf(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	//TODO:: should eval first or not
	maxStack = buildExpression.build(class, code, call.Args[0], context, state)
	if jvmSlotSize(call.Args[0].Value) == 2 {
		code.Codes[code.CodeLength] = cg.OP_pop2
		code.CodeLength++
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst(LucyFieldSignatureParser.typeOf(call.Args[0].Value), code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if 1 > maxStack {
		maxStack = 1
	}
	if e.IsStatementExpression {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}
	return
}
