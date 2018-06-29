package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.Expression.ExpressionValue.Type == ast.VariableTypeArray {
		return buildExpression.buildArrayMethodCall(class, code, e, context, state)
	}
	if call.Expression.ExpressionValue.Type == ast.VariableTypeMap {
		return buildExpression.buildMapMethodCall(class, code, e, context, state)
	}
	if call.Expression.ExpressionValue.Type == ast.VariableTypeJavaArray {
		return buildExpression.buildJavaArrayMethodCall(class, code, e, context, state)
	}

	pop := func(f *ast.Function) {
		if e.IsStatementExpression && f.NoReturnValue() == false {
			if len(e.ExpressionMultiValues) == 1 {
				if jvmSlotSize(e.ExpressionMultiValues[0]) == 1 {
					code.Codes[code.CodeLength] = cg.OP_pop
					code.CodeLength++
				} else {
					code.Codes[code.CodeLength] = cg.OP_pop2
					code.CodeLength++
				}
			} else { // > 1
				code.Codes[code.CodeLength] = cg.OP_pop
				code.CodeLength++
			}
		}
	}

	d := call.Method.Function.Descriptor
	if call.Class.LoadFromOutSide == false {
		d = JvmDescriptor.methodDescriptor(&call.Method.Function.Type)
	}
	if call.Method.IsStatic() {
		maxStack = buildExpression.buildCallArgs(class, code, call.Args, call.Method.Function.Type.ParameterList, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := buildExpression.expressionValueJvmSize(e); t > maxStack {
			maxStack = t
		}
		pop(call.Method.Function)
		return
	}

	maxStack, _ = buildExpression.build(class, code, call.Expression, context, state)
	// object ref
	state.pushStack(class, call.Expression.ExpressionValue)
	defer state.popStack(1)
	if call.Name == ast.ConstructionMethodName {
		state.popStack(1)
		v := &cg.StackMapUninitializedThisVariableInfo{} // make it right
		state.Stacks = append(state.Stacks, &cg.StackMapVerificationTypeInfo{
			Verify: v,
		})
	}
	stack := buildExpression.buildCallArgs(class, code, call.Args, call.Method.Function.Type.ParameterList, context, state)
	if t := stack + 1; t > maxStack {
		maxStack = t
	}
	if t := buildExpression.expressionValueJvmSize(e); t > maxStack {
		maxStack = t
	}
	if call.Name == ast.ConstructionMethodName { // call father construction method
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	if call.Class.IsInterface() {
		code.Codes[code.CodeLength] = cg.OP_invokeinterface
		class.InsertInterfaceMethodrefConst(cg.CONSTANT_InterfaceMethodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = interfaceMethodArgsCount(&call.Method.Function.Type)
		code.Codes[code.CodeLength+4] = 0
		code.CodeLength += 5
	} else {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	pop(call.Method.Function)
	return
}
