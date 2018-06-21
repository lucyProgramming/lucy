package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type MakeExpression struct {
	MakeClass *MakeClass
}

func (makeExpression *MakeExpression) build(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16, exits []*cg.Exit) {
	if e.IsCompileAuto == false {
		context.appendLimeNumberAndSourceFile(e.Pos, code, class)
	}
	switch e.Type {
	case ast.EXPRESSION_TYPE_TYPE_ALIAS:
		return // handled at ast stage
	case ast.EXPRESSION_TYPE_NULL:
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		maxStack = 1
	case ast.EXPRESSION_TYPE_BOOL:
		if e.Data.(bool) {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_0
		}
		code.CodeLength++
		maxStack = 1
	case ast.EXPRESSION_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_bipush
		code.Codes[code.CodeLength+1] = e.Data.(byte)
		code.CodeLength += 2
		maxStack = 1
	case ast.EXPRESSION_TYPE_INT, ast.EXPRESSION_TYPE_SHORT:
		loadInt32(class, code, e.Data.(int32))
		maxStack = 1
	case ast.EXPRESSION_TYPE_LONG:
		if e.Data.(int64) == 0 {
			code.Codes[code.CodeLength] = cg.OP_lconst_0
			code.CodeLength++
		} else if e.Data.(int64) == 1 {
			code.Codes[code.CodeLength] = cg.OP_lconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc2_w
			class.InsertLongConst(e.Data.(int64), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		maxStack = 2
	case ast.EXPRESSION_TYPE_FLOAT:
		if e.Data.(float32) == 0.0 {
			code.Codes[code.CodeLength] = cg.OP_fconst_0
			code.CodeLength++
		} else if e.Data.(float32) == 1.0 {
			code.Codes[code.CodeLength] = cg.OP_fconst_1
			code.CodeLength++
		} else if e.Data.(float32) == 2.0 {
			code.Codes[code.CodeLength] = cg.OP_fconst_2
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertFloatConst(e.Data.(float32), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		maxStack = 1
	case ast.EXPRESSION_TYPE_DOUBLE:
		if e.Data.(float64) == 0.0 {
			code.Codes[code.CodeLength] = cg.OP_dconst_0
			code.CodeLength++
		} else if e.Data.(float64) == 1.0 {
			code.Codes[code.CodeLength] = cg.OP_dconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc2_w
			class.InsertDoubleConst(e.Data.(float64), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		maxStack = 2
	case ast.EXPRESSION_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst(e.Data.(string), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if len([]byte(e.Data.(string))) > 65536 {
			panic("jvm max string length is 65536")
		}
		maxStack = 1
	//binary expression
	case ast.EXPRESSION_TYPE_LOGICAL_OR:
		fallthrough
	case ast.EXPRESSION_TYPE_LOGICAL_AND:
		maxStack, exits = makeExpression.buildLogical(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_OR:
		fallthrough
	case ast.EXPRESSION_TYPE_AND:
		fallthrough
	case ast.EXPRESSION_TYPE_XOR:
		fallthrough
	case ast.EXPRESSION_TYPE_LSH:
		fallthrough
	case ast.EXPRESSION_TYPE_RSH:
		fallthrough
	case ast.EXPRESSION_TYPE_ADD:
		fallthrough
	case ast.EXPRESSION_TYPE_SUB:
		fallthrough
	case ast.EXPRESSION_TYPE_MUL:
		fallthrough
	case ast.EXPRESSION_TYPE_DIV:
		fallthrough
	case ast.EXPRESSION_TYPE_MOD:
		maxStack = makeExpression.buildArithmetic(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_ASSIGN:
		maxStack = makeExpression.buildAssign(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_COLON_ASSIGN:
		maxStack = makeExpression.buildColonAssign(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_PLUS_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_MINUS_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_MUL_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_DIV_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_MOD_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_AND_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_OR_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_LSH_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_RSH_ASSIGN:
		fallthrough
	case ast.EXPRESSION_TYPE_XOR_ASSIGN:
		maxStack = makeExpression.buildOpAssign(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_EQ:
		fallthrough
	case ast.EXPRESSION_TYPE_NE:
		fallthrough
	case ast.EXPRESSION_TYPE_GE:
		fallthrough
	case ast.EXPRESSION_TYPE_GT:
		fallthrough
	case ast.EXPRESSION_TYPE_LE:
		fallthrough
	case ast.EXPRESSION_TYPE_LT:
		maxStack = makeExpression.buildRelations(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_INDEX:
		maxStack = makeExpression.buildIndex(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_SELECTION:
		maxStack = makeExpression.buildSelection(class, code, e, context, state)

	//
	case ast.EXPRESSION_TYPE_METHOD_CALL:
		maxStack = makeExpression.buildMethodCall(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_FUNCTION_CALL:
		maxStack = makeExpression.buildFunctionCall(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_INCREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_DECREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_PRE_INCREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_PRE_DECREMENT:
		maxStack = makeExpression.buildSelfIncrement(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_BIT_NOT:
		fallthrough
	case ast.EXPRESSION_TYPE_NEGATIVE:
		fallthrough
	case ast.EXPRESSION_TYPE_NOT:
		maxStack = makeExpression.buildUnary(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_IDENTIFIER:
		maxStack = makeExpression.buildIdentifier(class, code, e, context)
	case ast.EXPRESSION_TYPE_NEW:
		maxStack = makeExpression.buildNew(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_FUNCTION:
	case ast.EXPRESSION_TYPE_CHECK_CAST: // []byte(str)
		maxStack = makeExpression.buildTypeConversion(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_CONST: // const will analyse at ast stage
	case ast.EXPRESSION_TYPE_SLICE:
		maxStack = makeExpression.buildSlice(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_ARRAY:
		maxStack = makeExpression.buildArray(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_MAP:
		maxStack = makeExpression.buildMapLiteral(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_VAR:
		maxStack = makeExpression.buildVar(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_TYPE_ASSERT:
		maxStack = makeExpression.buildTypeAssert(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_TERNARY:
		maxStack = makeExpression.buildTernary(class, code, e, context, state)
	default:
		panic(e.OpName())
	}

	return
}

func (makeExpression *MakeExpression) valueJvmSize(e *ast.Expression) (size uint16) {
	if len(e.ExpressionMultiValues) > 1 {
		return 1
	}
	if e.ExpressionValue.RightValueValid() == false {
		return 0
	}
	return jvmSlotSize(e.ExpressionValue)
}

func (makeExpression *MakeExpression) buildExpressions(class *cg.ClassHighLevel, code *cg.AttributeCode,
	es []*ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := 0
	for _, e := range es {
		if e.MayHaveMultiValue() {
			length += len(e.ExpressionMultiValues)
			continue
		}
		length++
	}
	loadInt32(class, code, int32(length))
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst(java_root_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if 1 > maxStack {
		maxStack = 1
	}

	arrayListObject := state.newObjectVariableType(java_root_object_array)
	state.pushStack(class, arrayListObject)
	state.pushStack(class, arrayListObject)
	defer state.popStack(2)
	index := int32(0)
	for _, v := range es {
		currentStack := uint16(1)
		if v.MayHaveMultiValue() && len(v.ExpressionMultiValues) > 1 {
			stack, _ := makeExpression.build(class, code, v, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			for kk, _ := range v.ExpressionMultiValues {
				currentStack = 1
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				currentStack++
				stack = multiValuePacker.unPackObject(class, code, kk, context)
				if t := stack + currentStack; t > maxStack {
					maxStack = t
				}
				loadInt32(class, code, index)
				code.Codes[code.CodeLength] = cg.OP_swap
				code.Codes[code.CodeLength+1] = cg.OP_aastore
				code.CodeLength += 2
				index++
			}
			continue
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack++
		stack, es := makeExpression.build(class, code, v, context, state)
		if len(es) > 0 {
			fillOffsetForExits(es, code.CodeLength)
			state.pushStack(class, v.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.ExpressionValue.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.ExpressionValue)
		}
		loadInt32(class, code, index)
		code.Codes[code.CodeLength] = cg.OP_swap
		code.Codes[code.CodeLength+1] = cg.OP_aastore
		code.CodeLength += 2
		index++
	}
	return
}
