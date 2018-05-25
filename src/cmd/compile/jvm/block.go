package jvm

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context, state *StackMapState) {
	var deadend bool = false
	for _, s := range b.Statements {
		if deadend == true && s.Typ == ast.STATEMENT_TYPE_LABLE {
			jumpForwards := len(s.StatmentLable.BackPatches) > 0 // jump forward
			deadend = !jumpForwards
			//continue compile block from this lable statment
		}
		if deadend {
			continue
		}
		maxstack := m.buildStatement(class, code, b, s, context, state)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		if len(state.Stacks) > 0 {
			for _, v := range state.Stacks {
				fmt.Println(v.Verify)
			}
			panic(fmt.Sprintf("stack is not empty:%d", len(state.Stacks)))
		}
		if s.IsCallFatherContructionStatement { // special case
			state.Locals[0] = state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(class.Name))
			m.mkFieldDefaultValue(class, code, context, state)
		}
		//uncondition goto
		if m.statementIsUnConditionGoto(s) {
			deadend = true
			continue
		}
		//block deadend
		if s.Typ == ast.STATEMENT_TYPE_BLOCK {
			deadend = s.Block.DeadEnding
			continue
		}
		if s.Typ == ast.STATEMENT_TYPE_IF && s.StatementIf.ElseBlock != nil {
			t := s.StatementIf.Block.DeadEnding
			for _, v := range s.StatementIf.ElseIfList {
				t = t && v.Block.DeadEnding
			}
			t = t && s.StatementIf.ElseBlock.DeadEnding
			deadend = t
			continue
		}
		if s.Typ == ast.STATEMENT_TYPE_SWITCH && s.StatementSwitch.Default != nil {
			t := s.StatementSwitch.Default.DeadEnding
			for _, v := range s.StatementSwitch.StatmentSwitchCases {
				if v.Block != nil {
					t = t && v.Block.DeadEnding
				} else {
					//this will fallthrough
					t = t && false
					break
				}
			}
			t = t && s.StatementIf.ElseBlock.DeadEnding
			deadend = t
			continue
		}
	}
	// if b.IsFunctionTopBlock == true must a return at end
	if b.IsFunctionTopBlock == false && len(b.Defers) > 0 {
		m.buildDefers(class, code, context, b.Defers, state)
	}
	b.DeadEnding = deadend
	return
}

func (m *MakeClass) statementIsUnConditionGoto(s *ast.Statement) bool {
	return s.Typ == ast.STATEMENT_TYPE_RETURN ||
		s.Typ == ast.STATEMENT_TYPE_GOTO ||
		s.Typ == ast.STATEMENT_TYPE_CONTINUE ||
		s.Typ == ast.STATEMENT_TYPE_BREAK
}
