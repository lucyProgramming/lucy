package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildSwitchStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementSwitch, context *Context) (maxstack uint16) {
	maxstack, _ = m.MakeExpression.build(class, code, s.Condition, context)
	for _, c := range s.StatmentSwitchCases {
		for _, ee := range c.Matches {
			if ee.IsCall() && len(ee.VariableTypes) > 0 {
				continue
			}
		}
	}
	return
}
