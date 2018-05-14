package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildTernary(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	ternary := e.Data.(*ast.ExpressionTernary)
	var es []*cg.JumpBackPatch
	maxstack, es = m.build(class, code, ternary.Condition, context, state)
	if len(es) > 0 {
		backPatchEs(es, code.CodeLength)
		state.pushStack(class, ternary.Condition.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	exit := (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code)
	//true part
	stack, es := m.build(class, code, ternary.True, context, state)
	if len(es) > 0 {
		backPatchEs(es, code.CodeLength)
		state.pushStack(class, ternary.True.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if stack > maxstack {
		maxstack = stack
	}
	exit2 := (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)
	context.MakeStackMap(code, state, code.CodeLength)
	backPatchEs([]*cg.JumpBackPatch{exit}, code.CodeLength)
	stack, es = m.build(class, code, ternary.False, context, state)
	if len(es) > 0 {
		backPatchEs(es, code.CodeLength)
		state.pushStack(class, ternary.False.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if stack > maxstack {
		maxstack = stack
	}
	state.pushStack(class, e.Value)
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	backPatchEs([]*cg.JumpBackPatch{exit2}, code.CodeLength)
	return

}
