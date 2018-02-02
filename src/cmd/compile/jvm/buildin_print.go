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
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
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
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       specail_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 3
	currentStack := uint16(2)
	app := func(appendBlankSpace bool) {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/StringBuilder",
			Name:       "append",
			Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if appendBlankSpace {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertStringConst(" ", code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/StringBuilder",
				Name:       "append",
				Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
	for k, v := range call.Args {
		var variableType *ast.VariableType
		if (v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL) &&
			len(v.VariableTypes) > 1 {
			stack := m.buildFunctionCall(class, code, v, context)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			m.buildStoreArrayListAutoVar(code, context)
			for kk, tt := range v.VariableTypes {
				stack = m.unPackArraylist(class, code, kk, tt, context)
				if t := stack + currentStack; t > maxstack {
					maxstack = t
				}
				m.stackTop2String(class, code, tt)
				app(k < len(call.Args)-1 || kk < len(v.VariableTypes)-1)
			}
			continue
		}
		variableType = v.VariableType
		if v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL {
			variableType = v.VariableTypes[0]
		}
		stack, es := m.build(class, code, v, context)
		backPatchEs(es, code)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		m.stackTop2String(class, code, variableType)
		app(k < len(call.Args)-1)
	}
	// tostring
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Name:       "toString",
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// call println
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/io/PrintStream",
		Name:       "println",
		Descriptor: "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
