package jvm

import (
	//"fmt"
	//"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, s *ast.Statement,
	context *Context, state *StackMapState) (maxstack uint16) {
	//fmt.Println(s.Pos)
	switch s.Typ {
	case ast.STATEMENT_TYPE_EXPRESSION:
		if s.Expression.Typ == ast.EXPRESSION_TYPE_FUNCTION {
			return m.buildFunctionExpression(class, code, s.Expression, context, state)
		}
		maxstack, _ = m.MakeExpression.build(class, code, s.Expression, context, state)
	case ast.STATEMENT_TYPE_IF:
		s.StatementIf.BackPatchs = []*cg.JumpBackPatch{} //could compile multi times
		maxstack = m.buildIfStatement(class, code, s.StatementIf, context, state)
		if len(s.StatementIf.BackPatchs) > 0 {
			backPatchEs(s.StatementIf.BackPatchs, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_BLOCK: //new
		var ss *StackMapState
		if s.Block.HaveVariableDefinition() {
			ss = (&StackMapState{}).FromLast(state)
		} else {
			ss = state
		}
		m.buildBlock(class, code, s.Block, context, ss)
		state.addTop(ss)
	case ast.STATEMENT_TYPE_FOR:
		s.StatementFor.BackPatchs = []*cg.JumpBackPatch{} //could compile multi times
		maxstack = m.buildForStatement(class, code, s.StatementFor, context, state)
		if len(s.StatementFor.BackPatchs) > 0 {
			backPatchEs(s.StatementFor.BackPatchs, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_CONTINUE:
		m.buildDefers(class, code, context, s.StatementContinue.Defers, state)
		jumpTo(cg.OP_goto, code, s.StatementContinue.StatementFor.ContinueOPOffset)
	case ast.STATEMENT_TYPE_BREAK:
		m.buildDefers(class, code, context, s.StatementBreak.Defers, state)
		b := (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)
		if s.StatementBreak.StatementFor != nil {
			s.StatementBreak.StatementFor.BackPatchs = append(s.StatementBreak.StatementFor.BackPatchs, b)
		} else { // switch
			s.StatementBreak.StatementSwitch.BackPatchs = append(s.StatementBreak.StatementSwitch.BackPatchs, b)
		}
	case ast.STATEMENT_TYPE_RETURN:
		var haveVariableDefiniton bool
		for _, v := range s.StatementReturn.Defers {
			if v.Block.HaveVariableDefinition() {
				haveVariableDefiniton = true
				break
			}
		}
		ss := state
		if haveVariableDefiniton {
			ss = (&StackMapState{}).FromLast(state)
		}
		maxstack = m.buildReturnStatement(class, code, s.StatementReturn, context, ss)
		state.addTop(ss)
	case ast.STATEMENT_TYPE_SWITCH:
		s.StatementSwitch.BackPatchs = []*cg.JumpBackPatch{} //could compile multi times
		maxstack = m.buildSwitchStatement(class, code, s.StatementSwitch, context, state)
		if len(s.StatementSwitch.BackPatchs) > 0 {
			if code.CodeLength == context.LastStackMapOffset {
				code.Codes[code.CodeLength] = cg.OP_nop
				code.CodeLength++
			}
			backPatchEs(s.StatementSwitch.BackPatchs, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_GOTO:
		if s.StatementGoto.StatementLable.CodeOffsetGenerated {
			jumpTo(cg.OP_goto, code, s.StatementGoto.StatementLable.CodeOffset)
		} else {
			b := (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)
			s.StatementGoto.StatementLable.BackPatches = append(s.StatementGoto.StatementLable.BackPatches, b)
		}
	case ast.STATEMENT_TYPE_LABLE:
		s.StatmentLable.CodeOffsetGenerated = true
		s.StatmentLable.CodeOffset = code.CodeLength
		s.StatmentLable.BackPatches = []*cg.JumpBackPatch{} //could compile multi times
		if len(s.StatmentLable.BackPatches) > 0 {
			backPatchEs(s.StatmentLable.BackPatches, code.CodeLength) // back patch
		}
		context.MakeStackMap(code, state, code.CodeLength)
	case ast.STATEMENT_TYPE_DEFER: // nothing to do  ,defer will do after block is compiled
		s.Defer.StartPc = code.CodeLength
		s.Defer.StackMapState = (&StackMapState{}).FromLast(state)
	case ast.STATEMENT_TYPE_CLASS:
		s.Class.Name = m.newClassName(s.Class.Name)
		c := m.buildClass(s.Class)
		m.putClass(c.Name, c)
	}
	return
}
func (m *MakeClass) buildDefers(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context, ds []*ast.Defer, state *StackMapState) {
	index := len(ds) - 1
	for index >= 0 {
		var ss *StackMapState
		if ds[index].Block.HaveVariableDefinition() {
			ss = (&StackMapState{}).FromLast(state)
		} else {
			ss = state
		}
		m.buildBlock(class, code, &ds[index].Block, context, ss)
		index--
		state.addTop(ss)
	}
}
