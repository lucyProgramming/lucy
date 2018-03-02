package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

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
		backPatchEs(es, code.CodeLength)
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
