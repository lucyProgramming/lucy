package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode, args []*ast.Expression, parameters ast.ParameterList, context *Context, state *StackMapState) (maxstack uint16) {
	currentStack := uint16(0)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength) // let`s pop
	}()
	parameterIndex := 0
	for _, e := range args {
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
				if parameters[parameterIndex].Typ.IsNumber() && parameters[parameterIndex].Typ.Typ != t.Typ {
					m.numberTypeConverter(code, t.Typ, parameters[k].Typ.Typ)
				}
				currentStack += t.JvmSlotSize()
				state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, parameters[parameterIndex].Typ)...)
				parameterIndex++
			}
			continue
		}
		variableType := e.VariableType
		if e.MayHaveMultiValue() {
			variableType = e.VariableTypes[0]
		}
		stack, es := m.build(class, code, e, context, state)
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, parameters[parameterIndex].Typ)...)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(state, code.CodeLength))
		}
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		if parameters[parameterIndex].Typ.IsNumber() {
			if parameters[parameterIndex].Typ.Typ != variableType.Typ {
				m.numberTypeConverter(code, variableType.Typ, parameters[parameterIndex].Typ.Typ)
			}
		}
		currentStack += parameters[parameterIndex].Typ.JvmSlotSize()
		parameterIndex++
	}
	return
}
