package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildTernary(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	ternary := e.Data.(*ast.ExpressionTernary)
	var es []*cg.Exit
	maxStack, es = makeExpression.build(class, code, ternary.Condition, context, state)
	if len(es) > 0 {
		backfillExit(es, code.CodeLength)
		state.pushStack(class, ternary.Condition.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	exit := (&cg.Exit{}).FromCode(cg.OP_ifeq, code)
	//true part
	stack, es := makeExpression.build(class, code, ternary.True, context, state)
	if len(es) > 0 {
		backfillExit(es, code.CodeLength)
		state.pushStack(class, ternary.True.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if stack > maxStack {
		maxStack = stack
	}
	exit2 := (&cg.Exit{}).FromCode(cg.OP_goto, code)
	context.MakeStackMap(code, state, code.CodeLength)
	backfillExit([]*cg.Exit{exit}, code.CodeLength)
	stack, es = makeExpression.build(class, code, ternary.False, context, state)
	if len(es) > 0 {
		backfillExit(es, code.CodeLength)
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
	backfillExit([]*cg.Exit{exit2}, code.CodeLength)
	return

}
