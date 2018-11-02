package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildArray(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	arr := e.Data.(*ast.ExpressionArray)
	//	new array
	meta := ArrayMetas[e.Value.Array.Type]
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	{
		verify := &cg.StackMapVerificationTypeInfo{}
		unInit := &cg.StackMapUninitializedVariableInfo{}
		unInit.CodeOffset = uint16(code.CodeLength - 4)
		verify.Verify = unInit
		state.Stacks = append(state.Stacks, verify, verify)
	}
	loadInt32(class, code, int32(len(arr.Expressions)))
	newArrayBaseOnType(class, code, e.Value.Array)
	arrayObject := &ast.Type{
		Type:  ast.VariableTypeJavaArray,
		Array: e.Value.Array,
	}
	state.pushStack(class, arrayObject)
	maxStack = 3
	storeOP := storeArrayElementOp(e.Value.Array.Type)
	var index int32 = 0
	for _, v := range arr.Expressions {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, index) // load index
		state.pushStack(class, arrayObject)
		state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
		stack := buildExpression.build(class, code, v, context, state)
		state.popStack(2)
		if t := 5 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = storeOP
		code.CodeLength++
		index++
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      meta.className,
		Method:     specialMethodInit,
		Descriptor: meta.constructorFuncDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
