package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode, args []*ast.Expression, parameters ast.ParameterList, context *Context) (maxstack uint16) {
	currentStack := uint16(0)
	for k, e := range args {
		var variabletype *ast.VariableType
		if e.IsCall() && len(e.VariableTypes) > 1 {
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
				if parameters[k].Typ.IsNumber() && parameters[k].Typ.Typ != variabletype.Typ {
					m.numberTypeConverter(code, variabletype.Typ, parameters[k].Typ.Typ)
				}
				currentStack += parameters[k].Typ.JvmSlotSize()
			}
			continue
		}
		variabletype = e.VariableType
		if e.IsCall() {
			variabletype = e.VariableTypes[0]
		}
		stack, es := m.build(class, code, e, context)
		backPatchEs(es, code.CodeLength)
		if t := stack + currentStack; t > maxstack {
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
