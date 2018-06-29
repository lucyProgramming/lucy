package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildArithmetic(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if e.Type == ast.ExpressionTypeOr ||
		e.Type == ast.ExpressionTypeAnd ||
		e.Type == ast.ExpressionTypeXor {
		maxStack, _ = buildExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.ExpressionValue)
		stack, _ := buildExpression.build(class, code, bin.Right, context, state)
		if t := stack + jvmSlotSize(bin.Left.ExpressionValue); t > maxStack {
			maxStack = t
		}
		switch e.ExpressionValue.Type {
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeInt:
			fallthrough
		case ast.VariableTypeFloat:
			if e.Type == ast.ExpressionTypeAnd {
				code.Codes[code.CodeLength] = cg.OP_iand
			} else if e.Type == ast.ExpressionTypeOr {
				code.Codes[code.CodeLength] = cg.OP_ior
			} else {
				code.Codes[code.CodeLength] = cg.OP_ixor
			}
		case ast.VariableTypeDouble:
			fallthrough
		case ast.VariableTypeLong:
			if e.Type == ast.ExpressionTypeAnd {
				code.Codes[code.CodeLength] = cg.OP_land
			} else if e.Type == ast.ExpressionTypeOr {
				code.Codes[code.CodeLength] = cg.OP_lor
			} else {
				code.Codes[code.CodeLength] = cg.OP_lxor
			}
		}
		code.CodeLength++
		return
	}
	if e.Type == ast.ExpressionTypeAdd ||
		e.Type == ast.ExpressionTypeSub ||
		e.Type == ast.ExpressionTypeMul ||
		e.Type == ast.ExpressionTypeDiv ||
		e.Type == ast.ExpressionTypeMod {
		//handle string first
		if bin.Left.ExpressionValue.Type == ast.VariableTypeString ||
			bin.Right.ExpressionValue.Type == ast.VariableTypeString {
			return buildExpression.buildStrCat(class, code, bin, context, state)
		}
		maxStack, _ = buildExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, e.ExpressionValue)
		stack, _ := buildExpression.build(class, code, bin.Right, context, state)
		if t := jvmSlotSize(bin.Left.ExpressionValue) + stack; t > maxStack {
			maxStack = t
		}
		switch e.ExpressionValue.Type {
		case ast.VariableTypeByte:
			switch e.Type {
			case ast.ExpressionTypeAdd:
				code.Codes[code.CodeLength] = cg.OP_iadd
				code.Codes[code.CodeLength+1] = cg.OP_i2b
				code.CodeLength += 2
			case ast.ExpressionTypeSub:
				code.Codes[code.CodeLength] = cg.OP_isub
				code.Codes[code.CodeLength+1] = cg.OP_i2b
				code.CodeLength += 2
			case ast.ExpressionTypeMul:
				code.Codes[code.CodeLength] = cg.OP_imul
				code.Codes[code.CodeLength+1] = cg.OP_i2b
				code.CodeLength += 2
			case ast.ExpressionTypeDiv:
				code.Codes[code.CodeLength] = cg.OP_idiv
				code.CodeLength++
			case ast.ExpressionTypeMod:
				code.Codes[code.CodeLength] = cg.OP_irem
				code.CodeLength++
			}

		case ast.VariableTypeShort:
			switch e.Type {
			case ast.ExpressionTypeAdd:
				code.Codes[code.CodeLength] = cg.OP_iadd
				code.Codes[code.CodeLength+1] = cg.OP_i2s
				code.CodeLength += 2
			case ast.ExpressionTypeSub:
				code.Codes[code.CodeLength] = cg.OP_isub
				code.Codes[code.CodeLength+1] = cg.OP_i2s
				code.CodeLength += 2
			case ast.ExpressionTypeMul:
				code.Codes[code.CodeLength] = cg.OP_imul
				code.Codes[code.CodeLength+1] = cg.OP_i2s
				code.CodeLength += 2
			case ast.ExpressionTypeDiv:
				code.Codes[code.CodeLength] = cg.OP_idiv
				code.CodeLength++
			case ast.ExpressionTypeMod:
				code.Codes[code.CodeLength] = cg.OP_irem
				code.CodeLength++
			}
		case ast.VariableTypeInt:
			switch e.Type {
			case ast.ExpressionTypeAdd:
				code.Codes[code.CodeLength] = cg.OP_iadd
			case ast.ExpressionTypeSub:
				code.Codes[code.CodeLength] = cg.OP_isub
			case ast.ExpressionTypeMul:
				code.Codes[code.CodeLength] = cg.OP_imul
			case ast.ExpressionTypeDiv:
				code.Codes[code.CodeLength] = cg.OP_idiv
			case ast.ExpressionTypeMod:
				code.Codes[code.CodeLength] = cg.OP_irem
			}
			code.CodeLength++
		case ast.VariableTypeFloat:
			switch e.Type {
			case ast.ExpressionTypeAdd:
				code.Codes[code.CodeLength] = cg.OP_fadd
			case ast.ExpressionTypeSub:
				code.Codes[code.CodeLength] = cg.OP_fsub
			case ast.ExpressionTypeMul:
				code.Codes[code.CodeLength] = cg.OP_fmul
			case ast.ExpressionTypeDiv:
				code.Codes[code.CodeLength] = cg.OP_fdiv
			case ast.ExpressionTypeMod:
				code.Codes[code.CodeLength] = cg.OP_frem
			}
			code.CodeLength++
		case ast.VariableTypeDouble:
			switch e.Type {
			case ast.ExpressionTypeAdd:
				code.Codes[code.CodeLength] = cg.OP_dadd
			case ast.ExpressionTypeSub:
				code.Codes[code.CodeLength] = cg.OP_dsub
			case ast.ExpressionTypeMul:
				code.Codes[code.CodeLength] = cg.OP_dmul
			case ast.ExpressionTypeDiv:
				code.Codes[code.CodeLength] = cg.OP_ddiv
			case ast.ExpressionTypeMod:
				code.Codes[code.CodeLength] = cg.OP_drem
			}
			code.CodeLength++
		case ast.VariableTypeLong:
			switch e.Type {
			case ast.ExpressionTypeAdd:
				code.Codes[code.CodeLength] = cg.OP_ladd
			case ast.ExpressionTypeSub:
				code.Codes[code.CodeLength] = cg.OP_lsub
			case ast.ExpressionTypeMul:
				code.Codes[code.CodeLength] = cg.OP_lmul
			case ast.ExpressionTypeDiv:
				code.Codes[code.CodeLength] = cg.OP_ldiv
			case ast.ExpressionTypeMod:
				code.Codes[code.CodeLength] = cg.OP_lrem
			}
			code.CodeLength++
		}
		return
	}

	if e.Type == ast.ExpressionTypeLsh ||
		e.Type == ast.ExpressionTypeRsh {
		maxStack, _ = buildExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.ExpressionValue)
		stack, _ := buildExpression.build(class, code, bin.Right, context, state)
		if t := stack + jvmSlotSize(bin.Left.ExpressionValue); t > maxStack {
			maxStack = t
		}
		switch e.ExpressionValue.Type {
		case ast.VariableTypeByte:
			if e.Type == ast.ExpressionTypeLsh {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		case ast.VariableTypeShort:
			if e.Type == ast.ExpressionTypeLsh {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		case ast.VariableTypeInt:
			if e.Type == ast.ExpressionTypeLsh {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.CodeLength++
		case ast.VariableTypeLong:
			if e.Type == ast.ExpressionTypeLsh {
				code.Codes[code.CodeLength] = cg.OP_lshl
			} else {
				code.Codes[code.CodeLength] = cg.OP_lshr
			}
			code.CodeLength++
		}
		return
	}
	return
}
