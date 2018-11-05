package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildExpression) buildNew(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	switch e.Value.Type {
	case ast.VariableTypeArray:
		return this.buildNewArray(class, code, e, context, state)
	case ast.VariableTypeJavaArray:
		return this.buildNewJavaArray(class, code, e, context, state)
	case ast.VariableTypeMap:
		return this.buildNewMap(class, code, e, context)
	}
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	//new object
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
	maxStack += this.buildCallArgs(class, code, n.Args, n.VArgs, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	d := n.Construction.Function.JvmDescriptor
	if d == "" {
		d = Descriptor.methodDescriptor(&n.Construction.Function.Type)
	}
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      n.Type.Class.Name,
		Method:     specialMethodInit,
		Descriptor: d,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}

func (this *BuildExpression) buildNewMap(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context) (maxStack uint16) {
	maxStack = 2
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(mapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      mapClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}

func (this *BuildExpression) buildNewJavaArray(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
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
	maxStack = this.build(class, code, n.Args[0], context, state) // must be a integer
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
	if e.Value.Array.Type == ast.VariableTypeEnum {
		state.pushStack(class, e.Value)
		defer state.popStack(1)
		if t := 3 + setEnumArray(class, code, state, context, e.Value.Array.Enum); t > maxStack {
			maxStack = t
		}
	}
	return
}
func (this *BuildExpression) buildNewArray(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
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
	// get amount
	stack := this.build(class, code, n.Args[0], context, state) // must be a integer
	if t := 2 + stack; t > maxStack {
		maxStack = t
	}
	newArrayBaseOnType(class, code, e.Value.Array)
	if e.Value.Array.Type == ast.VariableTypeEnum {
		state.pushStack(class, &ast.Type{
			Type:  ast.VariableTypeJavaArray,
			Array: e.Value.Array,
		})

		if t := 3 + setEnumArray(class, code, state, context, e.Value.Array.Enum); t > maxStack {
			maxStack = t
		}
		state.popStack(1)
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      meta.className,
		Method:     specialMethodInit,
		Descriptor: meta.constructorFuncDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
