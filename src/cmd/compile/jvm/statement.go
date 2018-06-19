package jvm

import (
	//"fmt"
	//"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeClass *MakeClass) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, s *ast.Statement,
	context *Context, state *StackMapState) (maxStack uint16) {
	//fmt.Println(s.GetPos)
	switch s.Type {
	case ast.STATEMENT_TYPE_EXPRESSION:
		if s.Expression.Type == ast.EXPRESSION_TYPE_FUNCTION {
			return makeClass.buildFunctionExpression(class, code, s.Expression, context, state)
		}
		maxStack, _ = makeClass.makeExpression.build(class, code, s.Expression, context, state)
	case ast.STATEMENT_TYPE_IF:
		s.StatementIf.Exits = []*cg.Exit{} //could compile multi times
		maxStack = makeClass.buildIfStatement(class, code, s.StatementIf, context, state)
		if len(s.StatementIf.Exits) > 0 {
			backfillExit(s.StatementIf.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_BLOCK: //new
		var ss *StackMapState
		if s.Block.HaveVariableDefinition() {
			ss = (&StackMapState{}).FromLast(state)
		} else {
			ss = state
		}
		makeClass.buildBlock(class, code, s.Block, context, ss)
		state.addTop(ss)
	case ast.STATEMENT_TYPE_FOR:
		s.StatementFor.Exits = []*cg.Exit{} //could compile multi times
		maxStack = makeClass.buildForStatement(class, code, s.StatementFor, context, state)
		if len(s.StatementFor.Exits) > 0 {
			backfillExit(s.StatementFor.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_CONTINUE:
		makeClass.buildDefers(class, code, context, s.StatementContinue.Defers, state)
		jumpTo(cg.OP_goto, code, s.StatementContinue.StatementFor.ContinueOPOffset)
	case ast.STATEMENT_TYPE_BREAK:
		makeClass.buildDefers(class, code, context, s.StatementBreak.Defers, state)
		b := (&cg.Exit{}).FromCode(cg.OP_goto, code)
		if s.StatementBreak.StatementFor != nil {
			s.StatementBreak.StatementFor.Exits = append(s.StatementBreak.StatementFor.Exits, b)
		} else { // switch
			s.StatementBreak.StatementSwitch.Exits = append(s.StatementBreak.StatementSwitch.Exits, b)
		}
	case ast.STATEMENT_TYPE_RETURN:
		maxStack = makeClass.buildReturnStatement(class, code, s.StatementReturn, context, state)
	case ast.STATEMENT_TYPE_SWITCH:
		s.StatementSwitch.Exits = []*cg.Exit{} //could compile multi times
		maxStack = makeClass.buildSwitchStatement(class, code, s.StatementSwitch, context, state)
		if len(s.StatementSwitch.Exits) > 0 {
			if code.CodeLength == context.lastStackMapOffset {
				code.Codes[code.CodeLength] = cg.OP_nop
				code.CodeLength++
			}
			backfillExit(s.StatementSwitch.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_GOTO:
		if s.StatementGoTo.StatementLabel.CodeOffsetGenerated {
			jumpTo(cg.OP_goto, code, s.StatementGoTo.StatementLabel.CodeOffset)
		} else {
			b := (&cg.Exit{}).FromCode(cg.OP_goto, code)
			s.StatementGoTo.StatementLabel.Exits = append(s.StatementGoTo.StatementLabel.Exits, b)
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
		s.Class.Name = makeClass.newClassName(s.Class.Name)
		c := makeClass.buildClass(s.Class)
		makeClass.putClass(c.Name, c)
	}
	return
}
func (makeClass *MakeClass) buildDefers(class *cg.ClassHighLevel,
	code *cg.AttributeCode, context *Context, ds []*ast.StatementDefer, state *StackMapState) {
	index := len(ds) - 1
	for index >= 0 {
		var ss *StackMapState
		if ds[index].Block.HaveVariableDefinition() {
			ss = (&StackMapState{}).FromLast(state)
		} else {
			ss = state
		}
		makeClass.buildBlock(class, code, &ds[index].Block, context, ss)
		index--
		state.addTop(ss)
	}
}
