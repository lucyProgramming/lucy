package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildExpressionAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	bin := e.Data.(*ast.ExpressionBinary)
	left := bin.Left.Data.([]*ast.Expression)[0]
	right := bin.Right.Data.([]*ast.Expression)[0]
	maxStack, remainStack, op, target, className, name, descriptor := makeExpression.getLeftValue(class, code, left, context, state)
	stack, es := makeExpression.build(class, code, right, context, state)
	if len(es) > 0 {
		state.pushStack(class, right.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		fillOffsetForExits(es, code.CodeLength)
	}
	if t := remainStack + stack; t > maxStack {
		maxStack = t
	}
	currentStack := remainStack + jvmSlotSize(target)
	if e.IsStatementExpression == false {
		currentStack += makeExpression.controlStack2FitAssign(code, op, className, target)
		if currentStack > maxStack {
			maxStack = currentStack
		}
	}
	copyLeftValueOps(class, code, op, className, name, descriptor)
	return
}

// a,b,c = 122,fdfd2232,"hello";
func (makeExpression *MakeExpression) buildAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	rights := bin.Right.Data.([]*ast.Expression)
	lefts := bin.Left.Data.([]*ast.Expression)
	if e.IsStatementExpression == false || len(lefts) == 1 {
		return makeExpression.buildExpressionAssign(class, code, e, context, state)
	}
	if len(rights) == 1 {
		maxStack, _ = makeExpression.build(class, code, rights[0], context, state)
	} else {
		maxStack = makeExpression.buildExpressions(class, code, rights, context, state)
	}
	multiValuePacker.storeArrayListAutoVar(code, context)
	for k, v := range lefts {
		stackLength := len(state.Stacks)
		stack, remainStack, op, target, className, name, descriptor :=
			makeExpression.getLeftValue(class, code, v, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		stack = multiValuePacker.unPack(class, code, k, target, context)
		if t := remainStack + stack; t > maxStack {
			maxStack = t
		}
		copyLeftValueOps(class, code, op, className, name, descriptor)
		state.popStack(len(state.Stacks) - stackLength)
	}
	return
}

func (makeExpression *MakeExpression) controlStack2FitAssign(code *cg.AttributeCode, op []byte, className string,
	stackTopType *ast.Type) (increment uint16) {
	if op[0] == cg.OP_istore ||
		op[0] == cg.OP_lstore ||
		op[0] == cg.OP_fstore ||
		op[0] == cg.OP_dstore ||
		op[0] == cg.OP_astore ||
		op[0] == cg.OP_istore_0 ||
		op[0] == cg.OP_istore_1 ||
		op[0] == cg.OP_istore_2 ||
		op[0] == cg.OP_istore_3 ||
		op[0] == cg.OP_lstore_0 ||
		op[0] == cg.OP_lstore_1 ||
		op[0] == cg.OP_lstore_2 ||
		op[0] == cg.OP_lstore_3 ||
		op[0] == cg.OP_fstore_0 ||
		op[0] == cg.OP_fstore_1 ||
		op[0] == cg.OP_fstore_2 ||
		op[0] == cg.OP_fstore_3 ||
		op[0] == cg.OP_dstore_0 ||
		op[0] == cg.OP_dstore_1 ||
		op[0] == cg.OP_dstore_2 ||
		op[0] == cg.OP_dstore_3 ||
		op[0] == cg.OP_astore_0 ||
		op[0] == cg.OP_astore_1 ||
		op[0] == cg.OP_astore_2 ||
		op[0] == cg.OP_astore_3 {
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup
		} else {
			code.Codes[code.CodeLength] = cg.OP_dup2
			increment = 2
		}
		code.CodeLength++
		return
	}
	if op[0] == cg.OP_putstatic {
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup
		} else {
			code.Codes[code.CodeLength] = cg.OP_dup2
			increment = 2
		}
		code.CodeLength++
		return
	}
	if op[0] == cg.OP_putfield {
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup_x1
		} else {
			increment = 2
			code.Codes[code.CodeLength] = cg.OP_dup2_x1
		}
		code.CodeLength++
		return
	}
	if op[0] == cg.OP_iastore ||
		op[0] == cg.OP_lastore ||
		op[0] == cg.OP_fastore ||
		op[0] == cg.OP_dastore ||
		op[0] == cg.OP_aastore ||
		op[0] == cg.OP_bastore ||
		op[0] == cg.OP_castore ||
		op[0] == cg.OP_sastore {
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup_x2
			code.CodeLength++
		} else {
			increment = 2
			code.Codes[code.CodeLength] = cg.OP_dup2_x2
			code.CodeLength++
		}
		return
	}
	/*
		it is a flag indicate  map destination
		stack are ... mapRef kRef
	*/
	if className == java_hashmap_class {
		if jvmSlotSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup_x2
			code.CodeLength++
		} else {
			increment = 2
			code.Codes[code.CodeLength] = cg.OP_dup2_x2
			code.CodeLength++
		}
		return
	}
	return
}
