package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	store local var according on type and offset
*/
func storeLocalVariableOps(variableType ast.VariableTypeKind, offset uint16) []byte {
	if offset > 255 { // early check
		panic("over 255")
	}
	switch variableType {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		switch offset {
		case 0:
			return []byte{cg.OP_istore_0}
		case 1:
			return []byte{cg.OP_istore_1}
		case 2:
			return []byte{cg.OP_istore_2}
		case 3:
			return []byte{cg.OP_istore_3}
		default:
			return []byte{cg.OP_istore, byte(offset)}
		}
	case ast.VariableTypeLong:
		switch offset {
		case 0:
			return []byte{cg.OP_lstore_0}
		case 1:
			return []byte{cg.OP_lstore_1}
		case 2:
			return []byte{cg.OP_lstore_2}
		case 3:
			return []byte{cg.OP_lstore_3}
		default:
			return []byte{cg.OP_lstore, byte(offset)}
		}
	case ast.VariableTypeFloat:
		switch offset {
		case 0:
			return []byte{cg.OP_fstore_0}
		case 1:
			return []byte{cg.OP_fstore_1}
		case 2:
			return []byte{cg.OP_fstore_2}
		case 3:
			return []byte{cg.OP_fstore_3}
		default:
			return []byte{cg.OP_fstore, byte(offset)}
		}
	case ast.VariableTypeDouble:
		switch offset {
		case 0:
			return []byte{cg.OP_dstore_0}
		case 1:
			return []byte{cg.OP_dstore_1}
		case 2:
			return []byte{cg.OP_dstore_2}
		case 3:
			return []byte{cg.OP_dstore_3}
		default:
			return []byte{cg.OP_dstore, byte(offset)}
		}
	default:
		switch offset {
		case 0:
			return []byte{cg.OP_astore_0}
		case 1:
			return []byte{cg.OP_astore_1}
		case 2:
			return []byte{cg.OP_astore_2}
		case 3:
			return []byte{cg.OP_astore_3}
		default:
			return []byte{cg.OP_astore, byte(offset)}
		}
	}
}

func loadLocalVariableOps(variableType ast.VariableTypeKind, offset uint16) []byte {
	if offset > 255 { // early check
		panic("over 255")
	}

	switch variableType {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		switch offset {
		case 0:
			return []byte{cg.OP_iload_0}
		case 1:
			return []byte{cg.OP_iload_1}
		case 2:
			return []byte{cg.OP_iload_2}
		case 3:
			return []byte{cg.OP_iload_3}
		default:
			return []byte{cg.OP_iload, byte(offset)}
		}
	case ast.VariableTypeLong:
		switch offset {
		case 0:
			return []byte{cg.OP_lload_0}
		case 1:
			return []byte{cg.OP_lload_1}
		case 2:
			return []byte{cg.OP_lload_2}
		case 3:
			return []byte{cg.OP_lload_3}
		default:
			return []byte{cg.OP_lload, byte(offset)}
		}
	case ast.VariableTypeFloat:
		switch offset {
		case 0:
			return []byte{cg.OP_fload_0}
		case 1:
			return []byte{cg.OP_fload_1}
		case 2:
			return []byte{cg.OP_fload_2}
		case 3:
			return []byte{cg.OP_fload_3}
		default:
			return []byte{cg.OP_fload, byte(offset)}
		}
	case ast.VariableTypeDouble:
		switch offset {
		case 0:
			return []byte{cg.OP_dload_0}
		case 1:
			return []byte{cg.OP_dload_1}
		case 2:
			return []byte{cg.OP_dload_2}
		case 3:
			return []byte{cg.OP_dload_3}
		default:
			return []byte{cg.OP_dload, byte(offset)}
		}
	default:
		switch offset {
		case 0:
			return []byte{cg.OP_aload_0}
		case 1:
			return []byte{cg.OP_aload_1}
		case 2:
			return []byte{cg.OP_aload_2}
		case 3:
			return []byte{cg.OP_aload_3}
		default:
			return []byte{cg.OP_aload, byte(offset)}
		}
	}
}
