package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildMapLiteral(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	values := e.Data.(*ast.ExpressionMap).Values
	{
		t := state.newObjectVariableType(java_hashmap_class)
		state.pushStack(class, t)
		state.pushStack(class, t)
	}
	for _, v := range values {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack := uint16(2)
		variableType := v.Left.Value
		if v.Left.MayHaveMultiValue() {
			variableType = v.Left.Values[0]
		}
		stack, _ := m.build(class, code, v.Left, context, state)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if variableType.IsPointer() == false {
			typeConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
		}
		state.pushStack(class, state.newObjectVariableType(java_root_class))
		currentStack = 3 // stack is ... mapref mapref kref
		stack, es := m.build(class, code, v.Right, context, state)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			state.pushStack(class, v.Right.Value)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		state.popStack(1) // @46 line
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		variableType = v.Right.Value
		if variableType.IsPointer() == false {
			typeConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
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
