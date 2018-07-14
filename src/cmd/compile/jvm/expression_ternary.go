package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildTernary(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	ternary := e.Data.(*ast.ExpressionTernary)
	var es []*cg.Exit
	maxStack, es = buildExpression.build(class, code, ternary.Selection, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, ternary.Selection.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	falseExit := (&cg.Exit{}).FromCode(cg.OP_ifeq, code)
	//true part
	stack, es := buildExpression.build(class, code, ternary.True, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, ternary.True.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if stack > maxStack {
		maxStack = stack
	}
	trueExit := (&cg.Exit{}).FromCode(cg.OP_goto, code)
	context.MakeStackMap(code, state, code.CodeLength)
	writeExits([]*cg.Exit{falseExit}, code.CodeLength)
	stack, es = buildExpression.build(class, code, ternary.False, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, ternary.False.ExpressionValue)
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
