package jvm

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

/*
	function print
*/
func (m *MakeExpression) mkBuildinPrint(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_getstatic
	class.InsertFieldRef(cg.CONSTANT_Fieldref_info_high_level{
		Class: "java/lang/System",
		Name:  "out",
		Type:  "Ljava/io/PrintStream;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/lang/StringBuilder",
		Name:  `<init>`,
		Type:  "()V",
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
			Class: "java/lang/StringBuilder",
			Name:  "append",
			Type:  "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	// append crlf
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst("\n", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/lang/StringBuilder",
		Name:  "append",
		Type:  "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// tostring
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/lang/StringBuilder",
		Name:  "toString",
		Type:  "()Ljava/lang/String",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/io/PrintStream",
		Name:  "println",
		Type:  "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	return
}
func (m *MakeExpression) stackTop2String(class *cg.ClassHighLevel, code *cg.AttributeCode, typ *ast.VariableType) {
	switch typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class: "java/lang/String",
			Name:  "valueOf",
			Type:  "(Z)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_CHAR:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class: "java/lang/String",
			Name:  "valueOf",
			Type:  "(I)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class: "java/lang/String",
			Name:  "valueOf",
			Type:  "(J)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class: "java/lang/String",
			Name:  "valueOf",
			Type:  "(F)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class: "java/lang/String",
			Name:  "valueOf",
			Type:  "(D)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_STRING:
		return
	case ast.VARIABLE_TYPE_OBJECT:
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst(fmt.Sprintf("object@%s", typ.Class.Name), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst("[]", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	default:
		panic(1111111111)
	}
}
