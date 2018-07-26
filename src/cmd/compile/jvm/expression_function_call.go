package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildFunctionPointerCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	maxStack = buildExpression.build(class, code, call.Expression, context, state)
	stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/invoke/MethodHandle",
		Method:     functionPointerInvokeMethod,
		Descriptor: Descriptor.methodDescriptor(call.Expression.Value.FunctionType),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if e.IsStatementExpression {
		if call.Expression.Value.FunctionType.NoReturnValue() == false {
			if len(call.Expression.Value.FunctionType.ReturnList) == 1 {
				if jvmSlotSize(call.Expression.Value.FunctionType.ReturnList[0].Type) == 1 {
					code.Codes[code.CodeLength] = cg.OP_pop
					code.CodeLength++
				} else {
					code.Codes[code.CodeLength] = cg.OP_pop2
					code.CodeLength++
				}
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop
				code.CodeLength++
			}
		}
		return
	}
	return
}
func (buildExpression *BuildExpression) buildFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Function == nil {
		return buildExpression.buildFunctionPointerCall(class, code, e, context, state)
	}
	if call.Function.IsBuildIn {
		return buildExpression.mkBuildInFunctionCall(class, code, e, context, state)
	}
	if call.Function.TemplateFunction != nil {
		return buildExpression.buildTemplateFunctionCall(class, code, e, context, state)
	}

	if call.Expression != nil {
		if call.Expression.Type == ast.ExpressionTypeFunctionLiteral {
			maxStack = buildExpression.build(class, code, call.Expression, context, state)
		}
	}
	if call.Function.IsClosureFunction == false {
		maxStack = buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Function.ClassMethod.Class.Name,
			Method:     call.Function.ClassMethod.Name,
			Descriptor: call.Function.ClassMethod.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		//closure function call
		//load object
		if context.function.Closure.ClosureFunctionExist(call.Function) {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      call.Function.Name,
				Descriptor: "L" + call.Function.ClassMethod.Class.Name + ";",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, call.Function.ClosureVariableOffSet)...)
		}
		state.pushStack(class, state.newObjectVariableType(call.Function.ClassMethod.Class.Name))
		defer state.popStack(1)
		stack := buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Function.ClassMethod.Class.Name,
			Method:     call.Function.Name,
			Descriptor: call.Function.ClassMethod.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	if e.IsStatementExpression {
		if e.CallHasReturnValue() == false {
			// nothing to do
		} else if len(e.MultiValues) == 1 {
			if 2 == jvmSlotSize(e.MultiValues[0]) {
				code.Codes[code.CodeLength] = cg.OP_pop2
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop
			}
			code.CodeLength++
		} else { // > 1
			code.Codes[code.CodeLength] = cg.OP_pop // array list object on stack
			code.CodeLength++
		}
	}
	if e.CallHasReturnValue() == false { // nothing
	} else if len(e.MultiValues) == 1 {
		if t := jvmSlotSize(e.MultiValues[0]); t > maxStack {
			maxStack = t
		}
	} else { // > 1
		if 1 > maxStack {
			maxStack = 1
		}
	}
	return
}
