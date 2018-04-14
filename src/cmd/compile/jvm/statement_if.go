package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildIfStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementIF, context *Context, state *StackMapState) (maxstack uint16) {
	var es []*cg.JumpBackPatch
	maxstack, es = m.MakeExpression.build(class, code, s.Condition, context, state)
	backPatchEs(es, code.CodeLength)
	code.Codes[code.CodeLength] = cg.OP_ifeq
	codelength := code.CodeLength
	falseExit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	{
		ss := (&StackMapState{}).FromLast(state)
		m.buildBlock(class, code, s.Block, context, ss)
	}
	if len(s.ElseIfList) > 0 || s.ElseBlock != nil {
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	}
	for k, v := range s.ElseIfList {
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
			context.MakeStackMap(state, code.CodeLength)) // state is not change,all block var should be access from outside
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
		stack, es := m.MakeExpression.build(class, code, v.Condition, context, state)
		backPatchEs(es, code.CodeLength)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		codelength = code.CodeLength
		falseExit = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		code.CodeLength += 3
		ss := (&StackMapState{}).FromLast(state)
		m.buildBlock(class, code, v.Block, context, ss)
		if k != len(s.ElseIfList)-1 || s.ElseBlock != nil {
			s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
		}
	}
	if s.ElseBlock != nil {
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
			context.MakeStackMap(state, code.CodeLength)) // state is not change,all block var should be access from outside
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
		falseExit = nil
		ss := (&StackMapState{}).FromLast(state)
		m.buildBlock(class, code, s.ElseBlock, context, ss)
	}
	if falseExit != nil {
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
			context.MakeStackMap(state, code.CodeLength)) // state is not change,all block var should be access from outside
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
	}
	return
}
