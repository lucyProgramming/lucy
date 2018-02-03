package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	//init
	if s.Init != nil {
		stack, es := m.MakeExpression.build(class, code, s.Init, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
	}
	s.LoopBegin = code.CodeLength
	//condition
	if s.Condition != nil {
		stack, es := m.MakeExpression.build(class, code, s.Condition, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		appendBackPatch(&s.BackPatchs, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {

	}
	m.buildBlock(class, code, s.Block, context)
	if s.Post != nil {
		stack, es := m.MakeExpression.build(class, code, s.Init, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
	}
	code.Codes[code.CodeLength] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], s.LoopBegin)
	code.CodeLength += 3
	panic(s.LoopBegin)
	return
}
