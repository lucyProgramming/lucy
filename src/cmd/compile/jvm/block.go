package jvm

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context, state *StackMapState) {
	willNotExecuteToEnd := false
	for _, s := range b.Statements {
		if willNotExecuteToEnd == true && s.Type == ast.StatementTypeLabel {
			jumpForwards := len(s.StatementLabel.Exits) > 0 // jump forward
			willNotExecuteToEnd = !jumpForwards
			//continue compile block from this label statement
		}
		if willNotExecuteToEnd {
			continue
		}
		maxStack := buildPackage.buildStatement(class, code, b, s, context, state)
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
			buildPackage.mkNonStaticFieldDefaultValue(class, code, context, state)
		}
		//unCondition goto
		if buildPackage.statementIsUnConditionGoTo(s) {
			willNotExecuteToEnd = true
			continue
		}
		//block deadEnd
		if s.Type == ast.StatementTypeBlock {
			willNotExecuteToEnd = s.Block.WillNotExecuteToEnd
			continue
		}
		if s.Type == ast.StatementTypeIf && s.StatementIf.ElseBlock != nil {
			t := s.StatementIf.Block.WillNotExecuteToEnd
			for _, v := range s.StatementIf.ElseIfList {
				t = t && v.Block.WillNotExecuteToEnd
			}
			t = t && s.StatementIf.ElseBlock.WillNotExecuteToEnd
			willNotExecuteToEnd = t
			continue
		}
		if s.Type == ast.StatementTypeSwitch && s.StatementSwitch.Default != nil {
			t := s.StatementSwitch.Default.WillNotExecuteToEnd
			for _, v := range s.StatementSwitch.StatementSwitchCases {
				if v.Block != nil {
					t = t && v.Block.WillNotExecuteToEnd
				} else {
					//this will fallthrough
					t = false
					break
				}
			}
			t = t && s.StatementSwitch.Default.WillNotExecuteToEnd
			willNotExecuteToEnd = t
			continue
		}
	}
	// if b.IsFunctionTopBlock == true must a return at end
	if b.IsFunctionBlock == false && len(b.Defers) > 0 {
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		buildPackage.buildDefers(class, code, context, b.Defers, state)
	}
	b.WillNotExecuteToEnd = willNotExecuteToEnd
	return
}

func (buildPackage *BuildPackage) statementIsUnConditionGoTo(s *ast.Statement) bool {
	return s.Type == ast.StatementTypeReturn ||
		s.Type == ast.StatementTypeGoto ||
		s.Type == ast.StatementTypeContinue ||
		s.Type == ast.StatementTypeBreak
}
