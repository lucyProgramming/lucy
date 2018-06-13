package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) mkBuildinFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	switch call.Func.Name {
	case common.BUILD_IN_FUNCTION_PRINT:
		return m.mkBuildInPrint(class, code, e, context, state)
	case common.BUILD_IN_FUNCTION_PANIC:
		return m.mkBuildinPanic(class, code, e, context, state)
	case common.BUILD_IN_FUNCTION_CATCH:
		return m.mkBuildInCatch(class, code, e, context)
	case common.BUILD_IN_FUNCTION_MONITORENTER, common.BUILD_IN_FUNCTION_MONITOREXIT:
		maxStack, _ = m.build(class, code, call.Args[0], context, state)
		if call.Func.Name == common.BUILD_IN_FUNCTION_MONITORENTER {
			code.Codes[code.CodeLength] = cg.OP_monitorenter
		} else { // monitor enter on exit
			code.Codes[code.CodeLength] = cg.OP_monitorexit
		}
		code.CodeLength++
	case common.BUILD_IN_FUNCTION_PRINTF:
		return m.mkBuildInPrintf(class, code, e, context, state)
	case common.BUILD_IN_FUNCTION_SPRINTF:
		return m.mkBuildInSprintf(class, code, e, context, state)
	case common.BUILD_IN_FUNCTION_LEN:
		return m.mkBuildInLen(class, code, e, context, state)
	default:
		panic("unkown buildin function:" + call.Func.Name)
	}
	return
}

func (m *MakeExpression) mkBuildinPanic(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Args[0].Typ != ast.EXPRESSION_TYPE_NEW { // not new expression
		code.Codes[code.CodeLength] = cg.OP_new
		className := call.Args[0].Value.Class.Name
		class.InsertClassConst(className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
		{
			t := &cg.StackMap_verification_type_info{}
			tt := &cg.StackMap_Uninitialized_variable_info{}
			tt.CodeOffset = uint16(code.CodeLength - 4)
			t.Verify = tt
			state.Stacks = append(state.Stacks, t)
			state.Stacks = append(state.Stacks, t)
		}
		stack, _ := m.build(class, code, call.Args[0], context, state)
		state.popStack(2)
		maxStack = 2 + stack
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      className,
			Method:     special_method_init,
			Descriptor: "(Ljava/lang/Throwable;)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		maxStack, _ = m.build(class, code, call.Args[0], context, state)
	}
	code.Codes[code.CodeLength] = cg.OP_athrow
	code.CodeLength++
	context.MakeStackMap(code, state, code.CodeLength)
	return
}

func (m *MakeExpression) mkBuildInCatch(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context) (maxStack uint16) {
	if e.IsStatementExpression { // statement call
		maxStack = 1
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...)
		return
	}
	maxStack = 2
	//load to stack
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...) // load
	//set 2 null
	code.Codes[code.CodeLength] = cg.OP_aconst_null
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...) // store
	//check cast
	code.Codes[code.CodeLength] = cg.OP_checkcast
	if context.Defer.ExceptionClass != nil {
		class.InsertClassConst(context.Defer.ExceptionClass.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		class.InsertClassConst(ast.DEFAULT_EXCEPTION_CLASS, code.Codes[code.CodeLength+1:code.CodeLength+3])
	}
	code.CodeLength += 3
	return
}

func (m *MakeExpression) mkBuildInLen(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	maxStack, _ = m.build(class, code, call.Args[0], context, state)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 2 > maxStack {
		maxStack = 2
	}
	exit := (&cg.Exit{}).FromCode(cg.OP_ifnull, code)
	//binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 3)
	if call.Args[0].Value.Typ == ast.VARIABLE_TYPE_JAVA_ARRAY {
		code.Codes[code.CodeLength] = cg.OP_arraylength
		code.CodeLength++
	} else if call.Args[0].Value.Typ == ast.VARIABLE_TYPE_ARRAY {
		meta := ArrayMetas[call.Args[0].Value.ArrayType.Typ]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     "size",
			Descriptor: "()I",
		},
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else if call.Args[0].Value.Typ == ast.VARIABLE_TYPE_MAP {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Method:     "size",
			Descriptor: "()I",
		},
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else if call.Args[0].Value.Typ == ast.VARIABLE_TYPE_STRING {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "length",
			Descriptor: "()I",
		},
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	backfillExit([]*cg.Exit{exit}, code.CodeLength+3)
	state.pushStack(class, call.Args[0].Value)
	context.MakeStackMap(code, state, code.CodeLength+3)
	state.popStack(1)
	code.Codes[code.CodeLength] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 5)
	code.Codes[code.CodeLength+3] = cg.OP_pop
	code.Codes[code.CodeLength+4] = cg.OP_iconst_0
	code.CodeLength += 5
	state.pushStack(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	if e.IsStatementExpression {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}
	return
}
