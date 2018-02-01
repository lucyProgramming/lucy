package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.Statement, context *Context) (maxstack uint16) {
	switch s.Typ {
	case ast.STATEMENT_TYPE_EXPRESSION:
		var es [][]byte
		maxstack, es = m.MakeExpression.build(class, code, s.Expression, context)
		backPatchEs(es, code)
	case ast.STATEMENT_TYPE_IF:
		maxstack = m.buildIfStatement(class, code, s.StatementIf, context)
		backPatchEs(s.StatementIf.BackPatchs, code)
	case ast.STATEMENT_TYPE_BLOCK:
		m.buildBlock(class, code, s.Block, context)
	case ast.STATEMENT_TYPE_FOR:
		maxstack = m.buildForStatement(class, code, s.StatementFor, context)
		backPatchEs(s.StatementFor.BackPatchs, code)
	case ast.STATEMENT_TYPE_CONTINUE:
		code.Codes[code.CodeLength] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[1:3], s.StatementFor.LoopBegin)
		code.CodeLength += 3
	case ast.STATEMENT_TYPE_BREAK:
		code.Codes[code.CodeLength] = cg.OP_goto
		if s.StatementBreak.StatementFor != nil {
			appendBackPatch(&s.StatementFor.BackPatchs, code.Codes[code.CodeLength+1:code.CodeLength+3])
		} else { // switch
			appendBackPatch(&s.StatementSwitch.BackPatchs, code.Codes[code.CodeLength+1:code.CodeLength+3])
		}
		code.CodeLength += 3
	case ast.STATEMENT_TYPE_RETURN:
		maxstack = m.buildReturnStatement(class, code, s.StatementReturn, context)
	case ast.STATEMENT_TYPE_SWITCH:
		maxstack = m.buildSwitchStatement(class, code, s.StatementSwitch, context)
		backPatchEs(s.StatementSwitch.BackPatchs, code)
	case ast.STATEMENT_TYPE_SKIP: // skip this block
		panic("no skip")
	}

	return
}

func (m *MakeClass) buildSwitchStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementSwitch, context *Context) (maxstack uint16) {
	return
}
