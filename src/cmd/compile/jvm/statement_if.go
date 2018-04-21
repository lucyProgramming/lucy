package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildIfStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementIF, context *Context, state *StackMapState) (maxstack uint16) {
	var es []*cg.JumpBackPatch
	IfState := (&StackMapState{}).FromLast(state)
	maxstack, es = m.MakeExpression.build(class, code, s.Condition, context, IfState)
	if len(es) > 0 {
		backPatchEs(es, code.CodeLength)
		IfState.Stacks = append(IfState.Stacks,
			IfState.newStackMapVerificationTypeInfo(class, s.Condition.Value)...)
		context.MakeStackMap(code, IfState, code.CodeLength)
		IfState.popStack(1) // must be bool expression
	}
	code.Codes[code.CodeLength] = cg.OP_ifeq
	codelength := code.CodeLength
	falseExit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	m.buildBlock(class, code, s.Block, context, IfState)
	if len(s.ElseIfList) > 0 || s.ElseBlock != nil {
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	}
	state.addTop(IfState)
	for k, v := range s.ElseIfList {
		context.MakeStackMap(code, state, code.CodeLength) // state is not change,all block var should be access from outside
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
		elseIfState := (&StackMapState{}).FromLast(state)
		stack, es := m.MakeExpression.build(class, code, v.Condition, context, elseIfState)
		if len(es) > 0 {
			elseIfState.Stacks = append(elseIfState.Stacks,
				IfState.newStackMapVerificationTypeInfo(class, s.Condition.Value)...)
			backPatchEs(es, code.CodeLength)
			elseIfState.popStack(1)
		}
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		codelength = code.CodeLength
		falseExit = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		code.CodeLength += 3
		m.buildBlock(class, code, v.Block, context, elseIfState)
		if k != len(s.ElseIfList)-1 || s.ElseBlock != nil {
			s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
		}
		// when done
		state.addTop(elseIfState)
	}
	if s.ElseBlock != nil {
		context.MakeStackMap(code, state, code.CodeLength)
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
		falseExit = nil
		elseState := (&StackMapState{}).FromLast(state)
		m.buildBlock(class, code, s.ElseBlock, context, elseState)
		state.addTop(elseState)
	}
	if falseExit != nil {
		context.MakeStackMap(code, state, code.CodeLength)
		binary.BigEndian.PutUint16(falseExit, uint16(code.CodeLength-codelength))
	}
	return
}
