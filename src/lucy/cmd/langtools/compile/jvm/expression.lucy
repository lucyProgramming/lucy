package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type BuildExpression struct {
	BuildPackage *BuildPackage
}

func (this *BuildExpression) build(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	if e.IsCompileAuto == false {
		context.appendLimeNumberAndSourceFile(e.Pos, code, class)
	}
	switch e.Type {
	case ast.ExpressionTypeNull:
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		maxStack = 1
	case ast.ExpressionTypeBool:
		if e.Data.(bool) {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_0
		}
		code.CodeLength++
		maxStack = 1
	case ast.ExpressionTypeByte:
		code.Codes[code.CodeLength] = cg.OP_bipush
		code.Codes[code.CodeLength+1] = byte(e.Data.(int64))
		code.CodeLength += 2
		maxStack = 1
	case ast.ExpressionTypeInt, ast.ExpressionTypeShort, ast.ExpressionTypeChar:
		loadInt32(class, code, int32(e.Data.(int64)))
		maxStack = 1
	case ast.ExpressionTypeLong:
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
	case ast.ExpressionTypeFloat:
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
	case ast.ExpressionTypeDouble:
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
	case ast.ExpressionTypeString:
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst(e.Data.(string), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if len([]byte(e.Data.(string))) > 65536 {
			panic("jvm max string length is 65536")
		}
		maxStack = 1
	//binary expression
	case ast.ExpressionTypeLogicalOr:
		fallthrough
	case ast.ExpressionTypeLogicalAnd:
		maxStack = this.buildLogical(class, code, e, context, state)
	case ast.ExpressionTypeOr:
		fallthrough
	case ast.ExpressionTypeAnd:
		fallthrough
	case ast.ExpressionTypeXor:
		fallthrough
	case ast.ExpressionTypeLsh:
		fallthrough
	case ast.ExpressionTypeRsh:
		fallthrough
	case ast.ExpressionTypeAdd:
		fallthrough
	case ast.ExpressionTypeSub:
		fallthrough
	case ast.ExpressionTypeMul:
		fallthrough
	case ast.ExpressionTypeDiv:
		fallthrough
	case ast.ExpressionTypeMod:
		maxStack = this.buildArithmetic(class, code, e, context, state)
	//
	case ast.ExpressionTypeAssign:
		maxStack = this.buildAssign(class, code, e, context, state)
	case ast.ExpressionTypeVarAssign:
		maxStack = this.buildVarAssign(class, code, e, context, state)
	//
	case ast.ExpressionTypePlusAssign:
		fallthrough
	case ast.ExpressionTypeMinusAssign:
		fallthrough
	case ast.ExpressionTypeMulAssign:
		fallthrough
	case ast.ExpressionTypeDivAssign:
		fallthrough
	case ast.ExpressionTypeModAssign:
		fallthrough
	case ast.ExpressionTypeAndAssign:
		fallthrough
	case ast.ExpressionTypeOrAssign:
		fallthrough
	case ast.ExpressionTypeLshAssign:
		fallthrough
	case ast.ExpressionTypeRshAssign:
		fallthrough
	case ast.ExpressionTypeXorAssign:
		maxStack = this.buildOpAssign(class, code, e, context, state)
	//
	case ast.ExpressionTypeEq:
		fallthrough
	case ast.ExpressionTypeNe:
		fallthrough
	case ast.ExpressionTypeGe:
		fallthrough
	case ast.ExpressionTypeGt:
		fallthrough
	case ast.ExpressionTypeLe:
		fallthrough
	case ast.ExpressionTypeLt:
		maxStack = this.buildRelations(class, code, e, context, state)
	//
	case ast.ExpressionTypeIndex:
		maxStack = this.buildIndex(class, code, e, context, state)
	case ast.ExpressionTypeSelection:
		maxStack = this.buildSelection(class, code, e, context, state)
	//
	case ast.ExpressionTypeMethodCall:
		maxStack = this.buildMethodCall(class, code, e, context, state)
	case ast.ExpressionTypeFunctionCall:
		maxStack = this.buildFunctionCall(class, code, e, context, state)
	//
	case ast.ExpressionTypeIncrement:
		fallthrough
	case ast.ExpressionTypeDecrement:
		fallthrough
	case ast.ExpressionTypePrefixIncrement:
		fallthrough
	case ast.ExpressionTypePrefixDecrement:
		maxStack = this.buildSelfIncrement(class, code, e, context, state)
	//
	case ast.ExpressionTypeBitwiseNot:
		fallthrough
	case ast.ExpressionTypeNegative:
		fallthrough
	case ast.ExpressionTypeNot:
		maxStack = this.buildUnary(class, code, e, context, state)
	//
	case ast.ExpressionTypeIdentifier:
		maxStack = this.buildIdentifier(class, code, e, context)
	case ast.ExpressionTypeNew:
		maxStack = this.buildNew(class, code, e, context, state)
	case ast.ExpressionTypeFunctionLiteral:
		maxStack = this.BuildPackage.buildFunctionExpression(class, code, e, context, state)
	case ast.ExpressionTypeCheckCast: // []byte(str)
		maxStack = this.buildTypeConversion(class, code, e, context, state)
	case ast.ExpressionTypeConst:
		/*
		 analyse at ast stage
		*/
	case ast.ExpressionTypeSlice:
		maxStack = this.buildSlice(class, code, e, context, state)
	case ast.ExpressionTypeArray:
		maxStack = this.buildArray(class, code, e, context, state)
	case ast.ExpressionTypeMap:
		maxStack = this.buildMapLiteral(class, code, e, context, state)
	case ast.ExpressionTypeVar:
		maxStack = this.buildVar(class, code, e, context, state)
	case ast.ExpressionTypeTypeAssert:
		maxStack = this.buildTypeAssert(class, code, e, context, state)
	case ast.ExpressionTypeQuestion:
		maxStack = this.buildQuestion(class, code, e, context, state)
	case ast.ExpressionTypeVArgs:
		maxStack = this.build(class, code, e.Data.(*ast.Expression), context, state)
	default:
		panic("missing handle:" + e.Op)
	}
	return
}

func (this *BuildExpression) jvmSize(e *ast.Expression) (size uint16) {
	if len(e.MultiValues) > 1 {
		return 1
	}
	return jvmSlotSize(e.Value)
}

func (this *BuildExpression) buildExpressions(class *cg.ClassHighLevel, code *cg.AttributeCode,
	es []*ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := 0
	for _, e := range es {
		if e.HaveMultiValue() {
			length += len(e.MultiValues)
		} else {
			length++
		}
	}
	loadInt32(class, code, int32(length))
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst(javaRootClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if 1 > maxStack {
		maxStack = 1
	}
	arrayListObject := state.newObjectVariableType(javaRootObjectArray)
	state.pushStack(class, arrayListObject)
	defer state.popStack(1)
	index := int32(0)
	for _, v := range es {
		currentStack := uint16(1)
		if v.HaveMultiValue() {
			stack := this.build(class, code, v, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			autoVar := newMultiValueAutoVar(class, code, state)
			for kk, _ := range v.MultiValues {
				currentStack = 1
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				currentStack++
				stack = autoVar.unPack2Object(class, code, kk)
				if t := stack + currentStack; t > maxStack {
					maxStack = t
				}
				loadInt32(class, code, index)
				if 4 > maxStack { // current stack is  arrayRef arrayRef value index
					maxStack = 4
				}
				code.Codes[code.CodeLength] = cg.OP_swap
				code.Codes[code.CodeLength+1] = cg.OP_aastore
				code.CodeLength += 2
				index++
			}
			continue
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		state.pushStack(class, arrayListObject)
		currentStack++
		stack := this.build(class, code, v, context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Value.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.Value)
		}
		loadInt32(class, code, index)
		if 4 > maxStack { // current stack is  arrayRef arrayRef value index
			maxStack = 4
		}
		code.Codes[code.CodeLength] = cg.OP_swap
		code.Codes[code.CodeLength+1] = cg.OP_aastore
		code.CodeLength += 2
		state.popStack(1) // @270
		index++
	}
	return
}
