package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
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
		Name:       special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	defaultValues := e.Data.(*ast.ExpressionMap)
	for _, v := range defaultValues.Values {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack := uint16(2)
		variableType := v.Left.VariableType
		if v.Left.IsCall() {
			variableType = v.Left.VariableTypes[0]
		}
		stack, _ := m.build(class, code, v.Left, context)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if variableType.IsPointer() == false {
			PrimitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
		}
		currentStack = 3 // stack is ... mapref mapref kref
		stack, es := m.build(class, code, v.Right, context)
		backPatchEs(es, code.CodeLength)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		variableType = v.Right.VariableType
		if v.Left.IsCall() {
			variableType = v.Right.VariableTypes[0]
		}
		if variableType.IsPointer() == false {
			PrimitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
		}
		// put in hashmap
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Name:       "put",
			Descriptor: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_pop
		code.CodeLength += 4
	}
	return

}
