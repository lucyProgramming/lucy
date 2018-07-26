package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode,
	args []*ast.Expression, vArgs *ast.CallVArgs, context *Context, state *StackMapState) (maxStack uint16) {
	currentStack := uint16(0)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength) // let`s pop
	}()
	for _, e := range args {
		if e.MayHaveMultiValue() && len(e.MultiValues) > 1 {
			stack := buildExpression.build(class, code, e, context, state)
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
				state.pushStack(class, t)
			}
			continue
		}
		stack := buildExpression.build(class, code, e, context, state)
		if t := stack + currentStack; t > maxStack {
			maxStack = t
		}
		currentStack += jvmSlotSize(e.Value)
		state.pushStack(class, e.Value)
	}
	if vArgs == nil {
		return
	}
	if vArgs.NoArgs {
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		if t := 1 + currentStack; t > maxStack {
			maxStack = t
		}
	} else {
		if vArgs.IsJavaArray {
			stack := buildExpression.build(class, code, vArgs.Expressions[0], context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
		} else {
			loadInt32(class, code, int32(vArgs.Length))
			newArrayBaseOnType(class, code, vArgs.Type.Array)
			state.pushStack(class, vArgs.Type)
			currentStack++
			op := storeArrayElementByTypeOps(vArgs.Type.Array.Type)
			index := int32(0)
			for _, e := range vArgs.Expressions {
				if e.MayHaveMultiValue() && len(e.MultiValues) > 1 {
					stack := buildExpression.build(class, code, e, context, state)
					if t := stack + currentStack; t > maxStack {
						maxStack = t
					}
					multiValuePacker.storeMultiValueAutoVar(code, context)
					for kk, tt := range e.MultiValues {
						code.Codes[code.CodeLength] = cg.OP_dup
						code.CodeLength++
						loadInt32(class, code, index)
						currentStack += 2
						if t := currentStack + multiValuePacker.unPack(class, code, kk, tt, context); t > maxStack {
							maxStack = t
						}
						code.Codes[code.CodeLength] = op
						code.CodeLength++
						currentStack -= 2
						index++
					}
					continue
				}
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				state.pushStack(class, vArgs.Type)
				loadInt32(class, code, index)
				state.pushStack(class, &ast.Type{
					Type: ast.VariableTypeInt,
				})
				currentStack += 2
				stack := buildExpression.build(class, code, e, context, state)
				state.pushStack(class, e.Value)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				code.Codes[code.CodeLength] = op
				code.CodeLength++
				state.popStack(3)
				currentStack -= 2
				index++
			}
		}
	}

	return
}
