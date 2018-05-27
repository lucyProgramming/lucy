package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildExpressionAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	bin := e.Data.(*ast.ExpressionBinary)
	left := bin.Left.Data.([]*ast.Expression)[0]
	right := bin.Right.Data.([]*ast.Expression)[0]
	maxstack, remainStack, op, target, classname, name, descriptor := m.getLeftValue(class, code, left, context, state)
	stack, es := m.build(class, code, right, context, state)
	if len(es) > 0 {
		state.pushStack(class, right.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		backPatchEs(es, code.CodeLength)
	}
	if t := remainStack + stack; t > maxstack {
		maxstack = t
	}
	currentStack := remainStack + jvmSize(target)
	if e.IsStatementExpression == false {
		currentStack += m.controlStack2FitAssign(code, op, classname, target)
		if currentStack > maxstack {
			maxstack = currentStack
		}
	}
	if classname == java_hashmap_class && e.Value.IsPointer() == false {
		typeConverter.putPrimitiveInObject(class, code, target)
	}
	copyOPLeftValue(class, code, op, classname, name, descriptor)
	return
}

// a,b,c = 122,fdfd2232,"hello";
func (m *MakeExpression) buildAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	rights := bin.Right.Data.([]*ast.Expression)
	lefts := bin.Left.Data.([]*ast.Expression)
	if e.IsStatementExpression == false || len(lefts) == 1 {
		return m.buildExpressionAssign(class, code, e, context, state)
	}
	maxstack = m.buildExpressions(class, code, rights, context, state)
	arrayListPacker.storeArrayListAutoVar(code, context)
	for k, v := range lefts {
		stackLength := len(state.Stacks)
		stack, remainStack, op, target, classname, name, descriptor :=
			m.getLeftValue(class, code, v, context, state)
		if stack > maxstack {
			maxstack = stack
		}
		stack = arrayListPacker.unPack(class, code, k, target, context)
		if t := remainStack + stack; t > maxstack {
			maxstack = t
		}
		copyOPLeftValue(class, code, op, classname, name, descriptor)
		state.popStack(len(state.Stacks) - stackLength)
	}
	return
}

func (m *MakeExpression) controlStack2FitAssign(code *cg.AttributeCode, op []byte, classname string,
	stackTopType *ast.VariableType) (increment uint16) {
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
		if jvmSize(stackTopType) == 1 {
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
		if jvmSize(stackTopType) == 1 {
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
		if jvmSize(stackTopType) == 1 {
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
		if jvmSize(stackTopType) == 1 {
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
	if op[0] == cg.OP_invokevirtual && classname == java_hashmap_class {
		if jvmSize(stackTopType) == 1 {
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
