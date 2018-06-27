package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildUnary(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {

	if e.Type == ast.EXPRESSION_TYPE_NEGATIVE {
		maxStack, _ = makeExpression.build(class, code, e.Data.(*ast.Expression), context, state)
		switch e.ExpressionValue.Type {
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_ineg
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_fneg
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_dneg
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_lneg
		}
		code.CodeLength++
		return
	}
	if e.Type == ast.EXPRESSION_TYPE_BIT_NOT {
		ee := e.Data.(*ast.Expression)
		maxStack, _ = makeExpression.build(class, code, ee, context, state)
		if t := jvmSlotSize(ee.ExpressionValue) * 2; t > maxStack {
			maxStack = t
		}
		switch e.ExpressionValue.Type {
		case ast.VariableTypeByte:
			code.Codes[code.CodeLength] = cg.OP_bipush
			code.Codes[code.CodeLength+1] = 255
			code.Codes[code.CodeLength+2] = cg.OP_ixor
			code.CodeLength += 3
			if 2 > maxStack {
				maxStack = 2
			}
		case ast.VariableTypeShort:
			code.Codes[code.CodeLength] = cg.OP_sipush
			code.Codes[code.CodeLength+1] = 255
			code.Codes[code.CodeLength+2] = 255
			code.Codes[code.CodeLength+3] = cg.OP_ixor
			code.CodeLength += 4
			if 2 > maxStack {
				maxStack = 2
			}
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertIntConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = cg.OP_ixor
			code.CodeLength += 4
			if 2 > maxStack {
				maxStack = 2
			}
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_ldc2_w
			class.InsertLongConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = cg.OP_lxor
			code.CodeLength += 4
			if 4 > maxStack {
				maxStack = 4
			}
		}
		return
	}
	if e.Type == ast.EXPRESSION_TYPE_NOT {
		ee := e.Data.(*ast.Expression)
		var es []*cg.Exit
		maxStack, es = makeExpression.build(class, code, ee, context, state)
		if len(es) > 0 {
			fillOffsetForExits(es, code.CodeLength)
			state.pushStack(class, ee.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, ee.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength+8)
		state.popStack(1)
		code.Codes[code.CodeLength] = cg.OP_ifeq
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], uint16(7))
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], uint16(4))
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
	}
	return
}
