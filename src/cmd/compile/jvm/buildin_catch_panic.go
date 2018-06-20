package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) mkBuildInPanic(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Args[0].Type != ast.EXPRESSION_TYPE_NEW { // not new expression
		code.Codes[code.CodeLength] = cg.OP_new
		className := call.Args[0].ExpressionValue.Class.Name
		class.InsertClassConst(className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
		{
			t := &cg.StackMapVerificationTypeInfo{}
			tt := &cg.StackMapUninitializedVariableInfo{}
			tt.CodeOffset = uint16(code.CodeLength - 4)
			t.Verify = tt
			state.Stacks = append(state.Stacks, t)
			state.Stacks = append(state.Stacks, t)
		}
		stack, _ := makeExpression.build(class, code, call.Args[0], context, state)
		state.popStack(2)
		maxStack = 2 + stack
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      className,
			Method:     special_method_init,
			Descriptor: "(Ljava/lang/Throwable;)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		maxStack, _ = makeExpression.build(class, code, call.Args[0], context, state)
	}
	code.Codes[code.CodeLength] = cg.OP_athrow
	code.CodeLength++
	context.MakeStackMap(code, state, code.CodeLength)
	return
}

func (makeExpression *MakeExpression) mkBuildInCatch(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context) (maxStack uint16) {
	if e.IsStatementExpression { // statement call
		maxStack = 1
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVariableForException.Offset)...)
		return
	}
	maxStack = 2
	//load to stack
	copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVariableForException.Offset)...) // load
	//set 2 null
	code.Codes[code.CodeLength] = cg.OP_aconst_null
	code.CodeLength++
	copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVariableForException.Offset)...) // store
	//check cast
	code.Codes[code.CodeLength] = cg.OP_checkcast
	if context.Defer.ExceptionClass != nil {
		class.InsertClassConst(context.Defer.ExceptionClass.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		class.InsertClassConst(ast.DEFAULT_EXCEPTION_CLASS, code.Codes[code.CodeLength+1:code.CodeLength+3])
	}
	code.CodeLength += 3
	return
}
