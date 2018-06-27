package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildNew(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	if e.ExpressionValue.Type == ast.VariableTypeArray {
		return makeExpression.buildNewArray(class, code, e, context, state)
	}
	if e.ExpressionValue.Type == ast.VariableTypeJavaArray {
		return makeExpression.buildNewJavaArray(class, code, e, context, state)
	}
	if e.ExpressionValue.Type == ast.VariableTypeMap {
		return makeExpression.buildNewMap(class, code, e, context)
	}
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()

	//new class
	n := e.Data.(*ast.ExpressionNew)
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(n.Type.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	t := &cg.StackMapVerificationTypeInfo{}
	t.Verify = &cg.StackMapUninitializedVariableInfo{
		CodeOffset: uint16(code.CodeLength),
	}
	state.Stacks = append(state.Stacks, t, t)
	code.CodeLength += 4
	maxStack = 2
	if n.Args != nil && len(n.Args) > 0 {
		maxStack += makeExpression.buildCallArgs(class, code, n.Args, n.Construction.Function.Type.ParameterList, context, state)
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	if n.Construction == nil {
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      n.Type.Class.Name,
			Method:     specialMethodInit,
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		d := ""
		if n.Type.Class.LoadFromOutSide {
			d = n.Construction.Function.Descriptor
		} else {
			d = Descriptor.methodDescriptor(&n.Construction.Function.Type)
		}
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      n.Type.Class.Name,
			Method:     specialMethodInit,
			Descriptor: d,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	}
	code.CodeLength += 3
	return
}
func (makeExpression *MakeExpression) buildNewMap(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context) (maxStack uint16) {
	maxStack = 2
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(javaHashMapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaHashMapClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}

func (makeExpression *MakeExpression) buildNewJavaArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	dimensions := byte(0)
	{
		// get dimension
		t := e.ExpressionValue
		for t.Type == ast.VariableTypeJavaArray {
			dimensions++
			t = t.ArrayType
		}
	}
	n := e.Data.(*ast.ExpressionNew)
	maxStack, _ = makeExpression.build(class, code, n.Args[0], context, state) // must be a integer
	currentStack := uint16(1)
	for i := byte(0); i < dimensions-1; i++ {
		loadInt32(class, code, 0)
		currentStack++
		if currentStack > maxStack {
			maxStack = currentStack
		}
	}
	code.Codes[code.CodeLength] = cg.OP_multianewarray
	class.InsertClassConst(Descriptor.typeDescriptor(e.ExpressionValue), code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = dimensions
	code.CodeLength += 4
	return
}
func (makeExpression *MakeExpression) buildNewArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	//new
	n := e.Data.(*ast.ExpressionNew)
	meta := ArrayMetas[e.ExpressionValue.ArrayType.Type]
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxStack = 2
	{
		t := &cg.StackMapVerificationTypeInfo{}
		unInit := &cg.StackMapUninitializedVariableInfo{}
		unInit.CodeOffset = uint16(code.CodeLength - 4)
		t.Verify = unInit
		state.Stacks = append(state.Stacks, t, t) // 2 for dup
		defer state.popStack(2)
	}
	if n.IsConvertJavaArray2Array {
		stack, _ := makeExpression.build(class, code, n.Args[0], context, state) // must be a integer
		if t := 2 + stack; t > maxStack {
			maxStack = t
		}
	} else {
		// get amount
		stack, _ := makeExpression.build(class, code, n.Args[0], context, state) // must be a integer
		if t := 2 + stack; t > maxStack {
			maxStack = t
		}
		switch e.ExpressionValue.ArrayType.Type {
		case ast.VariableTypeBool:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ATYPE_T_BOOLEAN
			code.CodeLength += 2
		case ast.VariableTypeByte:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ATYPE_T_BYTE
			code.CodeLength += 2
		case ast.VariableTypeShort:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ATYPE_T_SHORT
			code.CodeLength += 2
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ATYPE_T_INT
			code.CodeLength += 2
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ATYPE_T_LONG
			code.CodeLength += 2
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ATYPE_T_FLOAT
			code.CodeLength += 2
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ATYPE_T_DOUBLE
			code.CodeLength += 2
		case ast.VariableTypeString:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(javaStringClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeMap:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(javaHashMapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeFunction:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(javaMethodHandleClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeObject:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(e.ExpressionValue.ArrayType.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeArray:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			meta := ArrayMetas[e.ExpressionValue.ArrayType.ArrayType.Type]
			class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeJavaArray:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(Descriptor.typeDescriptor(e.ExpressionValue.ArrayType), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.className,
		Method:     specialMethodInit,
		Descriptor: meta.constructorFuncDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
