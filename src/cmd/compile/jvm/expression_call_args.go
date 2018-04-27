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
		if e.MayHaveMultiValue() && len(e.Values) > 1 {
			stack, _ := m.build(class, code, e, context, state)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			arrayListPacker.buildStoreArrayListAutoVar(code, context)
			for k, t := range e.Values {
				stack = arrayListPacker.unPack(class, code, k, t, context)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				if parameters[parameterIndex].Typ.IsNumber() && parameters[parameterIndex].Typ.Typ != t.Typ {
					m.numberTypeConverter(code, t.Typ, parameters[k].Typ.Typ)
				}
				currentStack += jvmSize(t)
				state.pushStack(class, parameters[parameterIndex].Typ)
				parameterIndex++
			}
			continue
		}
		variableType := e.Value
		if e.MayHaveMultiValue() {
			variableType = e.Values[0]
		}
		stack, es := m.build(class, code, e, context, state)
		if len(es) > 0 {
			state.pushStack(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_BOOL})
			backPatchEs(es, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		if parameters[parameterIndex].Typ.IsNumber() {
			if parameters[parameterIndex].Typ.Typ != variableType.Typ {
				m.numberTypeConverter(code, variableType.Typ, parameters[parameterIndex].Typ.Typ)
			}
		}
		currentStack += jvmSize(parameters[parameterIndex].Typ)
		state.pushStack(class, parameters[parameterIndex].Typ)
		parameterIndex++
	}
	return
}
