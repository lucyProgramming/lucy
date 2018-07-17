package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildQuestion(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	question := e.Data.(*ast.ExpressionQuestion)
	var es []*cg.Exit
	maxStack, es = buildExpression.build(class, code, question.Selection, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, question.Selection.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	falseExit := (&cg.Exit{}).FromCode(cg.OP_ifeq, code)
	//true part
	stack, es := buildExpression.build(class, code, question.True, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, question.True.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if stack > maxStack {
		maxStack = stack
	}
	trueExit := (&cg.Exit{}).FromCode(cg.OP_goto, code)
	context.MakeStackMap(code, state, code.CodeLength)
	writeExits([]*cg.Exit{falseExit}, code.CodeLength)
	stack, es = buildExpression.build(class, code, question.False, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, question.False.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if stack > maxStack {
		maxStack = stack
	}
	state.pushStack(class, e.ExpressionValue)
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	writeExits([]*cg.Exit{trueExit}, code.CodeLength)
	return

}
