package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildStrCat(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.ExpressionBinary, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2 // current stack is 2
	currenStack := maxstack
	stack, es := m.build(class, code, e.Left, context)
	backPatchEs(es, code.CodeLength)
	if t := currenStack + stack; t > maxstack {
		maxstack = t
	}
	m.stackTop2String(class, code, e.Left.GetTheOnlyOneVariableType())
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       special_method_init,
		Descriptor: "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	currenStack = 1
	stack, es = m.build(class, code, e.Right, context)
	backPatchEs(es, code.CodeLength)
	if t := currenStack + stack; t > maxstack {
		maxstack = t
	}
	m.stackTop2String(class, code, e.Right.GetTheOnlyOneVariableType())
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       `append`,
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       `toString`,
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return maxstack
}
