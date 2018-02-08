package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	//init
	if s.Init != nil {
		code.MKLineNumber(s.Init.Pos.StartLine)
		stack, _ := m.MakeExpression.build(class, code, s.Init, context)
		if stack > maxstack {
			maxstack = stack
		}
	}
	s.LoopBegin = code.CodeLength
	s.ContinueOPOffset = s.LoopBegin
	//condition
	if s.Condition != nil {
		code.MKLineNumber(s.Condition.Pos.StartLine)
		stack, es := m.MakeExpression.build(class, code, s.Condition, context)
		backPatchEs(es, code.CodeLength)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		b := cg.JumpBackPatch{}
		b.CurrentCodeLength = code.CodeLength
		b.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		s.BackPatchs = append(s.BackPatchs, &b)
		code.CodeLength += 3
	} else {
	}
	m.buildBlock(class, code, s.Block, context)
	if s.Post != nil {
		code.MKLineNumber(s.Post.Pos.StartLine)
		s.ContinueOPOffset = code.CodeLength
		stack, _ := m.MakeExpression.build(class, code, s.Post, context)
		if stack > maxstack {
			maxstack = stack
		}
	}
	jumpto(cg.OP_goto, code, s.LoopBegin)
	return
}

//func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
//	//init
//	if s.Init != nil {
//		stack, _ := m.MakeExpression.build(class, code, s.Init, context)
//		if stack > maxstack {
//			maxstack = stack
//		}
//	}
//	s.LoopBegin = code.CodeLength
//	//condition
//	if s.Condition != nil {
//		stack, es := m.MakeExpression.build(class, code, s.Condition, context)
//		backPatchEs(es, code.CodeLength)
//		if stack > maxstack {
//			maxstack = stack
//		}
//		code.Codes[code.CodeLength] = cg.OP_ifeq
//		b := cg.JumpBackPatch{}
//		b.CurrentCodeLength = code.CodeLength
//		b.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
//		s.BackPatchs = append(s.BackPatchs, &b)
//		code.CodeLength += 3
//	} else {
//	}
//	m.buildBlock(class, code, s.Block, context)
//	if s.Post != nil {
//		stack, _ := m.MakeExpression.build(class, code, s.Post, context)
//		if stack > maxstack {
//			maxstack = stack
//		}
//	}
//	jumpto(cg.OP_goto, code, s.LoopBegin)
//	return
//}
