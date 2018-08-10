package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildIfStatement(class *cg.ClassHighLevel,
	code *cg.AttributeCode, s *ast.StatementIf, context *Context, state *StackMapState) (maxStack uint16) {
	conditionState := (&StackMapState{}).FromLast(state)
	defer state.addTop(conditionState)
	for _, v := range s.PrefixExpressions {
		stack := buildPackage.BuildExpression.build(class, code, v, context, conditionState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	var IfState *StackMapState
	if s.Block.HaveVariableDefinition() {
		IfState = (&StackMapState{}).FromLast(conditionState)
	} else {
		IfState = conditionState
	}
	stack := buildPackage.BuildExpression.build(class, code, s.Condition, context, IfState)
	if stack > maxStack {
		maxStack = stack
	}
	exit := (&cg.Exit{}).Init(cg.OP_ifeq, code)
	buildPackage.buildBlock(class, code, &s.Block, context, IfState)
	conditionState.addTop(IfState)
	if s.ElseBlock != nil || len(s.ElseIfList) > 0 {
		if s.Block.WillNotExecuteToEnd == false {
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
		}
	}
	for k, v := range s.ElseIfList {
		context.MakeStackMap(code, conditionState, code.CodeLength) // state is not change,all block var should be access from outside
		writeExits([]*cg.Exit{exit}, code.CodeLength)
		var elseIfState *StackMapState
		if v.Block.HaveVariableDefinition() {
			elseIfState = (&StackMapState{}).FromLast(conditionState)
		} else {
			elseIfState = conditionState
		}
		stack := buildPackage.BuildExpression.build(class, code, v.Condition, context, elseIfState)
		if stack > maxStack {
			maxStack = stack
		}
		exit = (&cg.Exit{}).Init(cg.OP_ifeq, code)
		buildPackage.buildBlock(class, code, v.Block, context, elseIfState)
		if s.ElseBlock != nil || k != len(s.ElseIfList)-1 {
			if v.Block.WillNotExecuteToEnd == false {
				s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
			}
		}
		// when done
		conditionState.addTop(elseIfState)
	}
	context.MakeStackMap(code, conditionState, code.CodeLength)
	writeExits([]*cg.Exit{exit}, code.CodeLength)
	if s.ElseBlock != nil {
		var elseBlockState *StackMapState
		if s.ElseBlock.HaveVariableDefinition() {
			elseBlockState = (&StackMapState{}).FromLast(conditionState)
		} else {
			elseBlockState = conditionState
		}
		buildPackage.buildBlock(class, code, s.ElseBlock, context, elseBlockState)
		conditionState.addTop(elseBlockState)
		if s.ElseBlock.WillNotExecuteToEnd == false {
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
		}
	}
	return
}
