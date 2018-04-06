package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildTypeAssert(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	assert := e.Data.(*ast.ExpressionTypeAssert)
	maxstack, _ = m.build(class, code, assert.Expression, context)
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

	code.Codes[code.CodeLength] = cg.OP_ifeq
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
	code.Codes[code.CodeLength+3] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 6)
	code.Codes[code.CodeLength+6] = cg.OP_swap
	code.Codes[code.CodeLength+7] = cg.OP_aconst_null
	code.Codes[code.CodeLength+8] = cg.OP_swap
	code.CodeLength += 9

	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_arrylist_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	if 5 > maxstack {
		maxstack = 5
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_dup_x1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_iconst_1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, e.VariableTypes[1])
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     "add",
		Descriptor: "(ILjava/lang/Object;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// store object
	code.Codes[code.CodeLength] = cg.OP_dup_x1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_iconst_1
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     "add",
		Descriptor: "(ILjava/lang/Object;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
