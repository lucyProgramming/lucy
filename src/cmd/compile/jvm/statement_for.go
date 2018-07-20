package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	s *ast.StatementFor, context *Context, state *StackMapState) (maxStack uint16) {
	if s.RangeAttr != nil {
		if s.RangeAttr.RangeOn.Value.Type == ast.VariableTypeArray ||
			s.RangeAttr.RangeOn.Value.Type == ast.VariableTypeJavaArray {
			return buildPackage.buildForRangeStatementForArray(class, code, s, context, state)
		} else { // for map
			return buildPackage.buildForRangeStatementForMap(class, code, s, context, state)
		}
	}
	forState := (&StackMapState{}).FromLast(state)
	defer func() {
		state.addTop(forState)
	}()
	//init
	if s.Init != nil {
		stack, _ := buildPackage.BuildExpression.build(class, code, s.Init, context, forState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	//condition
	var firstTimeExit *cg.Exit
	if s.Condition != nil {
		stack, es := buildPackage.BuildExpression.build(class, code, s.Condition, context, forState)
		if len(es) > 0 {
			writeExits(es, code.CodeLength)
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
	s.ContinueCodeOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	if s.Increment != nil {
		stack, _ := buildPackage.BuildExpression.build(class, code, s.Increment, context, forState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	if s.Condition != nil {
		stack, es := buildPackage.BuildExpression.build(class, code, s.Condition, context, forState)
		if len(es) > 0 {
			writeExits(es, code.CodeLength)
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
		writeExits([]*cg.Exit{firstTimeExit}, code.CodeLength)
		context.MakeStackMap(code, forState, code.CodeLength)
	}
	buildPackage.buildBlock(class, code, s.Block, context, forState)
	if s.Block.WillNotExecuteToEnd == false {
		jumpTo(cg.OP_goto, code, s.ContinueCodeOffset)
	}
	return
}
