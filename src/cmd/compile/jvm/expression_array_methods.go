package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildJavaArrayMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	maxStack = buildExpression.build(class, code, call.Expression, context, state)
	switch call.Name {
	case common.ArrayMethodSize:
		code.Codes[code.CodeLength] = cg.OP_arraylength
		code.CodeLength++
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}

func (buildExpression *BuildExpression) buildArrayMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length) // ref type
	}()
	call := e.Data.(*ast.ExpressionMethodCall)
	maxStack = buildExpression.build(class, code, call.Expression, context, state)
	state.pushStack(class, call.Expression.Value)
	switch call.Name {
	case common.ArrayMethodCap,
		common.ArrayMethodSize,
		common.ArrayMethodStart,
		common.ArrayMethodEnd:
		meta := ArrayMetas[call.Expression.Value.Array.Type]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     call.Name,
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	case common.ArrayMethodAppend:
		meta := ArrayMetas[call.Expression.Value.Array.Type]
		appendName := "append"
		appendDescriptor := meta.appendDescriptor
		for _, v := range call.Args {
			currentStack := uint16(1)
			if v.MayHaveMultiValue() && len(v.MultiValues) > 0 {
				stack := buildExpression.build(class, code, v, context, state)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				multiValuePacker.storeMultiValueAutoVar(code, context)
				for kk, t := range v.MultiValues {
					currentStack = 1
					if t := multiValuePacker.unPack(class, code, kk, t, context) + currentStack; t > maxStack {
						maxStack = t
					}
					if t := currentStack + jvmSlotSize(t); t > maxStack {
						maxStack = t
					}
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      meta.className,
						Method:     appendName,
						Descriptor: meta.appendDescriptor,
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				}
				continue
			}
			stack := buildExpression.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxStack {
				maxStack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.className,
				Method:     appendName,
				Descriptor: appendDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	case common.ArrayMethodAppendAll:
		meta := ArrayMetas[call.Expression.Value.Array.Type]
		for _, v := range call.Args {
			currentStack := uint16(1)
			appendName := "append"
			appendDescriptor := meta.appendAllDescriptor
			if v.MayHaveMultiValue() && len(v.MultiValues) > 1 {
				stack := buildExpression.build(class, code, v, context, state)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				multiValuePacker.storeMultiValueAutoVar(code, context)
				for kk, tt := range v.MultiValues {
					currentStack := uint16(1)
					stack = multiValuePacker.unPack(class, code, kk, tt, context)
					if t := currentStack + 2; t > maxStack {
						maxStack = t
					}
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      meta.className,
						Method:     appendName,
						Descriptor: appendDescriptor,
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				}
				continue
			}
			stack := buildExpression.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxStack {
				maxStack = t
			}
			//get elements field
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.className,
				Method:     appendName,
				Descriptor: appendDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}
