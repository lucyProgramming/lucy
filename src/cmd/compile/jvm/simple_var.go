package jvm

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	store local var according on type and offset
*/
func storeLocalVariableOps(variableType int, variableOffset uint16) []byte {
	switch variableType {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_istore_0}
		case 1:
			return []byte{cg.OP_istore_1}
		case 2:
			return []byte{cg.OP_istore_2}
		case 3:
			return []byte{cg.OP_istore_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_istore, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_LONG:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_lstore_0}
		case 1:
			return []byte{cg.OP_lstore_1}
		case 2:
			return []byte{cg.OP_lstore_2}
		case 3:
			return []byte{cg.OP_lstore_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_lstore, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_fstore_0}
		case 1:
			return []byte{cg.OP_fstore_1}
		case 2:
			return []byte{cg.OP_fstore_2}
		case 3:
			return []byte{cg.OP_fstore_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_fstore, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_dstore_0}
		case 1:
			return []byte{cg.OP_dstore_1}
		case 2:
			return []byte{cg.OP_dstore_2}
		case 3:
			return []byte{cg.OP_dstore_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_dstore, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		fallthrough
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_FUNCTION:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_astore_0}
		case 1:
			return []byte{cg.OP_astore_1}
		case 2:
			return []byte{cg.OP_astore_2}
		case 3:
			return []byte{cg.OP_astore_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_astore, byte(variableOffset)}
		}
	default:
		panic(fmt.Sprintf("typ:%v", variableType))
	}
}

func loadLocalVariableOps(variableType int, variableOffset uint16) []byte {
	switch variableType {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_iload_0}
		case 1:
			return []byte{cg.OP_iload_1}
		case 2:
			return []byte{cg.OP_iload_2}
		case 3:
			return []byte{cg.OP_iload_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_iload, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_LONG:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_lload_0}
		case 1:
			return []byte{cg.OP_lload_1}
		case 2:
			return []byte{cg.OP_lload_2}
		case 3:
			return []byte{cg.OP_lload_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_lload, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_fload_0}
		case 1:
			return []byte{cg.OP_fload_1}
		case 2:
			return []byte{cg.OP_fload_2}
		case 3:
			return []byte{cg.OP_fload_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_fload, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_dload_0}
		case 1:
			return []byte{cg.OP_dload_1}
		case 2:
			return []byte{cg.OP_dload_2}
		case 3:
			return []byte{cg.OP_dload_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_dload, byte(variableOffset)}
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_FUNCTION:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY:
		switch variableOffset {
		case 0:
			return []byte{cg.OP_aload_0}
		case 1:
			return []byte{cg.OP_aload_1}
		case 2:
			return []byte{cg.OP_aload_2}
		case 3:
			return []byte{cg.OP_aload_3}
		default:
			if variableOffset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_aload, byte(variableOffset)}
		}
	default:
		panic(fmt.Sprintf("typ:%d", variableType))
	}
}
