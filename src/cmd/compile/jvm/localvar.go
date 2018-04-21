package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) loadLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		return closure.loadLocalCloureVar(class, code, v)
	}
	maxstack = jvmSize(v.Typ)
	switch v.Typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_iload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_iload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_iload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_iload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_iload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_LONG:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_lload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_lload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_lload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_lload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_lload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_fload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_fload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_fload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_fload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_fload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_dload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_dload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_dload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_dload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_dload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY, ast.VARIABLE_TYPE_JAVA_ARRAY: //[]int
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_aload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_aload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_aload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_aload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	}
	return
}

func (m *MakeClass) storeLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		closure.storeLocalCloureVar(class, code, v)
		return
	}
	maxstack = jvmSize(v.Typ)
	switch v.Typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_istore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_istore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_istore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_istore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_istore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_LONG:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_lstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_lstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_lstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_lstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_lstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_fstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_fstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_fstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_fstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_fstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_dstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_dstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_dstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_dstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_dstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY, ast.VARIABLE_TYPE_JAVA_ARRAY: //[]int
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_astore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_astore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_astore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_astore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_astore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	}
	return
}
