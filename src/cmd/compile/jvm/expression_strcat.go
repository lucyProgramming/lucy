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
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(javaStringBuilderClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.pushStack(class, state.newObjectVariableType(javaStringBuilderClass))
	maxStack = 2 // current stack is 2
	currentStack := uint16(1)
	stack, es := buildExpression.build(class, code, e.Left, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, e.Left.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	if t := currentStack +
		buildExpression.stackTop2String(class, code, e.Left.Value, context, state); t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     "append",
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	stack, es = buildExpression.build(class, code, e.Right, context, state)
	if len(es) > 0 {
		writeExits(es, code.CodeLength)
		state.pushStack(class, e.Right.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	if t := currentStack + buildExpression.stackTop2String(class, code, e.Right.Value, context, state); t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     "append",
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     `toString`,
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
