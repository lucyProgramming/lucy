package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildExpression) buildStrCat(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	bin := e.Data.(*ast.ExpressionBinary)
	maxStack = this.build(class, code, bin.Left, context, state)
	if t := this.stackTop2String(class, code, bin.Left.Value, context, state); t > maxStack {
		maxStack = t
	}
	state.pushStack(class, state.newObjectVariableType(javaStringClass))
	stack := this.build(class, code, bin.Right, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	if t := 1 + this.stackTop2String(class, code,
		bin.Right.Value, context, state); t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      javaStringClass,
		Method:     `concat`,
		Descriptor: "(Ljava/lang/String;)Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return

}
