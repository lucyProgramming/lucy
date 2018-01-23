package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16, exits [][]byte) {
	maxstack = 2
	bin := e.Data.(*ast.ExpressionBinary)
	stack, es := m.build(class, code, bin.Left, context)
	if stack > maxstack {
		maxstack = stack
	}
	backPatchEs(es, code)
	if e.Typ == ast.EXPRESSION_TYPE_LOGICAL_OR {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifne
		exits = [][]byte{code.Codes[code.CodeLength+2 : code.CodeLength+4]}
		code.Codes[code.CodeLength+4] = cg.OP_pop // pop 0 on stack
		code.CodeLength += 5
		stack, es = m.build(class, code, bin.Right, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
	} else { //and
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifeq
		exits = [][]byte{code.Codes[code.CodeLength+2 : code.CodeLength+4]}
		code.Codes[code.CodeLength+4] = cg.OP_pop // pop 1 on stack
		code.CodeLength += 5
		stack, es = m.build(class, code, bin.Right, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
	}
	return
}
