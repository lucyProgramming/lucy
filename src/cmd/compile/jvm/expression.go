package jvm

import (
	"encoding/binary"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/common"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type MakeExpression struct {
}

func (m *MakeExpression) build(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16, exits [][]byte) {
	exits = [][]byte{}
	switch e.Typ {
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
		e.Data = int64(e.Data.(byte))
		fallthrough
	case ast.EXPRESSION_TYPE_INT:
		value := e.Data.(int64)
		if value == 0 {
			code.Codes[code.CodeLength] = cg.OP_iconst_0
			code.CodeLength += 1
		} else if value == 1 {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
			code.CodeLength += 1
		} else if value == 2 {
			code.Codes[code.CodeLength] = cg.OP_iconst_2
			code.CodeLength += 1
		} else if value == 3 {
			code.Codes[code.CodeLength] = cg.OP_iconst_3
			code.CodeLength += 1
		} else if value == 4 {
			code.Codes[code.CodeLength] = cg.OP_iconst_4
			code.CodeLength += 1
		} else if value == 5 {
			code.Codes[code.CodeLength] = cg.OP_iconst_5
			code.CodeLength += 1
		} else if -127 >= value && value <= 128 {
			code.Codes[code.CodeLength] = cg.OP_bipush
			code.Codes[code.CodeLength+1] = byte(value)
			code.CodeLength += 2
		} else if -32768 <= value && 32767 >= value {
			code.Codes[code.CodeLength] = cg.OP_sipush
			code.Codes[code.CodeLength+1] = byte(int16(value) >> 8)
			code.Codes[code.CodeLength+2] = byte(value)
			code.CodeLength += 3
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertIntConst(int32(value), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		maxstack = 1
	case ast.EXPRESSION_TYPE_FLOAT:
		if common.Float64Equal(e.Data.(float64), 0.0) {
			code.Codes[code.CodeLength] = cg.OP_fconst_0
			code.CodeLength++
		} else if common.Float64Equal(e.Data.(float64), 1.0) {
			code.Codes[code.CodeLength] = cg.OP_fconst_1
			code.CodeLength++
		} else if common.Float64Equal(e.Data.(float64), 2.0) {
			code.Codes[code.CodeLength] = cg.OP_fconst_2
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertFloatConst(float32(e.Data.(float64)), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	case ast.EXPRESSION_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst(e.Data.(string), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxstack = 1
	case ast.EXPRESSION_TYPE_ARRAY: // []bool{false,true}
		panic("11")
	//binary expression
	case ast.EXPRESSION_TYPE_LOGICAL_OR:
		fallthrough
	case ast.EXPRESSION_TYPE_LOGICAL_AND:
		maxstack, exits = m.buildLogical(class, code, e, context)
	case ast.EXPRESSION_TYPE_OR:
		fallthrough
	case ast.EXPRESSION_TYPE_AND:
		fallthrough
	case ast.EXPRESSION_TYPE_LEFT_SHIFT:
		fallthrough
	case ast.EXPRESSION_TYPE_RIGHT_SHIFT:
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
		maxstack = m.buildArithmetic(class, code, e, context)
	//
	case ast.EXPRESSION_TYPE_ASSIGN:
	case ast.EXPRESSION_TYPE_COLON_ASSIGN:

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
		maxstack = m.buildOpAssign(class, code, e, context)
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
		maxstack = m.buildRelations(class, code, e, context)
	//
	case ast.EXPRESSION_TYPE_INDEX:
		maxstack = m.buildIndex(class, code, e, context)
	case ast.EXPRESSION_TYPE_DOT:
		maxstack = m.buildDot(class, code, e, context)

	//
	case ast.EXPRESSION_TYPE_METHOD_CALL:
	case ast.EXPRESSION_TYPE_FUNCTION_CALL:
	//
	case ast.EXPRESSION_TYPE_INCREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_DECREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_PRE_INCREMENT:
		fallthrough
	case ast.EXPRESSION_TYPE_PRE_DECREMENT:
		maxstack = m.buildSelfIncrement(class, code, e, context)
	//
	case ast.EXPRESSION_TYPE_NEGATIVE:
		fallthrough
	case ast.EXPRESSION_TYPE_NOT:
		maxstack = m.buildUnary(class, code, e, context)
	//
	case ast.EXPRESSION_TYPE_IDENTIFIER:
		maxstack = m.buildIdentifer(class, code, e, context)
	case ast.EXPRESSION_TYPE_NEW:
		maxstack = m.buildNew(class, code, e, context)
	case ast.EXPRESSION_TYPE_LIST:
	case ast.EXPRESSION_TYPE_FUNCTION:
	case ast.EXPRESSION_TYPE_VAR:
	case ast.EXPRESSION_TYPE_CONVERTION_TYPE: // []byte(str)
	}
	return
}

func (m *MakeExpression) buildOpAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	//maxstack, _, op, target := m.buildLeftValue(class, code, e, context)
	return
}
func (m *MakeExpression) buildSelfIncrement(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	//ee := e.Data.(*ast.Expression)
	// m.leftValue(class, code, e)
	return
}

func (m *MakeExpression) buildUnary(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 2
	maxstack1, es := m.build(class, code, e.Data.(*ast.Expression), context) // in case !(xxx && a())
	if maxstack1 > maxstack {
		maxstack = maxstack1
	}
	backPatchEs(es, code)
	if e.Typ == ast.EXPRESSION_TYPE_NEGATIVE {
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_ineg
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_fneg
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dneg
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lneg
		}
		code.CodeLength++
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_NOT {
		code.Codes[code.CodeLength] = cg.OP_ifne                                      // length 1
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7) // length 2
		code.Codes[code.CodeLength+3] = cg.OP_iconst_1                                // length 1
		code.Codes[code.CodeLength+4] = cg.OP_goto                                    // length 1
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+8) // length 2
		code.Codes[code.CodeLength+7] = cg.OP_iconst_0                                // length 1
		code.CodeLength += 8
		return
	}
	return
}
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
			m.primitiveConverter(code, bin.Left.VariableType.Typ, target)
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
		}
		if target != bin.Right.VariableType.Typ {
			m.primitiveConverter(code, bin.Right.VariableType.Typ, target)
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
			m.primitiveConverter(code, bin.Left.Typ, e.VariableType.Typ)
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
		}
		if e.VariableType.Typ != bin.Right.VariableType.Typ {
			m.primitiveConverter(code, bin.Right.Typ, e.VariableType.Typ)
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

func (m *MakeExpression) stackTop2Byte(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_i2b
		code.CodeLength++
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2b
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_l2i
		code.CodeLength += 2
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Short(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
	case ast.VARIABLE_TYPE_SHORT:
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2s
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Int(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
	case ast.VARIABLE_TYPE_SHORT:
	case ast.VARIABLE_TYPE_INT:
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.CodeLength++
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.CodeLength++
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Float(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2f
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2f
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2f
		code.CodeLength++
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Long(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2l
		code.CodeLength++

	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2l
		code.CodeLength++

	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2l
		code.CodeLength++

	case ast.VARIABLE_TYPE_LONG:
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Double(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2d
		code.CodeLength++

	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2d
		code.CodeLength++

	case ast.VARIABLE_TYPE_DOUBLE:
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2d
		code.CodeLength++

	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) primitiveConverter(code *cg.AttributeCode, typ int, target int) {
	if typ == target {
		return
	}
	switch target {
	case ast.VARIABLE_TYPE_BYTE:
		m.stackTop2Byte(code, typ)
	case ast.VARIABLE_TYPE_SHORT:
		m.stackTop2Short(code, typ)
	case ast.VARIABLE_TYPE_INT:
		m.stackTop2Int(code, typ)
	case ast.VARIABLE_TYPE_LONG:
		m.stackTop2Long(code, typ)
	case ast.VARIABLE_TYPE_FLOAT:
		m.stackTop2Float(code, typ)
	case ast.VARIABLE_TYPE_DOUBLE:
		m.stackTop2Double(code, typ)
	default:
		panic(1)
	}
}
