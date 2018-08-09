package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildArray(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
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
		uninit := &cg.StackMapUninitializedVariableInfo{}
		uninit.CodeOffset = uint16(code.CodeLength - 4)
		verify.Verify = uninit
		state.Stacks = append(state.Stacks, verify, verify)
	}
	loadInt32(class, code, int32(arr.Length))
	newArrayBaseOnType(class, code, e.Value.Array)
	arrayObject := &ast.Type{}
	arrayObject.Type = ast.VariableTypeJavaArray
	arrayObject.Array = e.Value.Array
	state.pushStack(class, arrayObject)
	maxStack = 4
	storeOP := storeArrayElementOp(e.Value.Array.Type)
	var index int32 = 0
	for _, v := range arr.Expressions {
		if v.HaveMultiValue() {
			// stack top is array list
			stack := buildExpression.build(class, code, v, context, state)
			if t := 3 + stack; t > maxStack {
				maxStack = t
			}
			autoVar := newMultiValueAutoVar(class, code, state)
			for k, t := range v.MultiValues {
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				loadInt32(class, code, index) // load index
				stack := autoVar.unPack(class, code, k, t)
				if t := 5 + stack; t > maxStack {
					maxStack = t
				}
				code.Codes[code.CodeLength] = storeOP
				code.CodeLength++
				index++
			}
			continue
		}
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
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.className,
		Method:     specialMethodInit,
		Descriptor: meta.constructorFuncDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
