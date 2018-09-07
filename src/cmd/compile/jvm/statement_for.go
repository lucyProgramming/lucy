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
		stack := buildPackage.BuildExpression.build(class, code, s.Init, context, forState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	var firstTimeExit *cg.Exit
	if s.Condition != nil {
		stack, exit := buildPackage.BuildExpression.buildConditionNotOk(class, code, context, forState, s.Condition)
		if stack > maxStack {
			maxStack = stack
		}
		s.Exits = append(s.Exits, exit)
		firstTimeExit = (&cg.Exit{}).Init(cg.OP_goto, code) // goto body
	}
	s.ContinueCodeOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	if s.Increment != nil {
		stack := buildPackage.BuildExpression.build(class, code, s.Increment, context, forState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	if s.Condition != nil {
		stack, exit := buildPackage.BuildExpression.buildConditionNotOk(class, code, context, forState, s.Condition)
		if stack > maxStack {
			maxStack = stack
		}
		s.Exits = append(s.Exits, exit)
	}
	if firstTimeExit != nil {
		writeExits([]*cg.Exit{firstTimeExit}, code.CodeLength)
		context.MakeStackMap(code, forState, code.CodeLength)
	}
	buildPackage.buildBlock(class, code, s.Block, context, forState)
	if s.Block.WillNotExecuteToEnd == false {
		jumpTo(code, s.ContinueCodeOffset)
		if s.Condition == nil {
			context.MakeStackMap(code, forState, code.CodeLength)
		}
	}
	return
}
