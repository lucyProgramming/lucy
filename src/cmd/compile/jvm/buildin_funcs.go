package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) mkBuildinFunctionCall(call *ast.ExpressionFunctionCall, code cg.AttributeCode) {
	switch call.Func.Name {
	case "print":
		m.mkBuildinPrint(call, code)
	default:
		panic("unhandle buildin function" + call.Func.Name)
	}
}
