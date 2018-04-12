package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildIfStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementIF, context *Context) (maxstack uint16) {
	var es []*cg.JumpBackPatch
	maxstack, es = m.MakeExpression.build(class, code, s.Condition, context)
	backPatchEs(es, code.CodeLength)
	code.Codes[code.CodeLength] = cg.OP_ifeq
	context.Stacks = context.Stacks[0 : len(context.Stacks)-1] // should be 0
	codelength := code.CodeLength
	falseExit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	var stackMapState StackMapStateLocalsNumber
	stackMapState.FromContext(context)
	m.buildBlock(class, code, s.Block, context)
	if len(s.ElseIfList) > 0 || s.ElseBlock != nil {
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	}
	for k, v := range s.ElseIfList {
		// ifeq will goto there
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(&stackMapState, code.CodeLength))
		stackMapState.FromContext(context)
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
		stack, es := m.MakeExpression.build(class, code, v.Condition, context)
		backPatchEs(es, code.CodeLength)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		context.Stacks = context.Stacks[0 : len(context.Stacks)-1] // should be 0
		codelength = code.CodeLength
		falseExit = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		code.CodeLength += 3
		m.buildBlock(class, code, v.Block, context)
		if k != len(s.ElseIfList)-1 || s.ElseBlock != nil {
			s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
		}
	}
	if s.ElseBlock != nil {
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
			context.MakeStackMap(&stackMapState, code.CodeLength))
		stackMapState.FromContext(context)
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
		falseExit = nil
		m.buildBlock(class, code, s.ElseBlock, context)
	}
	if falseExit != nil {
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
			context.MakeStackMap(&stackMapState, code.CodeLength))
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
	}
	return
}
