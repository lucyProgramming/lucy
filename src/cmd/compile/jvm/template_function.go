package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildTemplateFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.TemplateFunctionCallPair.Entrance == nil {
		method := &cg.MethodHighLevel{}
		method.Class = class
		method.Name = class.NewMethodName(nameTemplateFunction(call.TemplateFunctionCallPair.Function))
		method.AccessFlags |= cg.ACC_CLASS_PUBLIC
		method.AccessFlags |= cg.ACC_CLASS_FINAL
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_BRIDGE
		if call.TemplateFunctionCallPair.Function.Type.VArgs != nil {
			method.AccessFlags |= cg.ACC_METHOD_VARARGS
		}
		method.Descriptor = Descriptor.methodDescriptor(&call.TemplateFunctionCallPair.Function.Type)
		method.Code = &cg.AttributeCode{}
		class.AppendMethod(method)
		call.TemplateFunctionCallPair.Function.Entrance = method
		//build function
		buildExpression.BuildPackage.buildFunction(class, nil, method, call.TemplateFunctionCallPair.Function)
		call.TemplateFunctionCallPair.Entrance = method
	}
	maxStack = buildExpression.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokestatic
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      call.TemplateFunctionCallPair.Entrance.Class.Name,
		Method:     call.TemplateFunctionCallPair.Entrance.Name,
		Descriptor: call.TemplateFunctionCallPair.Entrance.Descriptor,
	},
		code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
