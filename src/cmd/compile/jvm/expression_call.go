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
	if call.Func.IsClosureFunction() == false {
		maxstack = m.buildCallArgs(class, code, call.Args, call.Func.Typ.ParameterList, context)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Func.ClassMethod.Class.Name,
			Name:       call.Func.Name,
			Descriptor: call.Func.MkDescriptor(),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		//closure function call
	}
	if e.IsStatementExpression && len(e.VariableTypes) > 0 {
		if len(e.VariableTypes) == 1 {
			if 2 == e.VariableTypes[0].JvmSlotSize() {
				code.Codes[code.CodeLength] = cg.OP_pop2
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop
			}
			code.CodeLength++
		} else { // > 1
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}

func (m *MakeExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode, args []*ast.Expression, parameters ast.ParameterList, context *Context) (maxstack uint16) {
	currentStack := uint16(0)
	for k, e := range args {
		var variabletype *ast.VariableType
		if e.Typ == ast.EXPRESSION_TYPE_METHOD_CALL || e.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL && len(e.VariableTypes) > 1 {
			stack, _ := m.build(class, code, e, context)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			m.buildStoreArrayListAutoVar(code, context)
			for k, t := range e.VariableTypes {
				stack = m.unPackArraylist(class, code, k, t, context)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				if parameters[k].Typ.IsNumber() {
					if parameters[k].Typ.Typ != variabletype.Typ {
						m.numberTypeConverter(code, variabletype.Typ, parameters[k].Typ.Typ)
					}
				}
				currentStack += parameters[k].Typ.JvmSlotSize()
			}
			continue
		}
		variabletype = e.VariableType
		if e.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || e.Typ == ast.EXPRESSION_TYPE_METHOD_CALL {
			variabletype = e.VariableTypes[0]
		}
		ms, es := m.build(class, code, e, context)
		backPatchEs(es, code)
		if t := ms + currentStack; t > maxstack {
			maxstack = t
		}
		if parameters[k].Typ.IsNumber() {
			if parameters[k].Typ.Typ != variabletype.Typ {
				m.numberTypeConverter(code, variabletype.Typ, parameters[k].Typ.Typ)
			}
		}
		currentStack += parameters[k].Typ.JvmSlotSize()
	}
	return
}

func (m *MakeExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	if call.Method.IsStatic() {
		maxstack = m.buildCallArgs(class, code, call.Args, nil, context)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
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
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      call.Method.Func.ClassMethod.Class.Name,
		Name:       call.Name,
		Descriptor: call.Method.Func.Descriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
