package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildMethodCallJavaOnArray(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
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

func (buildExpression *BuildExpression) buildMethodCallOnArray(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length) // ref type
	}()
	call := e.Data.(*ast.ExpressionMethodCall)
	maxStack = buildExpression.build(class, code, call.Expression, context, state)
	state.pushStack(class, call.Expression.Value)
	switch call.Name {
	case common.ArrayMethodSize,
		common.ArrayMethodStart,
		common.ArrayMethodCap,
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
		for k, v := range call.Args {
			currentStack := uint16(1)
			if k != len(call.Args)-1 {
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				state.pushStack(class, call.Expression.Value)
				currentStack++
			}
			stack := buildExpression.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxStack {
				maxStack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.className,
				Method:     "append",
				Descriptor: meta.appendDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			if k != len(call.Args)-1 {
				state.popStack(1)
			}
		}
	case common.ArrayMethodAppendAll:
		meta := ArrayMetas[call.Expression.Value.Array.Type]
		for k, v := range call.Args {
			currentStack := uint16(1)
			if k != len(call.Args)-1 {
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				state.pushStack(class, call.Expression.Value)
				currentStack++
			}
			stack := buildExpression.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxStack {
				maxStack = t
			}
			//get elements field
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.className,
				Method:     "append",
				Descriptor: meta.appendAllDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			if k != len(call.Args)-1 {
				state.popStack(1)
			}
		}
	case common.ArrayMethodGetUnderlyingArray:
		meta := ArrayMetas[call.Expression.Value.Array.Type]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "elements",
			Descriptor: meta.elementsFieldDescriptor,
		},
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if meta.elementsFieldDescriptor != Descriptor.typeDescriptor(e.Value) {
			typeConverter.castPointer(class, code, e.Value)
		}
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}

	return
}
