package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall) (maxstack uint16) {
	switch call.Func.Name {
	case "print":
		return m.mkBuildinPrint(class, code, call)
	case "panic":
		return m.mkBuildinPanic(class, code, call)
	case "recover":
		return m.mkBuildinRecover(class, code, call)
	default:
		panic("unhandle buildin function" + call.Func.Name)
	}
}
