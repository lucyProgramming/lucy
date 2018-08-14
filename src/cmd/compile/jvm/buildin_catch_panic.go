package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) mkBuildInPanic(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Args[0].Type == ast.ExpressionTypeNew { // not new expression
		maxStack = buildExpression.build(class, code, call.Args[0], context, state)
	} else {
		code.Codes[code.CodeLength] = cg.OP_new
		className := call.Args[0].Value.Class.Name
		class.InsertClassConst(className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
		{
			verificationTypeInfo := &cg.StackMapVerificationTypeInfo{}
			uninitializedVariableInfo := &cg.StackMapUninitializedVariableInfo{}
			uninitializedVariableInfo.CodeOffset = uint16(code.CodeLength - 4)
			verificationTypeInfo.Verify = uninitializedVariableInfo
			state.Stacks = append(state.Stacks, verificationTypeInfo)
			state.Stacks = append(state.Stacks, verificationTypeInfo)
		}
		stack := buildExpression.build(class, code, call.Args[0], context, state)
		state.popStack(2)
		maxStack = 2 + stack
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      className,
			Method:     specialMethodInit,
			Descriptor: "(Ljava/lang/Throwable;)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	code.Codes[code.CodeLength] = cg.OP_athrow
	code.CodeLength++
	context.MakeStackMap(code, state, code.CodeLength)
	return
}

func (buildExpression *BuildExpression) mkBuildInCatch(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context) (maxStack uint16) {
	if e.IsStatementExpression { // statement call
		maxStack = 1
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		copyOPs(code,
			storeLocalVariableOps(ast.VariableTypeObject, context.exceptionVarOffset)...)
	} else {
		maxStack = 2
		//load to stack
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, context.exceptionVarOffset)...) // load
		//set 2 null
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, context.exceptionVarOffset)...) // store
		//check cast
		code.Codes[code.CodeLength] = cg.OP_checkcast
		if context.Defer.ExceptionClass != nil {
			class.InsertClassConst(context.Defer.ExceptionClass.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		} else {
			class.InsertClassConst(ast.DefaultExceptionClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		}
		code.CodeLength += 3
	}
	return
}
