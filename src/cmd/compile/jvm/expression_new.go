package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildNew(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	if e.Value.Typ == ast.VARIABLE_TYPE_ARRAY {
		return m.buildNewArray(class, code, e, context, state)
	}
	if e.Value.Typ == ast.VARIABLE_TYPE_MAP {
		return m.buildNewMap(class, code, e, context)
	}
	//new class
	n := e.Data.(*ast.ExpressionNew)
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(n.Typ.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2
	if context.method.CaptureFunctionLength > 0 {

	}
	if n.Args != nil && len(n.Args) > 0 {
		maxstack += m.buildCallArgs(class, code, n.Args, n.Construction.Func.Typ.ParameterList[context.method.CaptureFunctionLength:], context, state)
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	if n.Construction == nil {
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      n.Typ.Class.Name,
			Method:     special_method_init,
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		d := ""
		if n.Typ.Class.LoadFromOutSide {
			d = n.Construction.Func.Descriptor
		} else {
			d = Descriptor.methodDescriptor(n.Construction.Func)
		}
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      n.Typ.Class.Name,
			Method:     special_method_init,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	}
	code.CodeLength += 3
	return
}
func (m *MakeExpression) buildNewMap(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 2
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
func (m *MakeExpression) buildNewArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxstack uint16) {
	//new
	n := e.Data.(*ast.ExpressionNew)
	meta := ArrayMetas[e.Value.ArrayType.Typ]
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2
	{
		t := &cg.StackMap_verification_type_info{}
		unInit := &cg.StackMap_Uninitialized_variable_info{}
		unInit.Index = uint16(code.CodeLength - 4)
		t.Verify = unInit
		state.Stacks = append(state.Stacks, t, t) // 2 for dup
		defer state.popStack(2)
	}
	// call init
	stack, _ := m.build(class, code, n.Args[0], context, state) // must be a integer
	if t := 2 + stack; t > maxstack {
		maxstack = t
	}

	state.pushStack(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
	defer state.popStack(1)
	maxstack += stack
	currentStack := uint16(3)
	if currentStack > maxstack {
		maxstack = currentStack
	}
	if t := checkStackTopIfNagetiveThrowIndexOutOfRangeException(class, code, context, state) + currentStack; t > maxstack {
		maxstack = t
	}
	switch e.Value.ArrayType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_BOOLEAN
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_BYTE
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_SHORT
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_INT
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_LONG
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_FLOAT
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_DOUBLE
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(java_string_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_MAP:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(e.Value.ArrayType.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_ARRAY:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		meta := ArrayMetas[e.Value.ArrayType.ArrayType.Typ]
		class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Method:     special_method_init,
		Descriptor: meta.constructorFuncDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
