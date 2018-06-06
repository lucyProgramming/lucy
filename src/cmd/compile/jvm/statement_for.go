package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	s *ast.StatementFor, context *Context, state *StackMapState) (maxstack uint16) {
	if s.RangeAttr != nil {
		if s.RangeAttr.RangeOn.Value.Typ == ast.VARIABLE_TYPE_ARRAY ||
			s.RangeAttr.RangeOn.Value.Typ == ast.VARIABLE_TYPE_JAVA_ARRAY {
			return m.buildForRangeStatementForArray(class, code, s, context, state)
		} else { // for map
			return m.buildForRangeStatementForMap(class, code, s, context, state)
		}
	}
	forState := (&StackMapState{}).FromLast(state)
	defer func() {
		state.addTop(forState)
	}()
	//init
	if s.Init != nil {
		stack, _ := m.MakeExpression.build(class, code, s.Init, context, forState)
		if stack > maxstack {
			maxstack = stack
		}
	}
	//condition
	var firstExit *cg.JumpBackPatch
	if s.Condition != nil {
		stack, es := m.MakeExpression.build(class, code, s.Condition, context, forState)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			forState.pushStack(class, s.Condition.Value)
			context.MakeStackMap(code, forState, code.CodeLength)
			forState.popStack(1) // must be bool expression
		}
		if stack > maxstack {
			maxstack = stack
		}
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code))
		firstExit = (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)
	}
	s.ContinueOPOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	if s.Post != nil {
		stack, _ := m.MakeExpression.build(class, code, s.Post, context, forState)
		if stack > maxstack {
			maxstack = stack
		}
	}
	if s.Condition != nil {
		stack, es := m.MakeExpression.build(class, code, s.Condition, context, forState)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			forState.pushStack(class, s.Condition.Value)
			context.MakeStackMap(code, forState, code.CodeLength)
			forState.popStack(1) // must be bool expression
		}
		if stack > maxstack {
			maxstack = stack
		}
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code))
	}
	if firstExit != nil {
		backPatchEs([]*cg.JumpBackPatch{firstExit}, code.CodeLength)
		context.MakeStackMap(code, forState, code.CodeLength)
	}
	m.buildBlock(class, code, s.Block, context, forState)
	if s.Block.DeadEnding == false {
		jumpTo(cg.OP_goto, code, s.ContinueOPOffset)
	}
	return
}
