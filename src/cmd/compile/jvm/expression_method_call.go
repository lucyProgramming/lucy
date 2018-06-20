package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_ARRAY {
		return makeExpression.buildArrayMethodCall(class, code, e, context, state)
	}
	if call.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_MAP {
		return makeExpression.buildMapMethodCall(class, code, e, context, state)
	}
	if call.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_JAVA_ARRAY {
		return makeExpression.buildJavaArrayMethodCall(class, code, e, context, state)
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
	if call.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_PACKAGE {
		maxStack = makeExpression.buildCallArgs(class, code, call.Args, call.PackageFunction.Type.ParameterList, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Expression.ExpressionValue.Package.Name + "/main",
			Method:     call.Name,
			Descriptor: call.PackageFunction.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := makeExpression.valueJvmSize(e); t > maxStack {
			maxStack = t
		}
		pop(call.PackageFunction)
		return
	}

	d := call.Method.Func.Descriptor
	if call.Class.LoadFromOutSide == false {
		d = Descriptor.methodDescriptor(call.Method.Func)
	}
	if call.Method.IsStatic() {
		maxStack = makeExpression.buildCallArgs(class, code, call.Args, call.Method.Func.Type.ParameterList, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := makeExpression.valueJvmSize(e); t > maxStack {
			maxStack = t
		}
		pop(call.Method.Func)
		return
	}

	maxStack, _ = makeExpression.build(class, code, call.Expression, context, state)
	// object ref
	state.pushStack(class, call.Expression.ExpressionValue)
	defer state.popStack(1)
	if call.Name == ast.CONSTRUCTION_METHOD_NAME {
		state.popStack(1)
		v := &cg.StackMapUninitializedThisVariableInfo{} // make it right
		state.Stacks = append(state.Stacks, &cg.StackMapVerificationTypeInfo{
			Verify: v,
		})
	}
	stack := makeExpression.buildCallArgs(class, code, call.Args, call.Method.Func.Type.ParameterList, context, state)
	if t := stack + 1; t > maxStack {
		maxStack = t
	}
	if t := makeExpression.valueJvmSize(e); t > maxStack {
		maxStack = t
	}
	if call.Name == ast.CONSTRUCTION_METHOD_NAME { // call father construction method
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
		code.Codes[code.CodeLength+3] = interfaceMethodArgsCount(&call.Method.Func.Type)
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
	pop(call.Method.Func)
	return
}
