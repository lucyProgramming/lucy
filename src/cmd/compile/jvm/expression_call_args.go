package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode,
	args []*ast.Expression, parameters ast.ParameterList, context *Context, state *StackMapState) (maxStack uint16) {
	currentStack := uint16(0)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength) // let`s pop
	}()
	parameterIndex := 0
	for _, e := range args {
		if e.MayHaveMultiValue() && len(e.MultiValues) > 1 {
			stack, _ := buildExpression.build(class, code, e, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			multiValuePacker.storeMultiValueAutoVar(code, context)
			for k, t := range e.MultiValues {
				stack = multiValuePacker.unPack(class, code, k, t, context)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				currentStack += jvmSlotSize(t)
				state.pushStack(class, parameters[parameterIndex].Type)
				parameterIndex++
			}
			continue
		}
		stack, es := buildExpression.build(class, code, e, context, state)
		if len(es) > 0 {
			state.pushStack(class, &ast.Type{Type: ast.VariableTypeBool})
			writeExits(es, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := stack + currentStack; t > maxStack {
			maxStack = t
		}
		currentStack += jvmSlotSize(parameters[parameterIndex].Type)
		state.pushStack(class, parameters[parameterIndex].Type)
		parameterIndex++
	}
	return
}
