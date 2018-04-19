package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context, state *StackMapState) (maxstack uint16) {
	if s.StatmentForRangeAttr != nil {
		if s.StatmentForRangeAttr.Expression.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY ||
			s.StatmentForRangeAttr.Expression.VariableType.Typ == ast.VARIABLE_TYPE_JAVA_ARRAY {
			return m.buildForRangeStatementForArray(class, code, s, context, state)
		} else { // for map
			return m.buildForRangeStatementForMap(class, code, s, context, state)
		}
	}
	//init
	if s.Init != nil {
		stack, _ := m.MakeExpression.build(class, code, s.Init, context, state)
		if stack > maxstack {
			maxstack = stack
		}
	}
	forState := (&StackMapState{}).FromLast(state)
	loopBeginAt := code.CodeLength
	s.ContinueOPOffset = code.CodeLength
	context.MakeStackMap(code, state, code.CodeLength)
	//condition
	if s.Condition != nil {
		stack, es := m.MakeExpression.build(class, code, s.Condition, context, state)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			forState.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, s.Condition.VariableType)...)
			context.MakeStackMap(code, state, code.CodeLength)
			forState.popStack(1) // must be bool expression
		}
		if stack > maxstack {
			maxstack = stack
		}
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code))
	} else {

	}
	if s.Condition != nil {
		m.buildBlock(class, code, s.Block, context, forState)
	} else {
		m.buildBlock(class, code, s.Block, context, state)
	}
	if s.Post != nil {
		s.ContinueOPOffset = code.CodeLength
		//stack is here
		context.MakeStackMap(code, forState, code.CodeLength)
		stack, _ := m.MakeExpression.build(class, code, s.Post, context, forState)
		if stack > maxstack {
			maxstack = stack
		}
	}
	jumpto(cg.OP_goto, code, loopBeginAt)
	return
}
