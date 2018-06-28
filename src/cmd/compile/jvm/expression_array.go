package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildArray(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	arr := e.Data.(*ast.ExpressionArray)
	//	new array ,
	meta := ArrayMetas[e.ExpressionValue.Array.Type]
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4

	{
		t := &cg.StackMapVerificationTypeInfo{}
		tt := &cg.StackMapUninitializedVariableInfo{}
		tt.CodeOffset = uint16(code.CodeLength - 4)
		t.Verify = tt
		state.Stacks = append(state.Stacks, t, t)
	}

	loadInt32(class, code, int32(arr.Length))
	switch e.ExpressionValue.Array.Type {
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
	case ast.VariableTypeMap:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(javaMapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeString:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(javaStringClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeFunction:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(javaMethodHandleClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeObject:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(e.ExpressionValue.Array.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeArray:
		meta := ArrayMetas[e.ExpressionValue.Array.Array.Type]
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	arrayObject := &ast.Type{}
	arrayObject.Type = ast.VariableTypeJavaArray
	arrayObject.Array = e.ExpressionValue.Array
	state.pushStack(class, arrayObject)

	maxStack = 4

	store := func() {
		switch e.ExpressionValue.Array.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			code.Codes[code.CodeLength] = cg.OP_bastore
		case ast.VariableTypeShort:
			code.Codes[code.CodeLength] = cg.OP_sastore
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_iastore
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_lastore
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_fastore
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_dastore
		case ast.VariableTypeFunction:
			fallthrough
		case ast.VariableTypeMap:
			fallthrough
		case ast.VariableTypeString:
			fallthrough
		case ast.VariableTypeObject:
			fallthrough
		case ast.VariableTypeArray:
			code.Codes[code.CodeLength] = cg.OP_aastore
		}
		code.CodeLength++
	}
	var index int32 = 0
	for _, v := range arr.Expressions {
		if v.MayHaveMultiValue() && len(v.ExpressionMultiValues) > 1 {
			// stack top is array list
			stack, _ := makeExpression.build(class, code, v, context, state)
			if t := 3 + stack; t > maxStack {
				maxStack = t
			}
			multiValuePacker.storeMultiValueAutoVar(code, context)
			for k, t := range v.ExpressionMultiValues {
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				loadInt32(class, code, index) // load index
				stack := multiValuePacker.unPack(class, code, k, t, context)
				if t := 5 + stack; t > maxStack {
					maxStack = t
				}
				store()
				index++
			}
			continue
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, index) // load index
		state.pushStack(class, arrayObject)
		state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
		stack, es := makeExpression.build(class, code, v, context, state)
		if len(es) > 0 {
			fillOffsetForExits(es, code.CodeLength)
			state.pushStack(class, v.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1) // must be a logical expression
		}
		state.popStack(2)
		if t := 5 + stack; t > maxStack {
			maxStack = t
		}
		store()
		index++
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
