package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildMapIndex(class *cg.ClassHighLevel,
	code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	index := e.Data.(*ast.ExpressionIndex)
	maxStack, _ = buildExpression.build(class, code, index.Expression, context, state)
	currentStack := uint16(1)
	//build index
	state.pushStack(class, index.Expression.ExpressionValue)
	stack, _ := buildExpression.build(class, code, index.Index, context, state)
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	currentStack = 2 // mapref kref
	if index.Expression.ExpressionValue.Map.Key.IsPointer() == false {
		typeConverter.packPrimitives(class, code, index.Expression.ExpressionValue.Map.Key)
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaMapClass,
		Method:     "get",
		Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.popStack(1)
	if index.Expression.ExpressionValue.Map.Value.Type == ast.VariableTypeEnum {
		typeConverter.unPackPrimitives(class, code, index.Expression.ExpressionValue.Map.Value)
	} else if index.Expression.ExpressionValue.Map.Value.IsPointer() {
		typeConverter.castPointer(class, code, index.Expression.ExpressionValue.Map.Value)
	} else {
		code.Codes[code.CodeLength] = cg.OP_dup // incrment the stack
		code.CodeLength++
		if t := 1 + currentStack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_ifnonnull
		codeLength := code.CodeLength
		code.CodeLength += 3
		switch index.Expression.ExpressionValue.Map.Value.Type {
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
		code.Codes[code.CodeLength] = cg.OP_goto
		codeLength2 := code.CodeLength
		code.CodeLength += 3
		// no null goto here
		{
			state.pushStack(class, state.newObjectVariableType(javaRootClass))
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1) // pop java_root_class ref
		}
		binary.BigEndian.PutUint16(code.Codes[codeLength+1:codeLength+3], uint16(code.CodeLength-codeLength))
		typeConverter.unPackPrimitives(class, code, index.Expression.ExpressionValue.Map.Value)
		binary.BigEndian.PutUint16(code.Codes[codeLength2+1:codeLength2+3], uint16(code.CodeLength-codeLength2))
		{
			state.pushStack(class, e.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
	}
	return
}

func (buildExpression *BuildExpression) buildIndex(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	index := e.Data.(*ast.ExpressionIndex)
	if index.Expression.ExpressionValue.Type == ast.VariableTypeMap {
		return buildExpression.buildMapIndex(class, code, e, context, state)
	}
	maxStack, _ = buildExpression.build(class, code, index.Expression, context, state)
	state.pushStack(class, index.Expression.ExpressionValue)
	stack, _ := buildExpression.build(class, code, index.Index, context, state)
	if t := stack + 1; t > maxStack {
		maxStack = t
	}
	if index.Expression.ExpressionValue.Type == ast.VariableTypeArray {
		meta := ArrayMetas[e.ExpressionValue.Type]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     "get",
			Descriptor: meta.getMethodDescription,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.ExpressionValue.IsPointer() && e.ExpressionValue.Type != ast.VariableTypeString {
			typeConverter.castPointer(class, code, e.ExpressionValue)
		}
	} else {
		switch e.ExpressionValue.Type {
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
