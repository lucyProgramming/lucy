package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildTypeAssert(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	assert := e.Data.(*ast.ExpressionTypeAssert)
	maxstack, _ = m.build(class, code, assert.Expression, context, state)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_instanceof
	class.InsertClassConst(assert.Typ.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 3 > maxstack {
		maxstack = 3
	}
	{
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, assert.Expression.Value)...)
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.popStack(2)
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, assert.Expression.Value)...)
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

	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_arrylist_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	if 4 > maxstack {
		maxstack = 4
	}
	//
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// store object
	code.Codes[code.CodeLength] = cg.OP_dup_x1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     "add",
		Descriptor: "(Ljava/lang/Object;)Z",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	// store if ok
	code.Codes[code.CodeLength] = cg.OP_dup_x1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, &ast.VariableType{Typ: ast.VARIABLE_TYPE_BOOL})
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     "add",
		Descriptor: "(Ljava/lang/Object;)Z",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	return
}
