package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildIfStatement(class *cg.ClassHighLevel,
	code *cg.AttributeCode, s *ast.StatementIf, context *Context, state *StackMapState) (maxStack uint16) {
	ifState := (&StackMapState{}).FromLast(state)
	defer state.addTop(ifState)
	for _, v := range s.PrefixExpressions {
		stack := buildPackage.BuildExpression.build(class, code, v, context, ifState)
		if stack > maxStack {
			maxStack = stack
		}
	}
	trueBlockState := (&StackMapState{}).FromLast(ifState)
	stack, exit :=
		buildPackage.BuildExpression.buildConditionNotOk(class, code, context, trueBlockState, s.Condition)
	if stack > maxStack {
		maxStack = stack
	}
	buildPackage.buildBlock(class, code, &s.TrueBlock, context, trueBlockState)
	ifState.addTop(trueBlockState)
	if s.ElseBlock != nil || len(s.ElseIfList) > 0 {
		if s.TrueBlock.NotExecuteToLastStatement == false {
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
		}
	}
	for k, v := range s.ElseIfList {
		context.MakeStackMap(code, ifState, code.CodeLength) // state is not change,all block var should be access from outside
		writeExits([]*cg.Exit{exit}, code.CodeLength)
		elseIfState := (&StackMapState{}).FromLast(ifState)
		stack, exit =
			buildPackage.BuildExpression.buildConditionNotOk(class, code, context, elseIfState, v.Condition)
		if stack > maxStack {
			maxStack = stack
		}
		buildPackage.buildBlock(class, code, v.Block, context, elseIfState)
		if s.ElseBlock != nil || k != len(s.ElseIfList)-1 {
			if v.Block.NotExecuteToLastStatement == false {
				s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
			}
		}
		// when done
		ifState.addTop(elseIfState)
	}
	context.MakeStackMap(code, ifState, code.CodeLength)
	writeExits([]*cg.Exit{exit}, code.CodeLength)
	if s.ElseBlock != nil {
		elseBlockState := (&StackMapState{}).FromLast(ifState)
		buildPackage.buildBlock(class, code, s.ElseBlock, context, elseBlockState)
		ifState.addTop(elseBlockState)
	}
	return
}
