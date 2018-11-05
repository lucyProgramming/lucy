package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildExpression) buildMethodCall(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.FieldMethodHandler != nil {
		return this.buildMethodCallOnFieldHandler(class, code, e, context, state)
	}
	switch call.Expression.Value.Type {
	case ast.VariableTypeArray:
		return this.buildMethodCallOnArray(class, code, e, context, state)
	case ast.VariableTypeMap:
		return this.buildMethodCallOnMap(class, code, e, context, state)
	case ast.VariableTypeJavaArray:
		return this.buildMethodCallJavaOnArray(class, code, e, context, state)
	case ast.VariableTypePackage:
		return this.buildMethodCallOnPackage(class, code, e, context, state)
	case ast.VariableTypeDynamicSelector:
		return this.buildMethodCallOnDynamicSelector(class, code, e, context, state)
	case ast.VariableTypeClass:
		if call.Method.Function.JvmDescriptor == "" {
			call.Method.Function.JvmDescriptor = Descriptor.methodDescriptor(&call.Method.Function.Type)
		}
		maxStack = this.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: call.Method.Function.JvmDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := this.jvmSize(e); t > maxStack {
			maxStack = t
		}
		if t := popCallResult(code, e, &call.Method.Function.Type); t > maxStack {
			maxStack = t
		}
		return
	case ast.VariableTypeObject, ast.VariableTypeString:
		if call.Method.Function.JvmDescriptor == "" {
			call.Method.Function.JvmDescriptor = Descriptor.methodDescriptor(&call.Method.Function.Type)
		}
		maxStack = this.build(class, code, call.Expression, context, state)
		// object ref
		state.pushStack(class, call.Expression.Value)
		defer state.popStack(1)
		if call.Name == ast.SpecialMethodInit {
			state.popStack(1)
			v := &cg.StackMapUninitializedThisVariableInfo{} // make_node_objects it right
			state.Stacks = append(state.Stacks, &cg.StackMapVerificationTypeInfo{
				Verify: v,
			})
		}
		stack := this.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		if t := stack + 1; t > maxStack {
			maxStack = t
		}
		if t := this.jvmSize(e); t > maxStack {
			maxStack = t
		}
		if call.Name == ast.SpecialMethodInit { // call father construction method
			code.Codes[code.CodeLength] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
				Class:      call.Class.Name,
				Method:     call.Name,
				Descriptor: call.Method.Function.JvmDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			return
		}
		if call.Class.IsInterface() {
			code.Codes[code.CodeLength] = cg.OP_invokeinterface
			class.InsertInterfaceMethodrefConst(cg.ConstantInfoInterfaceMethodrefHighLevel{
				Class:      call.Class.Name,
				Method:     call.Name,
				Descriptor: call.Method.Function.JvmDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = interfaceMethodArgsCount(&call.Method.Function.Type)
			code.Codes[code.CodeLength+4] = 0
			code.CodeLength += 5
		} else {
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
				Class:      call.Class.Name,
				Method:     call.Name,
				Descriptor: call.Method.Function.JvmDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if t := popCallResult(code, e, &call.Method.Function.Type); t > maxStack {
			maxStack = t
		}
		return
	default:
		panic(call.Expression.Value.TypeString())
	}
	return
}
func (this *BuildExpression) buildMethodCallOnFieldHandler(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.FieldMethodHandler.IsStatic() == false {
		stack := this.build(class, code, call.Expression, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_getfield
		code.CodeLength++
	} else {
		code.Codes[code.CodeLength] = cg.OP_getstatic
		code.CodeLength++
	}
	class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
		Class:      call.Expression.Value.Class.Name,
		Field:      call.Name,
		Descriptor: Descriptor.typeDescriptor(call.FieldMethodHandler.Type),
	}, code.Codes[code.CodeLength:code.CodeLength+2])
	code.CodeLength += 2
	state.pushStack(class, state.newObjectVariableType(javaMethodHandleClass))
	defer state.popStack(1)
	stack := this.buildCallArgs(
		class, code, call.Args, call.VArgs,
		context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      javaMethodHandleClass,
		Method:     methodHandleInvokeMethodName,
		Descriptor: Descriptor.methodDescriptor(call.FieldMethodHandler.Type.FunctionType),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if t := popCallResult(code, e, call.FieldMethodHandler.Type.FunctionType); t > maxStack {
		maxStack = t
	}
	return
}
func (this *BuildExpression) buildMethodCallOnDynamicSelector(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.FieldMethodHandler != nil {
		if call.FieldMethodHandler.IsStatic() == false {
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
			if 1 > maxStack {
				maxStack = 1
			}
			code.Codes[code.CodeLength] = cg.OP_getfield
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_getstatic
			code.CodeLength++
		}
		class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
			Class:      call.Expression.Value.Class.Name,
			Field:      call.Name,
			Descriptor: Descriptor.typeDescriptor(call.FieldMethodHandler.Type),
		}, code.Codes[code.CodeLength:code.CodeLength+2])
		code.CodeLength += 2
		state.pushStack(class, state.newObjectVariableType(javaMethodHandleClass))
		defer state.popStack(1)
		stack := this.buildCallArgs(class, code, call.Args, call.VArgs,
			context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaMethodHandleClass,
			Method:     methodHandleInvokeMethodName,
			Descriptor: Descriptor.methodDescriptor(call.FieldMethodHandler.Type.FunctionType),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := popCallResult(code, e, call.FieldMethodHandler.Type.FunctionType); t > maxStack {
			maxStack = t
		}
	} else {
		currentStack := uint16(0)
		if call.Method.IsStatic() == false {
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
			state.pushStack(class, state.newObjectVariableType(call.Expression.Value.Class.Name))
			defer state.popStack(1)
			currentStack = 1
		}
		stack := this.buildCallArgs(class, code, call.Args, call.VArgs,
			context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if call.Method.IsStatic() {
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			code.CodeLength++
		}
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      call.Expression.Value.Class.Name,
			Method:     call.Name,
			Descriptor: Descriptor.methodDescriptor(&call.Method.Function.Type),
		}, code.Codes[code.CodeLength:code.CodeLength+2])
		code.CodeLength += 2
		if t := popCallResult(code, e, &call.Method.Function.Type); t > maxStack {
			maxStack = t
		}
	}
	return
}
func (this *BuildExpression) buildMethodCallOnPackage(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.PackageFunction != nil {
		stack := this.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      call.Expression.Value.Package.Name + "/main",
			Method:     call.Name,
			Descriptor: Descriptor.methodDescriptor(&call.PackageFunction.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := popCallResult(code, e, &call.PackageFunction.Type); t > maxStack {
			maxStack = t
		}
	} else {
		//call.PackageGlobalVariableFunction != nil
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
			Class:      call.Expression.Value.Package.Name + "/main",
			Field:      call.Name,
			Descriptor: Descriptor.typeDescriptor(call.PackageGlobalVariableFunction.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		state.pushStack(class, call.PackageGlobalVariableFunction.Type)
		defer state.popStack(1)
		stack := this.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      "java/lang/invoke/MethodHandle",
			Method:     methodHandleInvokeMethodName,
			Descriptor: Descriptor.methodDescriptor(call.PackageGlobalVariableFunction.Type.FunctionType),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := popCallResult(code, e, call.PackageGlobalVariableFunction.Type.FunctionType); t > maxStack {
			maxStack = t
		}
	}
	return
}
