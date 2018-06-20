package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildJavaArrayMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	maxStack, _ = makeExpression.build(class, code, call.Expression, context, state)
	switch call.Name {
	case common.ARRAY_METHOD_SIZE:
		code.Codes[code.CodeLength] = cg.OP_arraylength
		code.CodeLength++
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}

func (makeExpression *MakeExpression) buildArrayMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length) // ref type
	}()
	call := e.Data.(*ast.ExpressionMethodCall)
	maxStack, _ = makeExpression.build(class, code, call.Expression, context, state)
	state.pushStack(class, call.Expression.ExpressionValue)
	switch call.Name {
	case common.ARRAY_METHOD_CAP,
		common.ARRAY_METHOD_SIZE,
		common.ARRAY_METHOD_START,
		common.ARRAY_METHOD_END:
		meta := ArrayMetas[call.Expression.ExpressionValue.ArrayType.Type]
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
	case common.ARRAY_METHOD_APPEND:
		meta := ArrayMetas[call.Expression.ExpressionValue.ArrayType.Type]
		appendName := "append"
		appendDescriptor := meta.appendDescriptor
		for _, v := range call.Args {
			currentStack := uint16(1)
			if v.MayHaveMultiValue() && len(v.ExpressionMultiValues) > 0 {
				stack, _ := makeExpression.build(class, code, v, context, nil)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				multiValuePacker.storeArrayListAutoVar(code, context)
				for kk, t := range v.ExpressionMultiValues {
					currentStack = 1
					if t := multiValuePacker.unPack(class, code, kk, t, context) + currentStack; t > maxStack {
						maxStack = t
					}
					if t := currentStack + jvmSize(t); t > maxStack {
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
			stack, es := makeExpression.build(class, code, v, context, state)
			if len(es) > 0 {
				fillOffsetForExits(es, code.CodeLength)
				state.pushStack(class, v.ExpressionValue)
				context.MakeStackMap(code, state, code.CodeLength)
				state.popStack(1) // must be a logical expression
			}
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
	case common.ARRAY_METHOD_APPEND_ALL:
		meta := ArrayMetas[call.Expression.ExpressionValue.ArrayType.Type]
		for _, v := range call.Args {
			currentStack := uint16(1)
			appendName := "append"
			appendDescriptor := meta.appendAllDescriptor
			if v.MayHaveMultiValue() && len(v.ExpressionMultiValues) > 1 {
				stack, _ := makeExpression.build(class, code, v, context, state)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				multiValuePacker.storeArrayListAutoVar(code, context)
				for kk, tt := range v.ExpressionMultiValues {
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
			stack, _ := makeExpression.build(class, code, v, context, state)
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
