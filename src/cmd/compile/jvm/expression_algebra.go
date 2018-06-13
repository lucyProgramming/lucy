package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildArithmetic(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if e.Type == ast.EXPRESSION_TYPE_OR ||
		e.Type == ast.EXPRESSION_TYPE_AND ||
		e.Type == ast.EXPRESSION_TYPE_XOR {
		maxStack, _ = makeExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		stack, _ := makeExpression.build(class, code, bin.Right, context, state)
		if t := stack + jvmSize(bin.Left.Value); t > maxStack {
			maxStack = t
		}
		switch e.Value.Type {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			fallthrough
		case ast.VARIABLE_TYPE_FLOAT:
			if e.Type == ast.EXPRESSION_TYPE_AND {
				code.Codes[code.CodeLength] = cg.OP_iand
			} else if e.Type == ast.EXPRESSION_TYPE_OR {
				code.Codes[code.CodeLength] = cg.OP_ior
			} else {
				code.Codes[code.CodeLength] = cg.OP_ixor
			}
		case ast.VARIABLE_TYPE_DOUBLE:
			fallthrough
		case ast.VARIABLE_TYPE_LONG:
			if e.Type == ast.EXPRESSION_TYPE_AND {
				code.Codes[code.CodeLength] = cg.OP_land
			} else if e.Type == ast.EXPRESSION_TYPE_OR {
				code.Codes[code.CodeLength] = cg.OP_lor
			} else {
				code.Codes[code.CodeLength] = cg.OP_lxor
			}
		}
		code.CodeLength++
		return
	}
	if e.Type == ast.EXPRESSION_TYPE_ADD ||
		e.Type == ast.EXPRESSION_TYPE_SUB ||
		e.Type == ast.EXPRESSION_TYPE_MUL ||
		e.Type == ast.EXPRESSION_TYPE_DIV ||
		e.Type == ast.EXPRESSION_TYPE_MOD {
		//handle string first
		if bin.Left.Value.Type == ast.VARIABLE_TYPE_STRING ||
			bin.Right.Value.Type == ast.VARIABLE_TYPE_STRING {
			return makeExpression.buildStrCat(class, code, bin, context, state)
		}
		maxStack, _ = makeExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, e.Value)
		stack, _ := makeExpression.build(class, code, bin.Right, context, state)
		if t := jvmSize(bin.Left.Value) + stack; t > maxStack {
			maxStack = t
		}
		switch e.Value.Type {
		case ast.VARIABLE_TYPE_BYTE:
			switch e.Type {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_iadd
				code.Codes[code.CodeLength+1] = cg.OP_i2b
				code.CodeLength += 2
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_isub
				code.Codes[code.CodeLength+1] = cg.OP_i2b
				code.CodeLength += 2
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_imul
				code.Codes[code.CodeLength+1] = cg.OP_i2b
				code.CodeLength += 2
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_idiv
				code.CodeLength++
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_irem
				code.CodeLength++
			}

		case ast.VARIABLE_TYPE_SHORT:
			switch e.Type {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_iadd
				code.Codes[code.CodeLength+1] = cg.OP_i2s
				code.CodeLength += 2
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_isub
				code.Codes[code.CodeLength+1] = cg.OP_i2s
				code.CodeLength += 2
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_imul
				code.Codes[code.CodeLength+1] = cg.OP_i2s
				code.CodeLength += 2
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_idiv
				code.CodeLength++
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_irem
				code.CodeLength++
			}
		case ast.VARIABLE_TYPE_INT:
			switch e.Type {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_iadd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_isub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_imul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_idiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_irem
			}
			code.CodeLength++
		case ast.VARIABLE_TYPE_FLOAT:
			switch e.Type {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_fadd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_fsub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_fmul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_fdiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_frem
			}
			code.CodeLength++
		case ast.VARIABLE_TYPE_DOUBLE:
			switch e.Type {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_dadd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_dsub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_dmul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_ddiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_drem
			}
			code.CodeLength++
		case ast.VARIABLE_TYPE_LONG:
			switch e.Type {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_ladd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_lsub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_lmul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_ldiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_lrem
			}
			code.CodeLength++
		}
		return
	}

	if e.Type == ast.EXPRESSION_TYPE_LSH ||
		e.Type == ast.EXPRESSION_TYPE_RSH {
		maxStack, _ = makeExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		stack, _ := makeExpression.build(class, code, bin.Right, context, state)
		if t := stack + jvmSize(bin.Left.Value); t > maxStack {
			maxStack = t
		}
		switch e.Value.Type {
		case ast.VARIABLE_TYPE_BYTE:
			if e.Type == ast.EXPRESSION_TYPE_LSH {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_SHORT:
			if e.Type == ast.EXPRESSION_TYPE_LSH {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_INT:
			if e.Type == ast.EXPRESSION_TYPE_LSH {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.CodeLength++
		case ast.VARIABLE_TYPE_LONG:
			if e.Type == ast.EXPRESSION_TYPE_LSH {
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
