package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildMapIndex(class *cg.ClassHighLevel,
	code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	index := e.Data.(*ast.ExpressionIndex)
	maxStack, _ = makeExpression.build(class, code, index.Expression, context, state)
	currentStack := uint16(1)
	//build index
	state.pushStack(class, index.Expression.ExpressionValue)
	stack, _ := makeExpression.build(class, code, index.Index, context, state)
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	currentStack = 2 // mapref kref
	if index.Expression.ExpressionValue.Map.K.IsPointer() == false {
		typeConverter.packPrimitives(class, code, index.Expression.ExpressionValue.Map.K)
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Method:     "get",
		Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.popStack(1)
	if index.Expression.ExpressionValue.Map.V.Type == ast.VARIABLE_TYPE_ENUM {
		typeConverter.unPackPrimitives(class, code, index.Expression.ExpressionValue.Map.V)
	} else if index.Expression.ExpressionValue.Map.V.IsPointer() {
		typeConverter.castPointerTypeToRealType(class, code, index.Expression.ExpressionValue.Map.V)
	} else {
		code.Codes[code.CodeLength] = cg.OP_dup // incrment the stack
		code.CodeLength++
		if t := 1 + currentStack; t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_ifnonnull
		codeLength := code.CodeLength
		code.CodeLength += 3
		switch index.Expression.ExpressionValue.Map.V.Type {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_iconst_0
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_lconst_0
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_fconst_0
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength+1] = cg.OP_dconst_0
			code.CodeLength += 2
		}
		code.Codes[code.CodeLength] = cg.OP_goto
		codeLength2 := code.CodeLength
		code.CodeLength += 3
		// no null goto here
		{
			state.pushStack(class, state.newObjectVariableType(java_root_class))
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1) // pop java_root_class ref
		}
		binary.BigEndian.PutUint16(code.Codes[codeLength+1:codeLength+3], uint16(code.CodeLength-codeLength))
		typeConverter.unPackPrimitives(class, code, index.Expression.ExpressionValue.Map.V)
		binary.BigEndian.PutUint16(code.Codes[codeLength2+1:codeLength2+3], uint16(code.CodeLength-codeLength2))
		{
			state.pushStack(class, e.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
	}
	return
}

func (makeExpression *MakeExpression) buildIndex(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	index := e.Data.(*ast.ExpressionIndex)
	if index.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_MAP {
		return makeExpression.buildMapIndex(class, code, e, context, state)
	}
	maxStack, _ = makeExpression.build(class, code, index.Expression, context, state)
	state.pushStack(class, index.Expression.ExpressionValue)
	currentStack := uint16(1)
	if index.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_ARRAY {
		meta := ArrayMetas[e.ExpressionValue.Type]
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "end",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_dup_x1
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "start",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		state.pushStack(class, &ast.Type{Type: ast.VARIABLE_TYPE_INT})
		state.pushStack(class, &ast.Type{Type: ast.VARIABLE_TYPE_INT})
		currentStack = 3
	}
	stack, _ := makeExpression.build(class, code, index.Index, context, state)
	if t := stack + currentStack; t > maxStack {
		maxStack = t
	}
	if index.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_ARRAY {
		meta := ArrayMetas[e.ExpressionValue.Type]
		// stack arrayref  end start index
		code.Codes[code.CodeLength] = cg.OP_iadd
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_dup_x1
		code.CodeLength++
		{
			state.popStack(1)
			context.MakeStackMap(code, state, code.CodeLength+6)
			context.MakeStackMap(code, state, code.CodeLength+16)
		}
		code.Codes[code.CodeLength] = cg.OP_if_icmple
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
		code.Codes[code.CodeLength+3] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 13)
		code.Codes[code.CodeLength+6] = cg.OP_pop // incase stack over flow
		code.Codes[code.CodeLength+7] = cg.OP_pop
		code.Codes[code.CodeLength+8] = cg.OP_new
		class.InsertClassConst(java_index_out_of_range_exception_class, code.Codes[code.CodeLength+9:code.CodeLength+11])
		code.Codes[code.CodeLength+11] = cg.OP_dup
		code.Codes[code.CodeLength+12] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_index_out_of_range_exception_class,
			Method:     special_method_init,
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+13:code.CodeLength+15])
		code.Codes[code.CodeLength+15] = cg.OP_athrow
		// index not out of range
		code.Codes[code.CodeLength+16] = cg.OP_swap
		code.Codes[code.CodeLength+17] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "elements",
			Descriptor: meta.elementsFieldDescriptor,
		}, code.Codes[code.CodeLength+18:code.CodeLength+20])
		code.CodeLength += 20
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
	}
	switch e.ExpressionValue.Type {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_baload
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_saload
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_iaload
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_laload
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_faload
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_daload
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY:
		fallthrough
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		code.Codes[code.CodeLength] = cg.OP_aaload
	}
	code.CodeLength++
	if index.Expression.Type == ast.VARIABLE_TYPE_ARRAY &&
		e.ExpressionValue.IsPointer() && e.ExpressionValue.Type != ast.VARIABLE_TYPE_STRING {
		typeConverter.castPointerTypeToRealType(class, code, e.ExpressionValue)
	}
	return
}
