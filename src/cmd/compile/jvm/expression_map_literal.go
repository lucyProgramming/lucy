package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildMapLiteral(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(mapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxStack = 2
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      mapClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	values := e.Data.(*ast.ExpressionMap).KeyValuePairs
	hashMapObject := state.newObjectVariableType(mapClass)
	state.pushStack(class, hashMapObject)
	defer state.popStack(1)
	for _, v := range values {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack := uint16(2)
		state.pushStack(class, hashMapObject)
		stack := buildExpression.build(class, code, v.Key, context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Key.Value.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.Key.Value)
		}
		state.pushStack(class, state.newObjectVariableType(javaRootClass))
		currentStack = 3 // stack is ... mapref mapref kref
		stack = buildExpression.build(class, code, v.Value, context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Value.Value.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.Value.Value)
		}
		// put in hashmap
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      mapClass,
			Method:     "put",
			Descriptor: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_pop
		code.CodeLength += 4
		state.popStack(2)
	}
	return
}
