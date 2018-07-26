package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.Expression.Value.Type == ast.VariableTypeArray {
		return buildExpression.buildArrayMethodCall(class, code, e, context, state)
	}
	if call.Expression.Value.Type == ast.VariableTypeMap {
		return buildExpression.buildMapMethodCall(class, code, e, context, state)
	}
	if call.Expression.Value.Type == ast.VariableTypeJavaArray {
		return buildExpression.buildJavaArrayMethodCall(class, code, e, context, state)
	}
	pop := func(ft *ast.FunctionType) {
		if e.IsStatementExpression && ft.NoReturnValue() == false {
			if len(e.MultiValues) == 1 {
				if jvmSlotSize(e.MultiValues[0]) == 1 {
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
	if call.Expression.Value.Type == ast.VariableTypePackage {
		//if call.PackageFunction != nil {
		//	stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		//	if stack > maxStack {
		//		maxStack = stack
		//	}
		//	code.Codes[code.CodeLength] = cg.OP_invokestatic
		//	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		//		Class:      call.Expression.Value.Package.Name + "/main",
		//		Method:     call.Name,
		//		Descriptor: Descriptor.methodDescriptor(&call.PackageFunction.Type),
		//	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		//	code.CodeLength += 3
		//	pop(&call.PackageFunction.Type)
		//}
		if call.PackageGlobalVariableFunction != nil {
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      call.Expression.Value.Package.Name + "/main",
				Field:      call.Name,
				Descriptor: Descriptor.typeDescriptor(call.PackageGlobalVariableFunction.Type),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			state.pushStack(class, call.PackageGlobalVariableFunction.Type)
			defer state.popStack(1)
			stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
			if t := 1 + stack; t > maxStack {
				maxStack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/invoke/MethodHandle",
				Method:     functionPointerInvokeMethod,
				Descriptor: Descriptor.methodDescriptor(call.PackageGlobalVariableFunction.Type.FunctionType),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			pop(call.PackageGlobalVariableFunction.Type.FunctionType)
		}
		return
	}
	if call.FieldMethodHandler != nil {
		if call.FieldMethodHandler.IsStatic() == false {
			stack := buildExpression.build(class, code, call.Expression, context, state)
			if stack > maxStack {
				maxStack = stack
			}
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      call.Expression.Value.Class.Name,
				Field:      call.Name,
				Descriptor: Descriptor.typeDescriptor(call.FieldMethodHandler.Type),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else {
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      call.Expression.Value.Class.Name,
				Field:      call.Name,
				Descriptor: Descriptor.typeDescriptor(call.FieldMethodHandler.Type),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		state.pushStack(class, state.newObjectVariableType(javaMethodHandleClass))
		stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs,
			context, state)
		defer state.popStack(1)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaMethodHandleClass,
			Method:     functionPointerInvokeMethod,
			Descriptor: Descriptor.methodDescriptor(call.FieldMethodHandler.Type.FunctionType),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		pop(call.FieldMethodHandler.Type.FunctionType)
		return
	}
	d := call.Method.Function.Descriptor
	if call.Class.LoadFromOutSide == false {
		d = Descriptor.methodDescriptor(&call.Method.Function.Type)
	}
	if call.Method.IsStatic() {
		maxStack = buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := buildExpression.jvmSize(e); t > maxStack {
			maxStack = t
		}
		pop(&call.Method.Function.Type)
		return
	}

	maxStack = buildExpression.build(class, code, call.Expression, context, state)
	// object ref
	state.pushStack(class, call.Expression.Value)
	defer state.popStack(1)
	if call.Name == ast.SpecialMethodInit {
		state.popStack(1)
		v := &cg.StackMapUninitializedThisVariableInfo{} // make it right
		state.Stacks = append(state.Stacks, &cg.StackMapVerificationTypeInfo{
			Verify: v,
		})
	}
	stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
	if t := stack + 1; t > maxStack {
		maxStack = t
	}
	if t := buildExpression.jvmSize(e); t > maxStack {
		maxStack = t
	}
	if call.Name == ast.SpecialMethodInit { // call father construction method
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
	pop(&call.Method.Function.Type)
	return
}
