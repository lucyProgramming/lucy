package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildIdentifer(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifer)
	if identifier.Var.BeenCaptured > 0 {
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
		if identifier.Var.BeenCaptured > 0 {
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
					op = []byte{cg.OP_fstore, byte(identifier.Var.LocalValOffset)}
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
					op = []byte{cg.OP_dstore, byte(identifier.Var.LocalValOffset)}
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
					op = []byte{cg.OP_lstore, byte(identifier.Var.LocalValOffset)}
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
