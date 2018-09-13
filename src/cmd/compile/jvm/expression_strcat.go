package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildStrCat(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.ExpressionBinary,
	context *Context, state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if e.Left.IsString("") || e.Right.IsString("") {
		not := e.Left
		if e.Right.IsString("") == false {
			not = e.Right
		}
		maxStack = buildExpression.build(class, code, not, context, state)
		if t := buildExpression.stackTop2String(class, code, not.Value, context, state); t > maxStack {
			maxStack = t
		}
		return
	} else {
		maxStack = buildExpression.build(class, code, e.Left, context, state)
		if t := buildExpression.stackTop2String(class, code, e.Left.Value, context, state); t > maxStack {
			maxStack = t
		}
		state.pushStack(class, state.newObjectVariableType(javaStringClass))
		stack := buildExpression.build(class, code, e.Right, context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		if t := 1 + buildExpression.stackTop2String(class, code, e.Right.Value, context, state); t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     `concat`,
			Descriptor: "(Ljava/lang/String;)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
}
