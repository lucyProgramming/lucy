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
	//condition
	buildIntCompareCondition := func() {
		bin := s.Condition.Data.(*ast.ExpressionBinary)
		stack := buildPackage.BuildExpression.build(class, code, bin.Left, context, forState)
		if stack > maxStack {
			maxStack = stack
		}
		forState.pushStack(class, bin.Left.Value)
		stack = buildPackage.BuildExpression.build(class, code, bin.Right, context, forState)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		switch s.Condition.Type {
		case ast.ExpressionTypeEq:
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_if_icmpne, code))
		case ast.ExpressionTypeNe:
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_if_icmpeq, code))
		case ast.ExpressionTypeGe:
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_if_icmplt, code))
		case ast.ExpressionTypeGt:
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_if_icmple, code))
		case ast.ExpressionTypeLe:
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_if_icmpgt, code))
		case ast.ExpressionTypeLt:
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_if_icmpge, code))
		}
		forState.popStack(1)
	}
	var firstTimeExit *cg.Exit
	if s.Condition != nil {
		if s.Condition.Is2IntCompare() {
			buildIntCompareCondition()
		} else {
			stack := buildPackage.BuildExpression.build(class, code, s.Condition, context, forState)
			if stack > maxStack {
				maxStack = stack
			}
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_ifeq, code))
		}
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
		if s.Condition.Is2IntCompare() {
			buildIntCompareCondition()
		} else {
			stack := buildPackage.BuildExpression.build(class, code, s.Condition, context, forState)
			if stack > maxStack {
				maxStack = stack
			}
			s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_ifeq, code))
		}
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
