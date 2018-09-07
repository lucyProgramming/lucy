package jvm

import (
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
	var remainStack uint16
	var op []byte
	var leftValueKind LeftValueKind
	if left.IsIdentifier(ast.NoNameIdentifier) == false {
		maxStack, remainStack, op, leftValueKind =
			buildExpression.getLeftValue(class, code, left, context, state)
	}
	stack := buildExpression.build(class, code, right, context, state)
	if t := remainStack + stack; t > maxStack {
		maxStack = t
	}
	if left.IsIdentifier(ast.NoNameIdentifier) {
		if jvmSlotSize(right.Value) == 1 {
			code.Codes[code.CodeLength] = cg.OP_pop
		} else {
			code.Codes[code.CodeLength] = cg.OP_pop2
		}
		code.CodeLength++
	} else {
		currentStack := remainStack + jvmSlotSize(left.Value)
		if e.IsStatementExpression == false {
			currentStack += buildExpression.dupStackLeaveValueBelow(code, leftValueKind, left.Value)
			if currentStack > maxStack {
				maxStack = currentStack
			}
		}
		copyOPs(code, op...)
	}
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
		maxStack = buildExpression.build(class, code, rights[0], context, state)
	} else {
		maxStack = buildExpression.buildExpressions(class, code, rights, context, state)
	}
	autoVar := newMultiValueAutoVar(class, code, state)
	for k, v := range lefts {
		if v.IsIdentifier(ast.NoNameIdentifier) {
			continue
		}
		stackLength := len(state.Stacks)
		stack, remainStack, op, _ :=
			buildExpression.getLeftValue(class, code, v, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		stack = autoVar.unPack(class, code, k, v.Value)
		if t := remainStack + stack; t > maxStack {
			maxStack = t
		}
		copyOPs(code, op...)
		state.popStack(len(state.Stacks) - stackLength)
	}
	return
}

func (buildExpression *BuildExpression) dupStackLeaveValueBelow(code *cg.AttributeCode, leftValueKind LeftValueKind,
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
	}

	return
}
