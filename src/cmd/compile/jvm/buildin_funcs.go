package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	switch call.Func.Name {
	case common.BUILD_IN_FUNCTION_PRINT:
		return m.mkBuildinPrint(class, code, call, context, state)
	case common.BUILD_IN_FUNCTION_PANIC:
		return m.mkBuildinPanic(class, code, call, context, state)
	case common.BUILD_IN_FUNCTION_CATCH:
		return m.mkBuildinRecover(class, code, e, context)
	case common.BUILD_IN_FUNCTION_MONITORENTER, common.BUILD_IN_FUNCTION_MONITOREXIT:
		maxstack, _ = m.build(class, code, call.Args[0], context, state)
		if call.Func.Name == common.BUILD_IN_FUNCTION_MONITORENTER {
			code.Codes[code.CodeLength] = cg.OP_monitorenter
		} else {
			code.Codes[code.CodeLength] = cg.OP_monitorexit
		}
		code.CodeLength++
	default:
		panic("unkown buildin function:" + call.Func.Name)
	}
	return
}

func (m *MakeExpression) mkBuildinPanic(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall,
	context *Context, state *StackMapState) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_throwable_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	{
		t := &cg.StackMap_verification_type_info{}
		tt := &cg.StackMap_Uninitialized_variable_info{}
		tt.Index = uint16(code.CodeLength - 4)
		t.PayLoad = tt
		state.Stacks = append(state.Stacks, t)
		state.Stacks = append(state.Stacks, t)
		defer state.popStack(2)
	}
	stack, _ := m.build(class, code, call.Args[0], context, state)
	maxstack = 2 + stack
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_throwable_class,
		Method:     special_method_init,
		Descriptor: "(Ljava/lang/Throwable;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_athrow
	code.CodeLength++
	return
}

func (m *MakeExpression) mkBuildinRecover(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	if e.IsStatementExpression { // statement call
		maxstack = 1
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...)
		return
	}
	maxstack = 2
	//load to stack
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...) // load
	//set 2 null
	code.Codes[code.CodeLength] = cg.OP_aconst_null
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...) // load
	return
}
