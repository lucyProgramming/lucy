package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildQuestion(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	question := e.Data.(*ast.ExpressionQuestion)
	maxStack = buildExpression.build(class, code, question.Selection, context, state)
	falseExit := (&cg.Exit{}).Init(cg.OP_ifeq, code)
	//true part
	stack := buildExpression.build(class, code, question.True, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	trueExit := (&cg.Exit{}).Init(cg.OP_goto, code)
	context.MakeStackMap(code, state, code.CodeLength)
	writeExits([]*cg.Exit{falseExit}, code.CodeLength)
	stack = buildExpression.build(class, code, question.False, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	state.pushStack(class, e.Value)
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	writeExits([]*cg.Exit{trueExit}, code.CodeLength)
	return
}
