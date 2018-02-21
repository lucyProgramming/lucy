package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context) (maxstack uint16) {
	switch call.Func.Name {
	case "print":
		return m.mkBuildinPrint(class, code, call, context)
	case "panic":
		return m.mkBuildinPanic(class, code, call, context)
	case "catch":
		return m.mkBuildinRecover(class, code, call, context)
	default:
		panic("unhandle buildin function:" + call.Func.Name)
	}
}
