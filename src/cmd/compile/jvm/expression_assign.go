package jvm

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildExpressionAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	bin := e.Data.(*ast.ExpressionBinary)
	left := bin.Left.Data.([]*ast.Expression)[0]
	right := bin.Right.Data.([]*ast.Expression)[0]
	maxStack, remainStack, op, target, LeftValueKind := buildExpression.getLeftValue(class, code, left, context, state)
	stack, es := buildExpression.build(class, code, right, context, state)
	if len(es) > 0 {
		state.pushStack(class, right.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		writeExits(es, code.CodeLength)
	}
	if t := remainStack + stack; t > maxStack {
		maxStack = t
	}
	currentStack := remainStack + jvmSlotSize(target)
	if e.IsStatementExpression == false {
		currentStack += buildExpression.controlStack2FitAssign(code, LeftValueKind, target)
		if currentStack > maxStack {
			maxStack = currentStack
		}
	}
	copyOPs(code, op...)
	return
}

// a,b,c = 122,fdfd2232,"hello";
func (buildExpression *BuildExpression) buildAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	rights := bin.Right.Data.([]*ast.Expression)
	lefts := bin.Left.Data.([]*ast.Expression)
	if e.IsStatementExpression == false || len(lefts) == 1 {
		return buildExpression.buildExpressionAssign(class, code, e, context, state)
	}
	if len(rights) == 1 {
		maxStack, _ = buildExpression.build(class, code, rights[0], context, state)
	} else {
		maxStack = buildExpression.buildExpressions(class, code, rights, context, state)
	}
	multiValuePacker.storeMultiValueAutoVar(code, context)
	for k, v := range lefts {
		stackLength := len(state.Stacks)
		stack, remainStack, op, target, _ :=
			buildExpression.getLeftValue(class, code, v, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		stack = multiValuePacker.unPack(class, code, k, target, context)
		if t := remainStack + stack; t > maxStack {
			maxStack = t
		}
		copyOPs(code, op...)
		state.popStack(len(state.Stacks) - stackLength)
	}
	return
}

func (buildExpression *BuildExpression) controlStack2FitAssign(code *cg.AttributeCode, leftValueKind LeftValueKind,
	stackTopType *ast.Type) (increment uint16) {
	switch leftValueKind {
	case LeftValueTypeLocalVar:
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup
		} else {
			code.Codes[code.CodeLength] = cg.OP_dup2
			increment = 2
		}
		code.CodeLength++
	case LeftValueTypePutStatic:
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup
		} else {
			code.Codes[code.CodeLength] = cg.OP_dup2
			increment = 2
		}
		code.CodeLength++

	case LeftValueTypePutField:
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup_x1
		} else {
			increment = 2
			code.Codes[code.CodeLength] = cg.OP_dup2_x1
		}
		code.CodeLength++

	case LeftValueTypeArray, LeftValueTypeLucyArray:
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup_x2
			code.CodeLength++
		} else {
			increment = 2
			code.Codes[code.CodeLength] = cg.OP_dup2_x2
			code.CodeLength++
		}

	case LeftValueTypeMap:
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup_x2
			code.CodeLength++
		} else {
			increment = 2
			code.Codes[code.CodeLength] = cg.OP_dup2_x2
			code.CodeLength++
		}
	default:
		panic(fmt.Sprintf("unknow %d", leftValueKind))
	}

	return
}
