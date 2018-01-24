package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildStrCat(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.ExpressionBinary, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2 // current stack is 2
	stack, _ := m.build(class, code, e.Left, context)
	maxstack += stack
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       `<init>`,
		Descriptor: "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	stack, _ = m.build(class, code, e.Right, context)
	maxstack += stack
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       `append`,
		Descriptor: "(Ljava/lang/String;)java/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       `toString`,
		Descriptor: "()java/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	return maxstack
}
