package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) mkBuildInAssert(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	length := int32(len(call.Args))
	lengthOffset := code.MaxLocals
	code.MaxLocals++
	state.appendLocals(class, &ast.Type{
		Type: ast.VariableTypeInt,
	})
	loadInt32(class, code, length)
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, lengthOffset)...)
	stepOffset := code.MaxLocals
	code.MaxLocals++
	state.appendLocals(class, &ast.Type{
		Type: ast.VariableTypeInt,
	})
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, stepOffset)...)
	exits := []*cg.Exit{}
	for _, a := range call.Args {
		stack := buildExpression.build(class, code, a, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		exits = append(exits, (&cg.Exit{}).Init(cg.OP_ifeq, code))
		code.Codes[code.CodeLength] = cg.OP_iinc
		code.Codes[code.CodeLength+1] = byte(stepOffset)
		code.Codes[code.CodeLength+2] = 1
		code.CodeLength += 3
	}
	writeExits(exits, code.CodeLength)
	context.MakeStackMap(code, state, code.CodeLength)
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, lengthOffset)...)
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, stepOffset)...)
	okExit := (&cg.Exit{}).Init(cg.OP_if_icmpeq, code)
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst("assert failed,expression->'%d'", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	loadInt32(class, code, 1)
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst(javaRootClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, stepOffset)...)
	if 5 > maxStack {
		maxStack = 5
	}
	typeConverter.packPrimitives(class, code, &ast.Type{
		Type: ast.VariableTypeInt,
	})
	code.Codes[code.CodeLength] = cg.OP_aastore
	code.CodeLength++
	class.InsertMethodCall(code, cg.OP_invokestatic, javaStringClass,
		"format", "(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;")

	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(javaExceptionClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_dup_x1
	code.Codes[code.CodeLength+1] = cg.OP_swap
	code.CodeLength += 2
	class.InsertMethodCall(code, cg.OP_invokespecial, javaExceptionClass, specialMethodInit, "(Ljava/lang/String;)V")
	code.Codes[code.CodeLength] = cg.OP_athrow
	code.CodeLength++
	writeExits([]*cg.Exit{okExit}, code.CodeLength)
	context.MakeStackMap(code, state, code.CodeLength)
	return
}
