package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildFunctionPointerCall(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	maxStack = buildExpression.build(class, code, call.Expression, context, state)
	stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      "java/lang/invoke/MethodHandle",
		Method:     methodHandleInvokeMethodName,
		Descriptor: Descriptor.methodDescriptor(call.Expression.Value.FunctionType),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if t := popCallResult(code, e, call.Expression.Value.FunctionType); t > maxStack {
		maxStack = t
	}
	return
}
func (buildExpression *BuildExpression) buildFunctionCall(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Function == nil {
		return buildExpression.buildFunctionPointerCall(class, code, e, context, state)
	}
	if call.Function.TemplateFunction != nil {
		return buildExpression.buildTemplateFunctionCall(class, code, e, context, state)
	}
	if call.Function.IsBuildIn {
		return buildExpression.mkBuildInFunctionCall(class, code, e, context, state)
	}
	if call.Expression != nil &&
		call.Expression.Type == ast.ExpressionTypeFunctionLiteral {
		maxStack = buildExpression.build(class, code, call.Expression, context, state)
	}
	if call.Function.IsClosureFunction == false {
		maxStack = buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      call.Function.Entrance.Class.Name,
			Method:     call.Function.Entrance.Name,
			Descriptor: call.Function.Entrance.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		//closure function call
		//load object
		if context.function.Closure.ClosureFunctionExist(call.Function) {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
				Class:      class.Name,
				Field:      call.Function.Name,
				Descriptor: "L" + call.Function.Entrance.Class.Name + ";",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, call.Function.ClosureVariableOffSet)...)
		}
		state.pushStack(class, state.newObjectVariableType(call.Function.Entrance.Class.Name))
		defer state.popStack(1)
		stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      call.Function.Entrance.Class.Name,
			Method:     call.Function.Name,
			Descriptor: call.Function.Entrance.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	if t := popCallResult(code, e, &call.Function.Type); t > maxStack {
		maxStack = t
	}
	return
}
