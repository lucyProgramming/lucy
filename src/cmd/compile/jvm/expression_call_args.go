package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode, args []*ast.Expression, parameters ast.ParameterList, context *Context, state *StackMapState) (maxstack uint16) {
	currentStack := uint16(0)
	for k, e := range args {
		var variabletype *ast.VariableType
		if e.MayHaveMultiValue() && len(e.VariableTypes) > 1 {
			stack, _ := m.build(class, code, e, context, nil)
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
		if e.MayHaveMultiValue() {
			variabletype = e.VariableTypes[0]
		}
		stack, es := m.build(class, code, e, context, state)
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
