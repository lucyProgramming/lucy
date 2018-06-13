package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildTemplateFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.TemplateFunctionCallPair.Generated == nil {
		method := &cg.MethodHighLevel{}
		method.Class = class
		method.Name = class.NewFunctionName(nameTemplateFunction(call.TemplateFunctionCallPair.Function))
		method.AccessFlags |= cg.ACC_CLASS_PUBLIC
		method.AccessFlags |= cg.ACC_CLASS_FINAL
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_BRIDGE
		method.Descriptor = Descriptor.methodDescriptor(call.TemplateFunctionCallPair.Function)
		method.Code = &cg.AttributeCode{}
		class.AppendMethod(method)
		call.TemplateFunctionCallPair.Function.ClassMethod = method
		//build function
		makeExpression.MakeClass.buildFunction(class, nil, method, call.TemplateFunctionCallPair.Function)
		call.TemplateFunctionCallPair.Generated = method

	}
	maxStack = makeExpression.buildCallArgs(class, code, call.Args,
		call.TemplateFunctionCallPair.Function.Type.ParameterList, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokestatic
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      call.TemplateFunctionCallPair.Generated.Class.Name,
		Method:     call.TemplateFunctionCallPair.Generated.Name,
		Descriptor: call.TemplateFunctionCallPair.Generated.Descriptor,
	},
		code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
