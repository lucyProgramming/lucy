package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Func.Isbuildin {
		return m.mkBuildinFunctionCall(class, code, call, context)
	}
	return
}

func (m *MakeExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode, args []*ast.Expression, context *Context) (maxstack uint16) {
	stack := uint16(0)
	for _, e := range args {
		ms, es := m.build(class, code, e, context)
		backPatchEs(es, code)
		maxstack = ms + stack
	}
	return
}

func (m *MakeExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	return
}
