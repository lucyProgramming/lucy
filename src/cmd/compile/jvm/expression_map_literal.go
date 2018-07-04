package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildMapLiteral(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(javaMapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxStack = 2
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaMapClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	values := e.Data.(*ast.ExpressionMap).KeyValuePairs

	hashMapObject := state.newObjectVariableType(javaMapClass)
	state.pushStack(class, hashMapObject)

	for _, v := range values {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack := uint16(2)
		state.pushStack(class, hashMapObject)
		stack, _ := buildExpression.build(class, code, v.Left, context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Left.ExpressionValue.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.Left.ExpressionValue)
		}
		state.pushStack(class, state.newObjectVariableType(javaRootClass))
		currentStack = 3 // stack is ... mapref mapref kref
		stack, es := buildExpression.build(class, code, v.Right, context, state)
		if len(es) > 0 {
			writeExits(es, code.CodeLength)
			state.pushStack(class, v.Right.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		state.popStack(1) // @43 line
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Right.ExpressionValue.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.Right.ExpressionValue)
		}
		// put in hashmap
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaMapClass,
			Method:     "put",
			Descriptor: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_pop
		code.CodeLength += 4
		state.popStack(1) // @35
	}
	return
}
