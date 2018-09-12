package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildIndex(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	if index.Expression.Value.Type == ast.VariableTypeMap {
		return buildExpression.buildMapIndex(class, code, e, context, state)
	}
	if index.Expression.Value.Type == ast.VariableTypeString {
		return buildExpression.buildStringIndex(class, code, e, context, state)
	}
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	maxStack = buildExpression.build(class, code, index.Expression, context, state)
	state.pushStack(class, index.Expression.Value)
	stack := buildExpression.build(class, code, index.Index, context, state)
	if t := stack + 1; t > maxStack {
		maxStack = t
	}
	if index.Expression.Value.Type == ast.VariableTypeArray {
		meta := ArrayMetas[e.Value.Type]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     "get",
			Descriptor: meta.getMethodDescription,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.Value.IsPointer() && e.Value.Type != ast.VariableTypeString {
			typeConverter.castPointer(class, code, e.Value)
		}
	} else {
		switch e.Value.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			code.Codes[code.CodeLength] = cg.OP_baload
		case ast.VariableTypeShort:
			code.Codes[code.CodeLength] = cg.OP_saload
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_iaload
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_laload
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_faload
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_daload
		case ast.VariableTypeString:
			fallthrough
		case ast.VariableTypeObject:
			fallthrough
		case ast.VariableTypeMap:
			fallthrough
		case ast.VariableTypeArray:
			fallthrough
		case ast.VariableTypeFunction:
			fallthrough
		case ast.VariableTypeJavaArray:
			code.Codes[code.CodeLength] = cg.OP_aaload
		}
		code.CodeLength++
	}
	return
}

func (buildExpression *BuildExpression) buildStringIndex(class *cg.ClassHighLevel,
	code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	index := e.Data.(*ast.ExpressionIndex)
	maxStack = buildExpression.build(class, code, index.Expression, context, state)
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringClass,
		Method:     "getBytes",
		Descriptor: "()[B",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.pushStack(class, state.newObjectVariableType("[B"))
	stack := buildExpression.build(class, code, index.Index, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_baload
	code.CodeLength++
	return
}
func (buildExpression *BuildExpression) buildMapIndex(class *cg.ClassHighLevel,
	code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	index := e.Data.(*ast.ExpressionIndex)
	maxStack = buildExpression.build(class, code, index.Expression, context, state)
	currentStack := uint16(1)
	//build index
	state.pushStack(class, index.Expression.Value)
	stack := buildExpression.build(class, code, index.Index, context, state)
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	currentStack = 2 // mapref kref
	if index.Expression.Value.Map.K.IsPointer() == false {
		typeConverter.packPrimitives(class, code, index.Expression.Value.Map.K)
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      mapClass,
		Method:     "get",
		Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.popStack(1)
	if index.Expression.Value.Map.V.Type == ast.VariableTypeEnum {
		typeConverter.unPackPrimitives(class, code, index.Expression.Value.Map.V)
	} else if index.Expression.Value.Map.V.IsPointer() {
		typeConverter.castPointer(class, code, index.Expression.Value.Map.V)
	} else {
		code.Codes[code.CodeLength] = cg.OP_dup // increment the stack
		code.CodeLength++
		if 2 > maxStack { // stack is  ... valueObjectRef valueObjectRef
			maxStack = 2
		}
		noNullExit := (&cg.Exit{}).Init(cg.OP_ifnonnull, code)
		switch index.Expression.Value.Map.V.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_iconst_0
			code.CodeLength += 2
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_lconst_0
			code.CodeLength += 2
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_fconst_0
			code.CodeLength += 2
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_dconst_0
			code.CodeLength += 2
		}
		nullExit := (&cg.Exit{}).Init(cg.OP_goto, code)
		state.pushStack(class, state.newObjectVariableType(javaRootClass))
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1) // pop java_root_class ref
		writeExits([]*cg.Exit{noNullExit}, code.CodeLength)
		typeConverter.unPackPrimitives(class, code, index.Expression.Value.Map.V)
		writeExits([]*cg.Exit{nullExit}, code.CodeLength)
		state.pushStack(class, e.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
	}
	return
}
