package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16, exits []*cg.JumpBackPatch) {
	exits = []*cg.JumpBackPatch{}
	bin := e.Data.(*ast.ExpressionBinary)
	maxstack, es := m.build(class, code, bin.Left, context)
	if es != nil {
		exits = append(exits, es...)
	}
	if e.Typ == ast.EXPRESSION_TYPE_LOGICAL_OR {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		if 2 > maxstack { // dup increment stack
			maxstack = 2
		}
		code.Codes[code.CodeLength] = cg.OP_ifne // at this point,value is clear,leave 1 on stack
		exit := &cg.JumpBackPatch{}
		exit.CurrentCodeLength = code.CodeLength
		exit.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		exits = append(exits, exit)
		code.Codes[code.CodeLength+3] = cg.OP_pop // pop 0
		code.CodeLength += 4
		stack, es := m.build(class, code, bin.Right, context)
		if es != nil {
			exits = append(exits, es...)
		}
		if stack > maxstack {
			maxstack = stack
		}
	} else { //and
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		if 2 > maxstack { // dup increment stack
			maxstack = 2
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq // at this point,value is clear,leave 1 on stack
		exit := &cg.JumpBackPatch{}
		exit.CurrentCodeLength = code.CodeLength
		exit.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		exits = append(exits, exit)
		code.Codes[code.CodeLength+3] = cg.OP_pop // pop 1
		code.CodeLength += 4
		stack, es := m.build(class, code, bin.Right, context)
		if es != nil {
			exits = append(exits, es...)
		}
		if stack > maxstack {
			maxstack = stack
		}
	}
	return
}
