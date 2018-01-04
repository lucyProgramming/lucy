package jvm

import (
	"encoding/binary"
	"github.com/756445638/lucy/src/cmd/compile/ast"
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
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertIntConst(int32(value), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		maxstack = 1
	case ast.EXPRESSION_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertFloatConst(float32(e.Data.(float64)), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
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

	case ast.EXPRESSION_TYPE_DOT:
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

	case ast.EXPRESSION_TYPE_LIST:
	case ast.EXPRESSION_TYPE_FUNCTION:
	case ast.EXPRESSION_TYPE_VAR:
	case ast.EXPRESSION_TYPE_CONVERTION_TYPE: // []byte(str)
	}
	return
}

func (m *MakeExpression) buildDot(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	maxstack = 2
	stack, _ := m.build(class, code, index.Expression, context)
	if stack > maxstack {
		maxstack = stack
	}
	switch index.Expression.VariableType.Typ {
	case ast.VARIABLE_TYPE_OBJECT:

	}
	return
}
func (m *MakeExpression) buildIndex(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 2
	index := e.Data.(*ast.ExpressionIndex)
	stack, _ := m.build(class, code, index.Expression, context)
	if stack > maxstack {
		maxstack = stack
	}
	stack, _ = m.build(class, code, index.Expression, context)
	if stack+2 > maxstack {
		maxstack = stack + 2
	}
	switch e.VariableType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_baload
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_saload
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_iaload
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_laload
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_faload
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_daload
	case ast.VARIABLE_TYPE_STRING:
		panic(1)
	case ast.VARIABLE_TYPE_OBJECT:
		code.Codes[code.CodeLength] = cg.OP_aaload
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		code.Codes[code.CodeLength] = cg.OP_aaload
	}
	code.CodeLength++
	return
}

