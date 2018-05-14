package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type MakeExpression struct {
	MakeClass *MakeClass
}

func (m *MakeExpression) build(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16, exits []*cg.JumpBackPatch) {
	if e.IsCompileAuto == false {
		context.appendLimeNumberAndSourceFile(e.Pos, code, class)
	}
	switch e.Typ {
	case ast.EXPRESSION_TYPE_TYPE_ALIAS:
		return // handled at ast stage
	case ast.EXPRESSION_TYPE_NULL:
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		maxstack = 1
	case ast.EXPRESSION_TYPE_BOOL:
		if e.Data.(bool) {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_0
		}
		code.CodeLength++
		maxstack = 1
	case ast.EXPRESSION_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_bipush
		code.Codes[code.CodeLength+1] = e.Data.(byte)
		code.CodeLength += 2
		maxstack = 1
	case ast.EXPRESSION_TYPE_INT, ast.EXPRESSION_TYPE_SHORT:
		loadInt32(class, code, e.Data.(int32))
		maxstack = 1
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
		maxstack = 2
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
		maxstack = 1
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
		maxstack = 2
	case ast.EXPRESSION_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst(e.Data.(string), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if len([]byte(e.Data.(string))) > 65536 {
			panic("jvm max string length is 65536")
		}
		maxstack = 1
	//binary expression
	case ast.EXPRESSION_TYPE_LOGICAL_OR:
		fallthrough
	case ast.EXPRESSION_TYPE_LOGICAL_AND:
		maxstack, exits = m.buildLogical(class, code, e, context, state)
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
		maxstack = m.buildArithmetic(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_ASSIGN:
		maxstack = m.buildAssign(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_COLON_ASSIGN:
		maxstack = m.buildColonAssign(class, code, e, context, state)
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
		maxstack = m.buildOpAssign(class, code, e, context, state)
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
		maxstack = m.buildRelations(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_INDEX:
		maxstack = m.buildIndex(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_DOT:
		maxstack = m.buildDot(class, code, e, context, state)

	//
	case ast.EXPRESSION_TYPE_METHOD_CALL:
		maxstack = m.buildMethodCall(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_FUNCTION_CALL:
		maxstack = m.buildFunctionCall(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_INCREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_DECREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_PRE_INCREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_PRE_DECREMENT:
		maxstack = m.buildSelfIncrement(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_BITWISE_NOT:
		fallthrough
	case ast.EXPRESSION_TYPE_NEGATIVE:
		fallthrough
	case ast.EXPRESSION_TYPE_NOT:
		maxstack = m.buildUnary(class, code, e, context, state)
	//
	case ast.EXPRESSION_TYPE_IDENTIFIER:
		maxstack = m.buildIdentifer(class, code, e, context)
	case ast.EXPRESSION_TYPE_NEW:
		maxstack = m.buildNew(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_FUNCTION:
	case ast.EXPRESSION_TYPE_CHECK_CAST: // []byte(str)
		maxstack = m.buildTypeConvertion(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_CONST: // const will analyse at ast stage
	case ast.EXPRESSION_TYPE_SLICE:
		maxstack = m.buildSlice(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_ARRAY:
		maxstack = m.buildArray(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_MAP:
		maxstack = m.buildMapLiteral(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_VAR:
		maxstack = m.buildColonAssign(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_TYPE_ASSERT:
		maxstack = m.buildTypeAssert(class, code, e, context, state)
	case ast.EXPRESSION_TYPE_TERNARY:
		maxstack = m.buildTernary(class, code, e, context, state)
	default:
		panic(e.OpName())
	}

	return
}

func (m *MakeExpression) valueJvmSize(e *ast.Expression) (size uint16) {
	if len(e.Values) > 1 {
		return 1
	}
	if e.Value.RightValueValid() == false {
		return 0
	}
	return jvmSize(e.Value)
}
