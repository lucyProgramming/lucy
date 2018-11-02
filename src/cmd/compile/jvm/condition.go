package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	compile condition for false  &&  generate exit
*/
func (buildExpression *BuildExpression) buildConditionNotOk(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	context *Context,
	state *StackMapState,
	condition *ast.Expression) (maxStack uint16, exit *cg.Exit) {
	if condition.Is2IntCompare() {
		return buildExpression.build2IntCompareConditionNotOk(class, code, context, state, condition)
	} else if condition.IsCompare2Null() {
		return buildExpression.buildNullCompareConditionNotOk(class, code, context, state, condition)
	} else if condition.Is2StringCompare() {
		return buildExpression.buildStringCompareConditionNotOk(class, code, context, state, condition)
	} else if condition.Is2PointerCompare() {
		return buildExpression.buildPointerCompareConditionNotOk(class, code, context, state, condition)
	} else {
		maxStack = buildExpression.build(class, code, condition, context, state)
		exit = (&cg.Exit{}).Init(cg.OP_ifeq, code)
		return
	}
}

func (buildExpression *BuildExpression) build2IntCompareConditionNotOk(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	context *Context,
	state *StackMapState,
	condition *ast.Expression) (maxStack uint16, exit *cg.Exit) {
	bin := condition.Data.(*ast.ExpressionBinary)
	stack := buildExpression.build(class, code, bin.Left, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	state.pushStack(class, bin.Left.Value)
	stack = buildExpression.build(class, code, bin.Right, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	state.popStack(1)
	switch condition.Type {
	case ast.ExpressionTypeEq:
		exit = (&cg.Exit{}).Init(cg.OP_if_icmpne, code)
	case ast.ExpressionTypeNe:
		exit = (&cg.Exit{}).Init(cg.OP_if_icmpeq, code)
	case ast.ExpressionTypeGe:
		exit = (&cg.Exit{}).Init(cg.OP_if_icmplt, code)
	case ast.ExpressionTypeGt:
		exit = (&cg.Exit{}).Init(cg.OP_if_icmple, code)
	case ast.ExpressionTypeLe:
		exit = (&cg.Exit{}).Init(cg.OP_if_icmpgt, code)
	case ast.ExpressionTypeLt:
		exit = (&cg.Exit{}).Init(cg.OP_if_icmpge, code)
	}
	return
}
func (buildExpression *BuildExpression) buildNullCompareConditionNotOk(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	context *Context,
	state *StackMapState,
	condition *ast.Expression) (maxStack uint16, exit *cg.Exit) {
	var noNullExpression *ast.Expression
	bin := condition.Data.(*ast.ExpressionBinary)
	if bin.Left.Type != ast.ExpressionTypeNull {
		noNullExpression = bin.Left
	} else {
		noNullExpression = bin.Right
	}
	stack := buildExpression.build(class, code, noNullExpression, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	switch condition.Type {
	case ast.ExpressionTypeEq:
		exit = (&cg.Exit{}).Init(cg.OP_ifnonnull, code)
	case ast.ExpressionTypeNe:
		exit = (&cg.Exit{}).Init(cg.OP_ifnull, code)
	}
	return
}
func (buildExpression *BuildExpression) buildStringCompareConditionNotOk(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	context *Context,
	state *StackMapState,
	condition *ast.Expression) (maxStack uint16, exit *cg.Exit) {
	bin := condition.Data.(*ast.ExpressionBinary)
	stack := buildExpression.build(class, code, bin.Left, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	state.pushStack(class, bin.Left.Value)
	stack = buildExpression.build(class, code, bin.Right, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      javaStringClass,
		Method:     "compareTo",
		Descriptor: "(Ljava/lang/String;)I",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.popStack(1)
	switch condition.Type {
	case ast.ExpressionTypeEq:
		exit = (&cg.Exit{}).Init(cg.OP_ifne, code)
	case ast.ExpressionTypeNe:
		exit = (&cg.Exit{}).Init(cg.OP_ifeq, code)
	case ast.ExpressionTypeGe:
		exit = (&cg.Exit{}).Init(cg.OP_iflt, code)
	case ast.ExpressionTypeGt:
		exit = (&cg.Exit{}).Init(cg.OP_ifle, code)
	case ast.ExpressionTypeLe:
		exit = (&cg.Exit{}).Init(cg.OP_ifgt, code)
	case ast.ExpressionTypeLt:
		exit = (&cg.Exit{}).Init(cg.OP_ifge, code)
	}
	return
}
func (buildExpression *BuildExpression) buildPointerCompareConditionNotOk(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	context *Context,
	state *StackMapState,
	condition *ast.Expression) (maxStack uint16, exit *cg.Exit) {
	bin := condition.Data.(*ast.ExpressionBinary)
	stack := buildExpression.build(class, code, bin.Left, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	state.pushStack(class, bin.Left.Value)
	stack = buildExpression.build(class, code, bin.Right, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	switch condition.Type {
	case ast.ExpressionTypeEq:
		exit = (&cg.Exit{}).Init(cg.OP_if_acmpne, code)
	case ast.ExpressionTypeNe:
		exit = (&cg.Exit{}).Init(cg.OP_if_acmpeq, code)
	}
	state.popStack(1)
	return
}
