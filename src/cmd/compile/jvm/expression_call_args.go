package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode,
	args []*ast.Expression, parameters ast.ParameterList, context *Context, state *StackMapState) (maxStack uint16) {
	currentStack := uint16(0)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength) // let`s pop
	}()
	parameterIndex := 0
	for _, e := range args {
		if e.MayHaveMultiValue() && len(e.ExpressionMultiValues) > 1 {
			stack, _ := makeExpression.build(class, code, e, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			multiValuePacker.storeArrayListAutoVar(code, context)
			for k, t := range e.ExpressionMultiValues {
				stack = multiValuePacker.unPack(class, code, k, t, context)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				currentStack += jvmSize(t)
				state.pushStack(class, parameters[parameterIndex].Type)
				parameterIndex++
			}
			continue
		}
		stack, es := makeExpression.build(class, code, e, context, state)
		if len(es) > 0 {
			state.pushStack(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_BOOL})
			fillOffsetForExits(es, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := stack + currentStack; t > maxStack {
			maxStack = t
		}
		currentStack += jvmSize(parameters[parameterIndex].Type)
		state.pushStack(class, parameters[parameterIndex].Type)
		parameterIndex++
	}
	return
}
