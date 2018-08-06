package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	maxStack = buildExpression.build(class, code, bin.Left, context, state)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 2 > maxStack { // dup increment stack
		maxStack = 2
	}
	var exit *cg.Exit
	if e.Type == ast.ExpressionTypeLogicalOr {
		exit = (&cg.Exit{}).Init(cg.OP_ifne, code)
	} else {
		exit = (&cg.Exit{}).Init(cg.OP_ifeq, code)
	}
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	stack := buildExpression.build(class, code, bin.Right, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	state.pushStack(class, e.Value)
	writeExits([]*cg.Exit{exit}, code.CodeLength)
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	return
}
