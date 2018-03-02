package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.Expression.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY_INSTANCE {
		return m.buildArrayMethodCall(class, code, e, context)
	}
	if call.Expression.VariableType.Typ == ast.VARIABLE_TYPE_MAP {
		return m.buildMapMethodCall(class, code, e, context)
	}
	if call.Method.IsStatic() {
		maxstack = m.buildCallArgs(class, code, call.Args, nil, context)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Method.Func.ClassMethod.Class.Name,
			Name:       call.Name,
			Descriptor: call.Method.Func.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	maxstack, _ = m.build(class, code, call.Expression, context)
	stack := m.buildCallArgs(class, code, call.Args, nil, context)
	if t := stack + 1; t > maxstack {
		maxstack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      call.Method.Func.ClassMethod.Class.Name,
		Name:       call.Name,
		Descriptor: call.Method.Func.Descriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
