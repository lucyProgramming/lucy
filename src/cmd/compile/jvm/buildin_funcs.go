package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	switch call.Func.Name {
	case common.BUILD_IN_FUNCTION_PRINT:
		return m.mkBuildinPrint(class, code, call, context)
	case common.BUILD_IN_FUNCTION_PANIC:
		return m.mkBuildinPanic(class, code, call, context)
	case common.BUILD_IN_FUNCTION_CATCH:
		return m.mkBuildinRecover(class, code, e, context)
	case common.BUILD_IN_FUNCTION_MONITORENTER, common.BUILD_IN_FUNCTION_MONITOREXIT:
		maxstack, _ = m.build(class, code, call.Args[0], context)
		if call.Func.Name == common.BUILD_IN_FUNCTION_MONITORENTER {
			code.Codes[code.CodeLength] = cg.OP_monitorenter
		} else {
			code.Codes[code.CodeLength] = cg.OP_monitorexit
		}
		code.CodeLength++
	}
}
func (m *MakeExpression) mkBuildinPanic(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context) (maxstack uint16) {
	maxstack, _ = m.build(class, code, call.Args[0], context)
	code.Codes[code.CodeLength] = cg.OP_athrow
	code.CodeLength++
	return
}

/*
	function print
*/
func (m *MakeExpression) mkBuildinPrint(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_getstatic
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      "java/lang/System",
		Field:      "out",
		Descriptor: "Ljava/io/PrintStream;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 1
	if len(call.Args) == 0 {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/io/PrintStream",
			Method:     "println",
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}

	if len(call.Args) == 1 && call.Args[0].HaveOnlyOneValue() {
		stack, es := m.build(class, code, call.Args[0], context)
		backPatchEs(es, code.CodeLength)
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
		switch call.Args[0].GetTheOnlyOneVariableType().Typ {
		case ast.VARIABLE_TYPE_BOOL:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Z)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(I)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(J)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(F)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(D)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_STRING:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Ljava/lang/String;)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY:
			fallthrough
		case ast.VARIABLE_TYPE_MAP:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Object",
				Method:     "toString",
				Descriptor: "()Ljava/lang/String;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Ljava/lang/String;)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}

	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 3
	currentStack := uint16(2)
	app := func(appendBlankSpace bool) {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/StringBuilder",
			Method:     "append",
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
				Method:     "append",
				Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
	for k, v := range call.Args {
		var variableType *ast.VariableType
		if v.IsCall() && len(v.VariableTypes) > 1 {
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
		if v.IsCall() {
			variableType = v.VariableTypes[0]
		}
		stack, es := m.build(class, code, v, context)
		backPatchEs(es, code.CodeLength)
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
		Method:     "toString",
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// call println
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/io/PrintStream",
		Method:     "println",
		Descriptor: "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
func (m *MakeExpression) mkBuildinRecover(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	if e.IsStatementExpression { // first level statement
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...)
		return
	}
	maxstack = 2
	//load to stack
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...) // load
	//set 2 null
	code.Codes[code.CodeLength] = cg.OP_aconst_null
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...) // load
	return
}
