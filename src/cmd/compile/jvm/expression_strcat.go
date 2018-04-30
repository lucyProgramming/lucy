package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

//TODO:: stack maps
func (m *MakeExpression) buildStrCat(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.ExpressionBinary, context *Context, state *StackMapState) (maxstack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_string_builder_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_string_builder_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.pushStack(class, state.newObjectVariableType(java_string_builder_class))
	maxstack = 2 // current stack is 2
	currenStack := uint16(1)
	stack, es := m.build(class, code, e.Left, context, state)
	if len(es) > 0 {
		backPatchEs(es, code.CodeLength)
		state.pushStack(class, e.Left.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if t := currenStack + stack; t > maxstack {
		maxstack = t
	}
	m.stackTop2String(class, code, e.Left.Value, context, state)
	if t := currenStack + stack + jvmSize(e.Left.Value); t > maxstack {
		maxstack = t
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_string_builder_class,
		Method:     "append",
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	currenStack = 1
	stack, es = m.build(class, code, e.Right, context, state)
	if len(es) > 0 {
		backPatchEs(es, code.CodeLength)
		state.pushStack(class, e.Right.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	if t := currenStack + stack; t > maxstack {
		maxstack = t
	}
	m.stackTop2String(class, code, e.Right.Value, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_string_builder_class,
		Method:     "append",
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_string_builder_class,
		Method:     `toString`,
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return maxstack
}
