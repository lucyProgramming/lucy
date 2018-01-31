package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

/*
	function print
*/
func (m *MakeExpression) mkBuildinPrint(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_getstatic
	class.InsertFieldRef(cg.CONSTANT_Fieldref_info_high_level{
		Class:      "java/lang/System",
		Name:       "out",
		Descriptor: "Ljava/io/PrintStream;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       `<init>`,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 3
	stack := uint16(2)
	for _, v := range call.Args {
		var variableType *ast.VariableType
		if v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL {
			if len(v.VariableTypes) > 1 {
				panic(111)
			} else {
				variableType = v.VariableTypes[0]
			}
		} else {
			variableType = v.VariableType
		}
		s, es := m.build(class, code, v, context)
		backPatchEs(es, code)
		if stack+s > maxstack {
			maxstack = stack + s
		}
		stack += v.VariableType.JvmSlotSize()
		m.stackTop2String(class, code, variableType)
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/StringBuilder",
			Name:       "append",
			Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	// append crlf
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst("\n", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       "append",
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// tostring
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       "toString",
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/io/PrintStream",
		Name:       "println",
		Descriptor: "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
