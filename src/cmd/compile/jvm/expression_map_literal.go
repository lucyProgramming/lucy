package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildMapLiteral(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 2
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	values := e.Data.(*ast.ExpressionMap).Values
	for _, v := range values {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack := uint16(2)
		variableType := v.Left.VariableType
		if v.Left.MayHaveMultiValue() {
			variableType = v.Left.VariableTypes[0]
		}
		stack, _ := m.build(class, code, v.Left, context)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if variableType.IsPointer() == false {
			primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
		}
		currentStack = 3 // stack is ... mapref mapref kref
		stack, es := m.build(class, code, v.Right, context)
		backPatchEs(es, code.CodeLength)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		variableType = v.Right.VariableType
		if v.Right.MayHaveMultiValue() {
			variableType = v.Right.VariableTypes[0]
		}
		if variableType.IsPointer() == false {
			primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
		}
		// put in hashmap
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Method:     "put",
			Descriptor: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_pop
		code.CodeLength += 4
	}
	return
}
