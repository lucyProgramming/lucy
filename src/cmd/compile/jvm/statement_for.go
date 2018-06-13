package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeClass *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	s *ast.StatementFor, context *Context, state *StackMapState) (maxStack uint16) {
	if s.RangeAttr != nil {
		if s.RangeAttr.RangeOn.Value.Type == ast.VARIABLE_TYPE_ARRAY ||
			s.RangeAttr.RangeOn.Value.Type == ast.VARIABLE_TYPE_JAVA_ARRAY {
			return makeClass.buildForRangeStatementForArray(class, code, s, context, state)
		} else { // for map
			return makeClass.buildForRangeStatementForMap(class, code, s, context, state)
		}
	}
	forState := (&StackMapState{}).FromLast(state)
	defer func() {
		state.addTop(forState)
	}()
	//init
	if s.Init != nil {
		stack, _ := makeClass.makeExpression.build(class, code, s.Init, context, forState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	//condition
	var firstTimeExit *cg.Exit
	if s.Condition != nil {
		stack, es := makeClass.makeExpression.build(class, code, s.Condition, context, forState)
		if len(es) > 0 {
			backfillExit(es, code.CodeLength)
			forState.pushStack(class, s.Condition.Value)
			context.MakeStackMap(code, forState, code.CodeLength)
			forState.popStack(1) // must be bool expression
		}
		if stack > maxStack {
			maxStack = stack
		}
		s.Exits = append(s.Exits, (&cg.Exit{}).FromCode(cg.OP_ifeq, code))
		firstTimeExit = (&cg.Exit{}).FromCode(cg.OP_goto, code)
	}
	s.ContinueOPOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	if s.Post != nil {
		stack, _ := makeClass.makeExpression.build(class, code, s.Post, context, forState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	if s.Condition != nil {
		stack, es := makeClass.makeExpression.build(class, code, s.Condition, context, forState)
		if len(es) > 0 {
			backfillExit(es, code.CodeLength)
			forState.pushStack(class, s.Condition.Value)
			context.MakeStackMap(code, forState, code.CodeLength)
			forState.popStack(1) // must be bool expression
		}
		if stack > maxStack {
			maxStack = stack
		}
		s.Exits = append(s.Exits, (&cg.Exit{}).FromCode(cg.OP_ifeq, code))
	}
	if firstTimeExit != nil {
		backfillExit([]*cg.Exit{firstTimeExit}, code.CodeLength)
		context.MakeStackMap(code, forState, code.CodeLength)
	}
	makeClass.buildBlock(class, code, s.Block, context, forState)
	if s.Block.DeadEnding == false {
		jumpTo(cg.OP_goto, code, s.ContinueOPOffset)
	}
	return
}
