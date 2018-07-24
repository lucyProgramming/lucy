package jvm

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	store local var according on type and offset
*/
func storeLocalVariableOps(variableType ast.VariableTypeKind, variableOffset uint16) []byte {
	switch variableType {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
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
	case ast.VariableTypeLong:
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
	case ast.VariableTypeFloat:
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
	case ast.VariableTypeDouble:
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
	case ast.VariableTypeJavaArray:
		fallthrough
	case ast.VariableTypeString:
		fallthrough
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeMap:
		fallthrough
	case ast.VariableTypeArray:
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

func loadLocalVariableOps(variableType ast.VariableTypeKind, variableOffset uint16) []byte {
	switch variableType {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
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
	case ast.VariableTypeLong:
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
	case ast.VariableTypeFloat:
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
	case ast.VariableTypeDouble:
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
	case ast.VariableTypeString:
		fallthrough
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeMap:
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeArray:
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
