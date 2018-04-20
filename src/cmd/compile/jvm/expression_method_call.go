package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.Expression.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY {
		return m.buildArrayMethodCall(class, code, e, context, state)
	}
	if call.Expression.VariableType.Typ == ast.VARIABLE_TYPE_MAP {
		return m.buildMapMethodCall(class, code, e, context, state)
	}
	if call.Expression.VariableType.Typ == ast.VARIABLE_TYPE_JAVA_ARRAY {
		return m.buildJavaArrayMethodCall(class, code, e, context, state)
	}
	d := call.Method.Func.Descriptor
	if call.Class.LoadFromOutSide == false {
		d = Descriptor.methodDescriptor(call.Method.Func)
	}
	if call.Method.IsStatic() {
		maxstack = m.buildCallArgs(class, code, call.Args, call.Method.Func.Typ.ParameterList, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Method.Func.ClassMethod.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	maxstack, _ = m.build(class, code, call.Expression, context, nil)
	// object ref
	state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, call.Expression.VariableType)...)
	defer state.popStack(1)
	stack := m.buildCallArgs(class, code, call.Args, call.Method.Func.Typ.ParameterList, context, state)
	if t := stack + 1; t > maxstack {
		maxstack = t
	}
	if call.Class.IsInterface() {
		code.Codes[code.CodeLength] = cg.OP_invokeinterface
		class.InsertInterfaceMethodrefConst(cg.CONSTANT_InterfaceMethodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = interfaceMethodArgsCount(call.Method.Func.Typ)
		code.Codes[code.CodeLength+4] = 0
		code.CodeLength += 5
	} else {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Class.Name,
			Method:     call.Name,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	return
}
