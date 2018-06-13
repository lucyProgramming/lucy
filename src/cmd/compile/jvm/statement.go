package jvm

import (
	//"fmt"
	//"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, s *ast.Statement,
	context *Context, state *StackMapState) (maxStack uint16) {
	//fmt.Println(s.Pos)
	switch s.Typ {
	case ast.STATEMENT_TYPE_EXPRESSION:
		if s.Expression.Typ == ast.EXPRESSION_TYPE_FUNCTION {
			return m.buildFunctionExpression(class, code, s.Expression, context, state)
		}
		maxStack, _ = m.MakeExpression.build(class, code, s.Expression, context, state)
	case ast.STATEMENT_TYPE_IF:
		s.StatementIf.BackPatchs = []*cg.Exit{} //could compile multi times
		maxStack = m.buildIfStatement(class, code, s.StatementIf, context, state)
		if len(s.StatementIf.BackPatchs) > 0 {
			backfillExit(s.StatementIf.BackPatchs, code.CodeLength)
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
		s.StatementFor.Exits = []*cg.Exit{} //could compile multi times
		maxStack = m.buildForStatement(class, code, s.StatementFor, context, state)
		if len(s.StatementFor.Exits) > 0 {
			backfillExit(s.StatementFor.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_CONTINUE:
		m.buildDefers(class, code, context, s.StatementContinue.Defers, state)
		jumpTo(cg.OP_goto, code, s.StatementContinue.StatementFor.ContinueOPOffset)
	case ast.STATEMENT_TYPE_BREAK:
		m.buildDefers(class, code, context, s.StatementBreak.Defers, state)
		b := (&cg.Exit{}).FromCode(cg.OP_goto, code)
		if s.StatementBreak.StatementFor != nil {
			s.StatementBreak.StatementFor.Exits = append(s.StatementBreak.StatementFor.Exits, b)
		} else { // switch
			s.StatementBreak.StatementSwitch.Exits = append(s.StatementBreak.StatementSwitch.Exits, b)
		}
	case ast.STATEMENT_TYPE_RETURN:
		maxStack = m.buildReturnStatement(class, code, s.StatementReturn, context, state)
	case ast.STATEMENT_TYPE_SWITCH:
		s.StatementSwitch.Exits = []*cg.Exit{} //could compile multi times
		maxStack = m.buildSwitchStatement(class, code, s.StatementSwitch, context, state)
		if len(s.StatementSwitch.Exits) > 0 {
			if code.CodeLength == context.LastStackMapOffset {
				code.Codes[code.CodeLength] = cg.OP_nop
				code.CodeLength++
			}
			backfillExit(s.StatementSwitch.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_GOTO:
		if s.StatementGoto.StatementLable.CodeOffsetGenerated {
			jumpTo(cg.OP_goto, code, s.StatementGoto.StatementLable.CodeOffset)
		} else {
			b := (&cg.Exit{}).FromCode(cg.OP_goto, code)
			s.StatementGoto.StatementLable.Exits = append(s.StatementGoto.StatementLable.Exits, b)
		}
	case ast.STATEMENT_TYPE_LABLE:
		s.StatementLabel.CodeOffsetGenerated = true
		s.StatementLabel.CodeOffset = code.CodeLength
		s.StatementLabel.Exits = []*cg.Exit{} //could compile multi times
		if len(s.StatementLabel.Exits) > 0 {
			backfillExit(s.StatementLabel.Exits, code.CodeLength) // back patch
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
func (m *MakeClass) buildDefers(class *cg.ClassHighLevel,
	code *cg.AttributeCode, context *Context, ds []*ast.Defer, state *StackMapState) {
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
