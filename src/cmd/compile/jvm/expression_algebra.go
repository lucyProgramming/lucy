package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildArithmetic(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	if e.Typ == ast.EXPRESSION_TYPE_OR || e.Typ == ast.EXPRESSION_TYPE_AND {
		maxstack, _ = m.build(class, code, bin.Left, context)
		size := bin.Left.VariableType.JvmSlotSize()
		stack, _ := m.build(class, code, bin.Right, context)
		if t := stack + size; t > maxstack {
			maxstack = t
		}
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			fallthrough
		case ast.VARIABLE_TYPE_FLOAT:
			if e.Typ == ast.EXPRESSION_TYPE_AND {
				code.Codes[code.CodeLength] = cg.OP_iand
			} else {
				code.Codes[code.CodeLength] = cg.OP_ior
			}
		case ast.VARIABLE_TYPE_DOUBLE:
			fallthrough
		case ast.VARIABLE_TYPE_LONG:
			if e.Typ == ast.EXPRESSION_TYPE_AND {
				code.Codes[code.CodeLength] = cg.OP_land
			} else {
				code.Codes[code.CodeLength] = cg.OP_lor
			}
		}
		code.CodeLength++
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_ADD || e.Typ == ast.EXPRESSION_TYPE_SUB ||
		e.Typ == ast.EXPRESSION_TYPE_MUL || e.Typ == ast.EXPRESSION_TYPE_DIV ||
		e.Typ == ast.EXPRESSION_TYPE_MOD {
		//handle string first
		if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_STRING || bin.Right.VariableType.Typ == ast.VARIABLE_TYPE_STRING {
			return m.buildStrCat(class, code, bin, context)
		}
		maxstack = 4
		stack, _ := m.build(class, code, bin.Left, context)
		if stack > maxstack {
			maxstack = stack
		}
		//number type,no doubt
		if e.VariableType.Typ != bin.Left.VariableType.Typ {
			m.numberTypeConverter(code, bin.Left.VariableType.Typ, e.VariableType.Typ)
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if t := 2 + stack; t > maxstack {
			maxstack = t
		}
		if e.VariableType.Typ != bin.Right.VariableType.Typ {
			m.numberTypeConverter(code, bin.Right.VariableType.Typ, e.VariableType.Typ)
		}
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			switch e.Typ {
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
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_SHORT:
			switch e.Typ {
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
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_INT:
			switch e.Typ {
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
			switch e.Typ {
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
			switch e.Typ {
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
			switch e.Typ {
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

	if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == ast.EXPRESSION_TYPE_RIGHT_SHIFT {
		maxstack, _ = m.build(class, code, bin.Left, context)
		if e.VariableType.Typ != bin.Left.VariableType.Typ {
			m.numberTypeConverter(code, bin.Left.Typ, e.VariableType.Typ)
		}
		size := e.VariableType.JvmSlotSize()
		stack, _ := m.build(class, code, bin.Right, context)
		if t := stack + size; t > maxstack {
			maxstack = t
		}
		if e.VariableType.Typ != bin.Right.VariableType.Typ {
			m.numberTypeConverter(code, bin.Right.Typ, e.VariableType.Typ)
		}
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_SHORT:
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		case ast.VARIABLE_TYPE_INT:
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
			code.CodeLength++
		case ast.VARIABLE_TYPE_LONG:
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
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
