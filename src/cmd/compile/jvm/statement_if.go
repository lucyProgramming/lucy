package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildIfStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementIF, context *Context) (maxstack uint16) {
	stack, es := m.MakeExpression.build(class, code, s.Condition, context)
	backPatchEs(es, code)
	if stack > maxstack {
		maxstack = stack
	}
	code.Codes[code.CodeLength] = cg.OP_ifeq
	falseExit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	m.buildBlock(class, code, s.Block, context)
	for _, v := range s.ElseIfList {
		backPatchEs([][]byte{falseExit}, code)
		stack, es := m.MakeExpression.build(class, code, v.Condition, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		falseExit = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		code.CodeLength += 3
		m.buildBlock(class, code, v.Block, context)
	}
	if s.ElseBlock != nil {
		backPatchEs([][]byte{falseExit}, code)
		falseExit = nil
		m.buildBlock(class, code, s.ElseBlock, context)
	}
	if falseExit != nil {
		backPatchEs([][]byte{falseExit}, code)
	}
	return
}
