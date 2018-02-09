package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildIfStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementIF, context *Context) (maxstack uint16) {
	var es []*cg.JumpBackPatch
	//code.MKLineNumber(s.Condition.Pos.StartLine)
	maxstack, es = m.MakeExpression.build(class, code, s.Condition, context)
	backPatchEs(es, code.CodeLength)
	code.Codes[code.CodeLength] = cg.OP_ifeq
	codelength := code.CodeLength
	falseExit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	m.buildBlock(class, code, s.Block, context)
	s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	for _, v := range s.ElseIfList {
		binary.BigEndian.PutUint16(falseExit, code.CodeLength-codelength)
		//code.MKLineNumber(v.Condition.Pos.StartLine)
		stack, es := m.MakeExpression.build(class, code, v.Condition, context)
		backPatchEs(es, code.CodeLength)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		codelength = code.CodeLength
		falseExit = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		code.CodeLength += 3
		m.buildBlock(class, code, v.Block, context)
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	}
	if s.ElseBlock != nil {
		binary.BigEndian.PutUint16(falseExit, code.CodeLength-codelength)
		falseExit = nil
		m.buildBlock(class, code, s.ElseBlock, context)
	}
	if falseExit != nil {
		binary.BigEndian.PutUint16(falseExit, code.CodeLength-codelength)
	}
	return
}
