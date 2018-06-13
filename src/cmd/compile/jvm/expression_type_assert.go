package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildTypeAssert(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	assertOn := e.Data.(*ast.ExpressionTypeAssert)
	maxStack, _ = makeExpression.build(class, code, assertOn.Expression, context, state)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_instanceof
	if assertOn.Type.Type == ast.VARIABLE_TYPE_OBJECT {
		class.InsertClassConst(assertOn.Type.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else if assertOn.Type.Type == ast.VARIABLE_TYPE_ARRAY { // arrays
		meta := ArrayMetas[assertOn.Type.ArrayType.Type]
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		class.InsertClassConst(Descriptor.typeDescriptor(assertOn.Type), code.Codes[code.CodeLength+1:code.CodeLength+3])
	}
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4

	{
		state.pushStack(class, assertOn.Expression.ExpressionValue)
		state.pushStack(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.popStack(2)
		state.pushStack(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
		state.pushStack(class, assertOn.Expression.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength+11)
		state.popStack(2)
	}
	code.Codes[code.CodeLength] = cg.OP_ifeq
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
	code.Codes[code.CodeLength+3] = cg.OP_swap
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7)
	code.Codes[code.CodeLength+7] = cg.OP_pop
	code.Codes[code.CodeLength+8] = cg.OP_pop
	code.Codes[code.CodeLength+9] = cg.OP_iconst_0
	code.Codes[code.CodeLength+10] = cg.OP_aconst_null
	code.CodeLength += 11
	loadInt(class, code, 2)
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst(java_root_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3

	// store object
	code.Codes[code.CodeLength] = cg.OP_dup_x1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	if 5 > maxStack {
		maxStack = 5
	}
	code.Codes[code.CodeLength] = cg.OP_aastore
	code.CodeLength++

	// store if ok
	code.Codes[code.CodeLength] = cg.OP_dup_x1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	typeConverter.putPrimitiveInObject(class, code, &ast.VariableType{Type: ast.VARIABLE_TYPE_BOOL})
	code.Codes[code.CodeLength] = cg.OP_iconst_1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_aastore
	code.CodeLength++
	return
}
