package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildNew(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	n := e.Data.(*ast.ExpressionNew)
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses(n.Typ.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2
	stackneed := maxstack
	size := uint16(0)
	for _, v := range n.Args {
		if v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || ast.EXPRESSION_TYPE_METHOD_CALL == v.Typ {
			panic(1)
		}
		size = e.VariableType.JvmSlotSize()
		stack, es := m.build(class, code, v, context)
		if stackneed+stack > maxstack {
			maxstack = stackneed + stack
		}
		stackneed += size
		backPatchEs(es, code.CodeLength)
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	methodref := cg.CONSTANT_Methodref_info_high_level{
		Class:      n.Typ.Class.Name,
		Name:       n.Construction.Func.Name,
		Descriptor: n.Construction.Func.Descriptor,
	}
	class.InsertMethodRefConst(methodref, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
