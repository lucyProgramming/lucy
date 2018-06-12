package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildTemplateFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.TemplateFunctionCallPair.TemplateFunctionCallPairGenerated == nil {
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
		m.MakeClass.buildFunction(class, nil, method, call.TemplateFunctionCallPair.Function)
		call.TemplateFunctionCallPair.TemplateFunctionCallPairGenerated = method

	}
	maxstack = m.buildCallArgs(class, code, call.Args,
		call.TemplateFunctionCallPair.Function.Typ.ParameterList, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokestatic
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      call.TemplateFunctionCallPair.TemplateFunctionCallPairGenerated.Class.Name,
		Method:     call.TemplateFunctionCallPair.TemplateFunctionCallPairGenerated.Name,
		Descriptor: call.TemplateFunctionCallPair.TemplateFunctionCallPairGenerated.Descriptor,
	},
		code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