func (m *MakeExpression) buildNew(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 1
	return
}
func (m *MakeExpression) buildIdentifer(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifer)
	if identifier.Var.BeenCaptured {
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_aload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_aload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_aload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_aload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("local object var out of range")
		}
		code.Codes[code.CodeLength] = cg.OP_iconst_1
		switch identifier.Var.Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			code.Codes[code.CodeLength+1] = cg.OP_baload
		case ast.VARIABLE_TYPE_SHORT:
			code.Codes[code.CodeLength+1] = cg.OP_saload
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength+1] = cg.OP_iaload
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength+1] = cg.OP_faload
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength+1] = cg.OP_daload
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength+1] = cg.OP_laload
		case ast.VARIABLE_TYPE_OBJECT:
			code.Codes[code.CodeLength+1] = cg.OP_aaload
		}
		code.CodeLength += 2
	} else {
		switch identifier.Var.Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			if identifier.Var.LocalValOffset == 0 {
				code.Codes[code.CodeLength] = cg.OP_iload_0
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 1 {
				code.Codes[code.CodeLength] = cg.OP_iload_1
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 2 {
				code.Codes[code.CodeLength] = cg.OP_iload_2
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 3 {
				code.Codes[code.CodeLength] = cg.OP_iload_3
				code.CodeLength++
			} else if identifier.Var.LocalValOffset < 255 {
				code.Codes[code.CodeLength] = cg.OP_iload
				code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
				code.CodeLength += 2
			} else {
				panic("local int var out of range")
			}
		case ast.VARIABLE_TYPE_FLOAT:
			if identifier.Var.LocalValOffset == 0 {
				code.Codes[code.CodeLength] = cg.OP_fload_0
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 1 {
				code.Codes[code.CodeLength] = cg.OP_fload_1
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 2 {
				code.Codes[code.CodeLength] = cg.OP_fload_2
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 3 {
				code.Codes[code.CodeLength] = cg.OP_fload_3
				code.CodeLength++
			} else if identifier.Var.LocalValOffset < 255 {
				code.Codes[code.CodeLength] = cg.OP_fload
				code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
				code.CodeLength += 2
			} else {
				panic("local float var out of range")
			}
		case ast.VARIABLE_TYPE_DOUBLE:
			if identifier.Var.LocalValOffset == 0 {
				code.Codes[code.CodeLength] = cg.OP_dload_0
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 1 {
				code.Codes[code.CodeLength] = cg.OP_dload_1
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 2 {
				code.Codes[code.CodeLength] = cg.OP_dload_2
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 3 {
				code.Codes[code.CodeLength] = cg.OP_dload_3
				code.CodeLength++
			} else if identifier.Var.LocalValOffset < 255 {
				code.Codes[code.CodeLength] = cg.OP_dload
				code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
				code.CodeLength += 2
			} else {
				panic("local double var out of range")
			}
			maxstack = 2
		case ast.VARIABLE_TYPE_LONG:
			if identifier.Var.LocalValOffset == 0 {
				code.Codes[code.CodeLength] = cg.OP_lload_0
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 1 {
				code.Codes[code.CodeLength] = cg.OP_lload_1
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 2 {
				code.Codes[code.CodeLength] = cg.OP_lload_2
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 3 {
				code.Codes[code.CodeLength] = cg.OP_lload_3
				code.CodeLength++
			} else if identifier.Var.LocalValOffset < 255 {
				code.Codes[code.CodeLength] = cg.OP_lload
				code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
				code.CodeLength += 2
			} else {
				panic("local double var out of range")
			}
			maxstack = 2
		case ast.VARIABLE_TYPE_OBJECT:
			if identifier.Var.LocalValOffset == 0 {
				code.Codes[code.CodeLength] = cg.OP_aload_0
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 1 {
				code.Codes[code.CodeLength] = cg.OP_aload_1
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 2 {
				code.Codes[code.CodeLength] = cg.OP_aload_2
				code.CodeLength++
			} else if identifier.Var.LocalValOffset == 3 {
				code.Codes[code.CodeLength] = cg.OP_aload_3
				code.CodeLength++
			} else if identifier.Var.LocalValOffset < 255 {
				code.Codes[code.CodeLength] = cg.OP_aload
				code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
				code.CodeLength += 2
			} else {
				panic("local object var out of range")
			}
		}
	}
	return
}

func (m *MakeExpression) buildLeftValue(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16, left_value_type int, op []byte, target *ast.VariableType) {
	switch e.Typ {
	case ast.EXPRESSION_TYPE_IDENTIFIER:
		identifier := e.Data.(*ast.ExpressionIdentifer)
		if identifier.Var.BeenCaptured {
			panic(1)
		} else {
			switch identifier.Var.Typ.Typ {
			case ast.VARIABLE_TYPE_BOOL:
				fallthrough
			case ast.VARIABLE_TYPE_BYTE:
				fallthrough
			case ast.VARIABLE_TYPE_SHORT:
				fallthrough
			case ast.VARIABLE_TYPE_INT:
				if identifier.Var.LocalValOffset == 0 {
					op = []byte{cg.OP_istore_0}
				} else if identifier.Var.LocalValOffset == 1 {
					op = []byte{cg.OP_istore_1}
				} else if identifier.Var.LocalValOffset == 2 {
					op = []byte{cg.OP_istore_2}
				} else if identifier.Var.LocalValOffset == 3 {
					op = []byte{cg.OP_istore_3}
				} else if identifier.Var.LocalValOffset <= 255 {
					op = []byte{cg.OP_istore, byte(identifier.Var.LocalValOffset)}
				} else {
					panic("local int var out of range")
				}
			case ast.VARIABLE_TYPE_FLOAT:
				if identifier.Var.LocalValOffset == 0 {
					op = []byte{cg.OP_fstore_0}
				} else if identifier.Var.LocalValOffset == 1 {
					op = []byte{cg.OP_fstore_1}
				} else if identifier.Var.LocalValOffset == 2 {
					op = []byte{cg.OP_fstore_2}
				} else if identifier.Var.LocalValOffset == 3 {
					op = []byte{cg.OP_fstore_3}
				} else if identifier.Var.LocalValOffset <= 255 {
					op = []byte{cg.OP_fstore, identifier.Var.LocalValOffset}
				} else {
					panic("local float var out of range")
				}
			case ast.VARIABLE_TYPE_DOUBLE:
				if identifier.Var.LocalValOffset == 0 {
					op = []byte{cg.OP_dstore_0}
				} else if identifier.Var.LocalValOffset == 1 {
					op = []byte{cg.OP_dstore_1}
				} else if identifier.Var.LocalValOffset == 2 {
					op = []byte{cg.OP_dstore_2}
				} else if identifier.Var.LocalValOffset == 3 {
					op = []byte{cg.OP_dstore_3}
				} else if identifier.Var.LocalValOffset <= 255 {
					op = []byte{cg.OP_dstore, identifier.Var.LocalValOffset}
				} else {
					panic("local float var out of range")
				}
			case ast.VARIABLE_TYPE_LONG:
				if identifier.Var.LocalValOffset == 0 {
					op = []byte{cg.OP_lstore_0}
				} else if identifier.Var.LocalValOffset == 1 {
					op = []byte{cg.OP_lstore_1}
				} else if identifier.Var.LocalValOffset == 2 {
					op = []byte{cg.OP_lstore_2}
				} else if identifier.Var.LocalValOffset == 3 {
					op = []byte{cg.OP_lstore_3}
				} else if identifier.Var.LocalValOffset <= 255 {
					op = []byte{cg.OP_lstore, identifier.Var.LocalValOffset}
				} else {
					panic("local float var out of range")
				}
			}
			target = identifier.Var.Typ
		}
	case ast.EXPRESSION_TYPE_INDEX:
		maxstack = 2
		index := e.Data.(*ast.ExpressionIndex)
		stack, _ := m.build(class, code, index.Expression, context)
		if stack > maxstack {
			maxstack = stack
		}
		stack, _ = m.build(class, code, index.Index, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
		}
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			op = []byte{cg.OP_bastore}
		case ast.VARIABLE_TYPE_SHORT:
			op = []byte{cg.OP_sastore}
		case ast.VARIABLE_TYPE_INT:
			op = []byte{cg.OP_iastore}
		case ast.VARIABLE_TYPE_FLOAT:
			op = []byte{cg.OP_fastore}
		case ast.VARIABLE_TYPE_DOUBLE:
			op = []byte{cg.OP_dastore}
		case ast.VARIABLE_TYPE_LONG:
			op = []byte{cg.OP_lastore}
		}
	case ast.EXPRESSION_TYPE_DOT:

	default:
		panic(2)
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
		return
	}
	return
}
func (m *MakeExpression) buildRelations(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 4
	bin := e.Data.(*ast.ExpressionBinary)
	maxstack1, _ := m.build(class, code, bin.Left, context)
	if maxstack1 > maxstack {
		maxstack = maxstack1
	}
	if bin.Left.VariableType.IsNumber() {
		if bin.Left.VariableType.IsInteger() {
			m.primitiveConverter(code, bin.Left.VariableType.Typ, ast.VARIABLE_TYPE_LONG)
		}
		maxstack2, _ := m.build(class, code, bin.Left, context)
		if maxstack2+2 > maxstack {
			maxstack = maxstack2 + 2
		}
		if bin.Right.VariableType.IsInteger() {
			m.primitiveConverter(code, bin.Right.VariableType.Typ, ast.VARIABLE_TYPE_LONG)
		}

		return
	}
	if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_BOOL {

		return
	}
	panic(1)
	return
}

func (m *MakeExpression) buildArithmetic(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	maxstack = 4
	maxstack1, _ := m.build(class, code, bin.Left, context)
	if maxstack1 > maxstack {
		maxstack = maxstack1
	}
	if e.Typ == ast.EXPRESSION_TYPE_OR || e.Typ == ast.EXPRESSION_TYPE_AND {
		maxstack2, _ := m.build(class, code, bin.Right, context)
		if maxstack2+2 > maxstack {
			maxstack = maxstack2 + 2
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
		//
		target := m.tt2What(bin.Left.Typ, bin.Right.Typ)
		if target > 0 {
			m.primitiveConverter(code, bin.Left.Typ, target)
		}
		maxstack2, _ := m.build(class, code, bin.Right, context)
		if maxstack2 > maxstack {
			maxstack = maxstack2
		}
		if target > 0 {
			m.primitiveConverter(code, bin.Right.Typ, target)
		}
		if target == -1 {
			target = bin.Left.Typ
		}
		switch target {
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
		maxstack2, _ := m.build(class, code, bin.Right, context)
		if maxstack2+2 > maxstack {
			maxstack = maxstack2 + 2
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
	exits = [][]byte{}
	bin := e.Data.(*ast.ExpressionBinary)
	var stack1, stack2 uint16
	var exits1, exits2 [][]byte
	stack1, exits1 = m.build(class, code, bin.Left, context)
	exits = append(exits, exits1...)
	if e.Typ == ast.EXPRESSION_TYPE_LOGICAL_OR {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifge
		exits = append(exits, code.Codes[code.CodeLength+2:])
		code.CodeLength += 4
		stack2, exits2 = m.build(class, code, bin.Right, context)
		exits = append(exits, exits2...)
	} else { //and

	}
	if stack1 > stack2 {
		maxstack = stack1
	} else {
		maxstack = stack2
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

func (m *MakeExpression) tt2What(t1, t2 int) int {
	if t1 == t2 {
		return -1
	}
	if t1 == ast.VARIABLE_TYPE_DOUBLE || t2 == ast.VARIABLE_TYPE_DOUBLE {
		return ast.VARIABLE_TYPE_DOUBLE
	}
	if t1 == ast.VARIABLE_TYPE_FLOAT || t2 == ast.VARIABLE_TYPE_FLOAT {
		return ast.VARIABLE_TYPE_FLOAT
	}
	if t1 == ast.VARIABLE_TYPE_LONG || t2 == ast.VARIABLE_TYPE_LONG {
		return ast.VARIABLE_TYPE_LONG
	}
	if t1 == ast.VARIABLE_TYPE_INT || t2 == ast.VARIABLE_TYPE_INT {
		return ast.VARIABLE_TYPE_INT
	}
	if t1 == ast.VARIABLE_TYPE_SHORT || t2 == ast.VARIABLE_TYPE_SHORT {
		return ast.VARIABLE_TYPE_SHORT
	}
	return -1
}
