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
	maxstack, _ = m.build(class, code, slice.Expression, context, state)
	meta := ArrayMetas[slice.Expression.Value.ArrayType.Typ]
	state.pushStack(class, state.newObjectVariableType(meta.classname))
	// build start
	if slice.Start != nil {
		stack, _ := m.build(class, code, slice.Start, context, state)
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
	} else {
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
	}
	state.pushStack(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
	if slice.End != nil {
		stack, _ := m.build(class, code, slice.End, context, state)
		if t := 2 + stack; t > maxstack {
			maxstack = t
		}
	} else {
		code.Codes[code.CodeLength] = cg.OP_iconst_m1
		code.CodeLength++
		if 3 > maxstack { // stack top ->  ...arrayref start end
			maxstack = 3
		}
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
