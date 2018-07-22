package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildNew(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	if e.Value.Type == ast.VariableTypeArray {
		return buildExpression.buildNewArray(class, code, e, context, state)
	}
	if e.Value.Type == ast.VariableTypeJavaArray {
		return buildExpression.buildNewJavaArray(class, code, e, context, state)
	}
	if e.Value.Type == ast.VariableTypeMap {
		return buildExpression.buildNewMap(class, code, e, context)
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
	maxStack += buildExpression.buildCallArgs(class, code, n.Args, n.VArgs, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokespecial
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
	code.CodeLength += 3
	return
}
func (buildExpression *BuildExpression) buildNewMap(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context) (maxStack uint16) {
	maxStack = 2
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(javaMapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaMapClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}

func (buildExpression *BuildExpression) buildNewJavaArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	dimensions := byte(0)
	{
		// get dimension
		t := e.Value
		for t.Type == ast.VariableTypeJavaArray {
			dimensions++
			t = t.Array
		}
	}
	n := e.Data.(*ast.ExpressionNew)
	maxStack, _ = buildExpression.build(class, code, n.Args[0], context, state) // must be a integer
	currentStack := uint16(1)
	for i := byte(0); i < dimensions-1; i++ {
		loadInt32(class, code, 0)
		currentStack++
		if currentStack > maxStack {
			maxStack = currentStack
		}
	}
	code.Codes[code.CodeLength] = cg.OP_multianewarray
	class.InsertClassConst(Descriptor.typeDescriptor(e.Value), code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = dimensions
	code.CodeLength += 4
	return
}
func (buildExpression *BuildExpression) buildNewArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	//new
	n := e.Data.(*ast.ExpressionNew)
	meta := ArrayMetas[e.Value.Array.Type]
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
		stack, _ := buildExpression.build(class, code, n.Args[0], context, state) // must be a integer
		if t := 2 + stack; t > maxStack {
			maxStack = t
		}
	} else {
		// get amount
		stack, _ := buildExpression.build(class, code, n.Args[0], context, state) // must be a integer
		if t := 2 + stack; t > maxStack {
			maxStack = t
		}
		switch e.Value.Array.Type {
		case ast.VariableTypeBool:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ArrayTypeBoolean
			code.CodeLength += 2
		case ast.VariableTypeByte:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ArrayTypeByte
			code.CodeLength += 2
		case ast.VariableTypeShort:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ArrayTypeShort
			code.CodeLength += 2
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ArrayTypeInt
			code.CodeLength += 2
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ArrayTypeLong
			code.CodeLength += 2
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ArrayTypeFloat
			code.CodeLength += 2
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_newarray
			code.Codes[code.CodeLength+1] = ArrayTypeDouble
			code.CodeLength += 2
		case ast.VariableTypeString:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(javaStringClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeMap:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(javaMapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeFunction:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(javaMethodHandleClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeObject:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(e.Value.Array.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeArray:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			meta := ArrayMetas[e.Value.Array.Array.Type]
			class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeJavaArray:
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(Descriptor.typeDescriptor(e.Value.Array), code.Codes[code.CodeLength+1:code.CodeLength+3])
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
