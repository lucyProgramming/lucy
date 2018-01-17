package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinPanic(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context) (maxstack uint16) {
	return
}
