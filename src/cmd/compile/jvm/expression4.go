package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildRelations(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 4
	bin := e.Data.(*ast.ExpressionBinary)
	if bin.Left.VariableType.IsNumber() { // number types
		stack, _ := m.build(class, code, bin.Left, context)
		if stack > maxstack {
			maxstack = stack
		}
		target := bin.Left.VariableType.NumberTypeConvertRule(bin.Right.VariableType)
		if target == ast.VARIABLE_TYPE_INT || target == ast.VARIABLE_TYPE_SHORT || target == ast.VARIABLE_TYPE_BYTE {
			target = ast.VARIABLE_TYPE_LONG
		}
		if target != bin.Left.VariableType.Typ {
			m.numberTypeConverter(code, bin.Left.VariableType.Typ, target)
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
		}
		if target != bin.Right.VariableType.Typ {
			m.numberTypeConverter(code, bin.Right.VariableType.Typ, target)
		}
		switch target {
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lcmp
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_fcmpl
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dcmpl
		}
		code.CodeLength++
		if e.Typ == ast.EXPRESSION_TYPE_GT || e.Typ == ast.EXPRESSION_TYPE_LE { // > and <=
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength] = cg.OP_ifge
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifge
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+8)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		}
		if e.Typ == ast.EXPRESSION_TYPE_EQ || e.Typ == ast.EXPRESSION_TYPE_NE { // == and !=
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength] = cg.OP_ifeq
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifne
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+8)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		}
		if e.Typ == ast.EXPRESSION_TYPE_EQ || e.Typ == ast.EXPRESSION_TYPE_NE { // < and >=
			if e.Typ == ast.EXPRESSION_TYPE_LE {
				code.Codes[code.CodeLength] = cg.OP_iflt
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifge
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+8)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		}
		return
	}
	if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_BOOL { // bool type
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_if_icmpeq
		} else {
			code.Codes[code.CodeLength] = cg.OP_if_icmpne
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+8)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}
	if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_NULL || bin.Right.VariableType.Typ == ast.VARIABLE_TYPE_NULL {
		var stack uint16
		if bin.Left.VariableType.Typ != ast.VARIABLE_TYPE_NULL {
			stack, _ = m.build(class, code, bin.Left, context)
		} else {
			stack, _ = m.build(class, code, bin.Right, context)
		}
		if stack > maxstack {
			maxstack = stack
		}
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_ifnull
		} else { // ne
			code.Codes[code.CodeLength] = cg.OP_ifnonnull
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+8)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}
	if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_OBJECT || ast.VARIABLE_TYPE_ARRAY_INSTANCE == bin.Left.VariableType.Typ { //
		maxstack = uint16(1)
		stack, _ := m.build(class, code, bin.Left, context)
		if stack > maxstack {
			maxstack = stack
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if stack+1 > maxstack {
			maxstack = stack + 1
		}
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_if_acmpeq
		} else { // ne
			code.Codes[code.CodeLength] = cg.OP_if_acmpne
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+8)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}
	panic("missing")
	return
}

func (m *MakeExpression) buildArithmetic(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	maxstack = 2
	stack, _ := m.build(class, code, bin.Left, context)
	if stack > maxstack {
		maxstack = stack
	}
	if e.Typ == ast.EXPRESSION_TYPE_OR || e.Typ == ast.EXPRESSION_TYPE_AND {
		stack, _ := m.build(class, code, bin.Right, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
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
		default:
			panic("~~~~~~~~~~~~")
		}
		code.CodeLength++
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_ADD || e.Typ == ast.EXPRESSION_TYPE_SUB || e.Typ == ast.EXPRESSION_TYPE_MUL ||
		e.Typ == ast.EXPRESSION_TYPE_DIV || e.Typ == ast.EXPRESSION_TYPE_MOD {
		if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_STRING {
			panic(1)
		}
		if e.VariableType.Typ != bin.Left.Typ {
			m.numberTypeConverter(code, bin.Left.Typ, e.VariableType.Typ)
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
		}
		if e.VariableType.Typ != bin.Right.VariableType.Typ {
			m.numberTypeConverter(code, bin.Right.Typ, e.VariableType.Typ)
		}
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
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
		default:
			panic("~~~~~~~~~~~~")
		}
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == ast.EXPRESSION_TYPE_RIGHT_SHIFT {
		stack, _ := m.build(class, code, bin.Right, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
		}
		if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_LONG { // long
			if bin.Right.VariableType.Typ != ast.VARIABLE_TYPE_LONG {
				m.stackTop2Long(code, bin.Right.VariableType.Typ)
			}
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
				code.Codes[code.CodeLength] = cg.OP_lshl
			} else {
				code.Codes[code.CodeLength] = cg.OP_lshr
			}
		} else { // int
			if bin.Right.VariableType.Typ != ast.VARIABLE_TYPE_INT {
				m.stackTop2Int(code, bin.Right.VariableType.Typ)
			}
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
		}
		code.CodeLength++
		return
	}
	return
}

func (m *MakeExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16, exits [][]byte) {
	maxstack = 1
	exits = [][]byte{}
	bin := e.Data.(*ast.ExpressionBinary)
	var stack uint16
	var exits1 [][]byte
	stack, exits1 = m.build(class, code, bin.Left, context)
	if stack > maxstack {
		maxstack = stack
	}
	exits = append(exits, exits1...)
	if e.Typ == ast.EXPRESSION_TYPE_LOGICAL_OR {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifne
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		stack, exits1 = m.build(class, code, bin.Left, context)
		if stack+1 > maxstack {
			maxstack = stack + 1
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifne
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
	} else { //and
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifeq
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		stack, exits1 = m.build(class, code, bin.Left, context)
		if stack+1 > maxstack {
			maxstack = stack + 1
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifeq
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
	}
	return
}
