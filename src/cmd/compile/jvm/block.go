package jvm

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context, state *StackMapState) {
	var deadEnd bool = false
	for _, s := range b.Statements {
		if deadEnd == true && s.Typ == ast.STATEMENT_TYPE_LABLE {
			jumpForwards := len(s.StatementLabel.Exits) > 0 // jump forward
			deadEnd = !jumpForwards
			//continue compile block from this label statement
		}
		if deadEnd {
			continue
		}
		maxStack := m.buildStatement(class, code, b, s, context, state)
		if maxStack > code.MaxStack {
			code.MaxStack = maxStack
		}
		if len(state.Stacks) > 0 {
			for _, v := range state.Stacks {
				fmt.Println(v.Verify)
			}
			panic(fmt.Sprintf("stack is not empty:%d", len(state.Stacks)))
		}
		if s.IsCallFatherConstructionStatement { // special case
			state.Locals[0] = state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(class.Name))
			m.mkFieldDefaultValue(class, code, context, state)
		}
		//unCondition goto
		if m.statementIsUnConditionGoto(s) {
			deadEnd = true
			continue
		}
		//block deadEnd
		if s.Typ == ast.STATEMENT_TYPE_BLOCK {
			deadEnd = s.Block.DeadEnding
			continue
		}
		if s.Typ == ast.STATEMENT_TYPE_IF && s.StatementIf.ElseBlock != nil {
			t := s.StatementIf.Block.DeadEnding
			for _, v := range s.StatementIf.ElseIfList {
				t = t && v.Block.DeadEnding
			}
			t = t && s.StatementIf.ElseBlock.DeadEnding
			deadEnd = t
			continue
		}
		if s.Typ == ast.STATEMENT_TYPE_SWITCH && s.StatementSwitch.Default != nil {
			t := s.StatementSwitch.Default.DeadEnding
			for _, v := range s.StatementSwitch.StatementSwitchCases {
				if v.Block != nil {
					t = t && v.Block.DeadEnding
				} else {
					//this will fallthrough
					t = t && false
					break
				}
			}
			t = t && s.StatementSwitch.Default.DeadEnding
			deadEnd = t
			continue
		}
	}
	// if b.IsFunctionTopBlock == true must a return at end
	if b.IsFunctionTopBlock == false && len(b.Defers) > 0 {
		m.buildDefers(class, code, context, b.Defers, state)
	}
	b.DeadEnding = deadEnd
	return
}

func (m *MakeClass) statementIsUnConditionGoto(s *ast.Statement) bool {
	return s.Typ == ast.STATEMENT_TYPE_RETURN ||
		s.Typ == ast.STATEMENT_TYPE_GOTO ||
		s.Typ == ast.STATEMENT_TYPE_CONTINUE ||
		s.Typ == ast.STATEMENT_TYPE_BREAK
}
