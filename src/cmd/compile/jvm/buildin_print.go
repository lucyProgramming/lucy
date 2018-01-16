package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

/*
	function print
*/
func (m *MakeExpression) mkBuildinPrint(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall) (maxstack uint16) {
	maxstack = 1
	code.Codes[code.CodeLength] = cg.OP_getstatic
	class.InsertFieldRef(cg.CONSTANT_Fieldref_info_high_level{
		Class:       "java/lang/System",
		NameAndType: "out:Ljava/io/PrintStream;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	maxstack += 2
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:       "java/lang/System",
		NameAndType: `"<init>":()V`,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3

	return
}
func (m *MakeExpression) stackTop2String(class *cg.ClassHighLevel, code *cg.AttributeCode, typ *ast.VariableType) {
	switch typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_ifne
		binary.BigEndian.PutUint16()
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_CHAR:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		fallthrough
	case ast.VARIABLE_TYPE_LONG:

	case ast.VARIABLE_TYPE_FLOAT:
	case ast.VARIABLE_TYPE_DOUBLE:
	case ast.VARIABLE_TYPE_STRING:
		return
	case ast.VARIABLE_TYPE_OBJECT:

	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
	case ast.VARIABLE_TYPE_FUNCTION:
	}
}
