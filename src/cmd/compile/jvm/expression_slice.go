package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildSlice(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	slice := e.Data.(*ast.ExpressionSlice)
	maxstack, _ = m.build(class, code, slice.Expression, context)
	// build start
	stack, _ := m.build(class, code, slice.Start, context)
	if t := 1 + stack; t > maxstack {
		maxstack = t
	}
	stack, _ = m.build(class, code, slice.End, context)
	if t := 2 + stack; t > maxstack {
		maxstack = t
	}
	meta := ArrayMetas[e.VariableType.CombinationType.Typ]
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Name:       "slice",
		Descriptor: meta.sliceDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
