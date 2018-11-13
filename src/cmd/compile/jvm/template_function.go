package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildExpression) buildTemplateFunctionCall(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.TemplateFunctionCallPair.Entrance == nil {
		method := &cg.MethodHighLevel{}
		method.Class = class
		method.Name = class.NewMethodName(nameTemplateFunction(call.TemplateFunctionCallPair.Function))
		method.AccessFlags |= cg.AccClassPublic
		method.AccessFlags |= cg.AccClassFinal
		method.AccessFlags |= cg.AccMethodStatic
		method.AccessFlags |= cg.AccMethodBridge
		if call.TemplateFunctionCallPair.Function.Type.VArgs != nil {
			method.AccessFlags |= cg.AccMethodVarargs
		}
		method.Descriptor = Descriptor.methodDescriptor(&call.TemplateFunctionCallPair.Function.Type)
		method.Code = &cg.AttributeCode{}
		class.AppendMethod(method)
		call.TemplateFunctionCallPair.Function.Entrance = method
		//build function
		this.BuildPackage.buildFunction(class, nil, method, call.TemplateFunctionCallPair.Function)
		call.TemplateFunctionCallPair.Entrance = method
	}
	maxStack = this.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokestatic
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      call.TemplateFunctionCallPair.Entrance.Class.Name,
		Method:     call.TemplateFunctionCallPair.Entrance.Name,
		Descriptor: call.TemplateFunctionCallPair.Entrance.Descriptor,
	},
		code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if e.IsStatementExpression {
		if call.TemplateFunctionCallPair.Function.Type.VoidReturn() == false {
			if len(call.TemplateFunctionCallPair.Function.Type.ReturnList) > 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
				code.CodeLength++
			} else {
				if jvmSlotSize(e.Value) == 1 {
					code.Codes[code.CodeLength] = cg.OP_pop
					code.CodeLength++
				} else {
					code.Codes[code.CodeLength] = cg.OP_pop2
					code.CodeLength++
				}
			}
		}
	}
	return
}
