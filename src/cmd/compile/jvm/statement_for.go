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
	loopBeginAt := code.CodeLength
	s.ContinueOPOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	//condition
	stack, es := m.MakeExpression.build(class, code, s.Condition, context, forState)
	if len(es) > 0 {
		backPatchEs(es, code.CodeLength)
		state.pushStack(class, s.Condition.Value)
		context.MakeStackMap(code, forState, code.CodeLength)
		forState.popStack(1) // must be bool expression
	}
	if stack > maxstack {
		maxstack = stack
	}
	s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code))

	m.buildBlock(class, code, s.Block, context, forState)
	if s.Post != nil {
		s.ContinueOPOffset = code.CodeLength
		//stack is here
		context.MakeStackMap(code, forState, code.CodeLength)
		stack, _ := m.MakeExpression.build(class, code, s.Post, context, forState)
		if stack > maxstack {
			maxstack = stack
		}
		jumpTo(cg.OP_goto, code, loopBeginAt)
	} else {
		if s.Block.DeadEnding == false { //will execute to here
			jumpTo(cg.OP_goto, code, loopBeginAt)
		}
	}
	return
}
