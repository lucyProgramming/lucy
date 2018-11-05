package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildExpression) buildStringSlice(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	slice := e.Data.(*ast.ExpressionSlice)
	maxStack = this.build(class, code, slice.ExpressionOn, context, state)
	state.pushStack(class, state.newObjectVariableType(javaStringClass))
	// build start
	stack := this.build(class, code, slice.Start, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	if slice.End != nil {
		state.pushStack(class, slice.Start.Value)
		stack = this.build(class, code, slice.End, context, state)
		if t := 2 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "substring",
			Descriptor: "(II)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "substring",
			Descriptor: "(I)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	return
}

func (this *BuildExpression) buildSlice(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	slice := e.Data.(*ast.ExpressionSlice)
	if slice.ExpressionOn.Value.Type == ast.VariableTypeString {
		return this.buildStringSlice(class, code, e, context, state)
	}
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	meta := ArrayMetas[e.Value.Array.Type]
	maxStack = this.build(class, code, slice.ExpressionOn, context, state)
	state.pushStack(class, slice.ExpressionOn.Value)
	if slice.End != nil {
		// build start
		stack := this.build(class, code, slice.Start, context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		state.pushStack(class, slice.Start.Value)
		stack = this.build(class, code, slice.End, context, state)
		if t := 3 + stack; t > maxStack {
			maxStack = t
		}
	} else {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      meta.className,
			Method:     "size",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		state.pushStack(class, slice.Start.Value)
		stack := this.build(class, code, slice.Start, context, state)
		if t := 2 + stack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      meta.className,
		Method:     "slice",
		Descriptor: meta.sliceDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
