package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildIfStatement(class *cg.ClassHighLevel,
	code *cg.AttributeCode, s *ast.StatementIf, context *Context, state *StackMapState) (maxStack uint16) {
	var es []*cg.Exit
	conditionState := (&StackMapState{}).FromLast(state)
	defer state.addTop(conditionState)
	for _, v := range s.PrefixExpressions {
		stack, _ := buildPackage.BuildExpression.build(class, code, v, context, conditionState)
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
	maxStack, es = buildPackage.BuildExpression.build(class, code, s.Condition, context, IfState)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		IfState.pushStack(class, s.Condition.Value)
		context.MakeStackMap(code, IfState, code.CodeLength)
		IfState.popStack(1) // must be bool expression
	}
	code.Codes[code.CodeLength] = cg.OP_ifeq
	codeLength := code.CodeLength
	exit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	buildPackage.buildBlock(class, code, &s.Block, context, IfState)
	conditionState.addTop(IfState)
	if s.ElseBlock != nil || len(s.ElseIfList) > 0 {
		if s.Block.WillNotExecuteToEnd == false {
			s.Exits = append(s.Exits, (&cg.Exit{}).FromCode(cg.OP_goto, code))
		}
	}

	for k, v := range s.ElseIfList {
		context.MakeStackMap(code, conditionState, code.CodeLength) // state is not change,all block var should be access from outside
		binary.BigEndian.PutUint16(exit, uint16(code.CodeLength-codeLength))
		var elseIfState *StackMapState
		if v.Block.HaveVariableDefinition() {
			elseIfState = (&StackMapState{}).FromLast(conditionState)
		} else {
			elseIfState = conditionState
		}
		stack, es := buildPackage.BuildExpression.build(class, code, v.Condition, context, elseIfState)
		if len(es) > 0 {
			writeExits(es, code.CodeLength)
			elseIfState.pushStack(class, s.Condition.Value)
			context.MakeStackMap(code, elseIfState, code.CodeLength)
			elseIfState.popStack(1)
		}
		if stack > maxStack {
			maxStack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		codeLength = code.CodeLength
		exit = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		code.CodeLength += 3
		buildPackage.buildBlock(class, code, v.Block, context, elseIfState)
		if s.ElseBlock != nil || k != len(s.ElseIfList)-1 {
			if v.Block.WillNotExecuteToEnd == false {
				s.Exits = append(s.Exits, (&cg.Exit{}).FromCode(cg.OP_goto, code))
			}
		}
		// when done
		conditionState.addTop(elseIfState)
	}
	context.MakeStackMap(code, conditionState, code.CodeLength)
	binary.BigEndian.PutUint16(exit, uint16(code.CodeLength-codeLength))
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
			s.Exits = append(s.Exits, (&cg.Exit{}).FromCode(cg.OP_goto, code))
		}
	}
	return
}
