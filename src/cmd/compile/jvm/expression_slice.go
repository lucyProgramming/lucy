package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildSlice(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	slice := e.Data.(*ast.ExpressionSlice)
	meta := ArrayMetas[slice.Expression.Value.ArrayType.Typ]
	maxstack, _ = m.build(class, code, slice.Expression, context, state)
	state.pushStack(class, state.newObjectVariableType(meta.classname))
	// build start
	stack, _ := m.build(class, code, slice.Start, context, state)
	if t := 1 + stack; t > maxstack {
		maxstack = t
	}
	state.pushStack(class, slice.Start.Value)
	stack, _ = m.build(class, code, slice.End, context, state)
	if t := 2 + stack; t > maxstack {
		maxstack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Method:     "slice",
		Descriptor: meta.sliceDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
