package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

/*
	make a default construction
*/
func mkClassDefaultContruction(class *cg.ClassHighLevel) {
	method := &cg.MethodHighLevel{}
	method.Name = specail_method_init
	method.Descriptor = "()V"
	method.AccessFlags |= cg.ACC_METHOD_PRIVATE
	length := 5
	method.Code.Codes = make([]byte, length)
	method.Code.Codes[0] = cg.OP_aload_0
	method.Code.Codes[1] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      class.SuperClass,
		Name:       specail_method_init,
		Descriptor: "()V",
	}, method.Code.Codes[2:4])
	method.Code.Codes[4] = cg.OP_return
	method.Code.MaxStack = 1
	method.Code.MaxLocals = 1
	method.Code.CodeLength = uint16(length)
	class.AppendMethod(method)
}

func backPatchEs(es []*cg.JumpBackPatch, to uint16) {
	for _, e := range es {
		offset := int16(int(to) - int(e.CurrentCodeLength))
		e.Bs[0] = byte(offset >> 8)
		e.Bs[1] = byte(offset)
	}
}

func jumpto(op byte, code *cg.AttributeCode, to uint16) {
	code.Codes[code.CodeLength] = op
	b := &cg.JumpBackPatch{}
	b.CurrentCodeLength = code.CodeLength
	b.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	backPatchEs([]*cg.JumpBackPatch{b}, to)
	code.CodeLength += 3
}

/*
	store local var according on type and offset
*/
func storeSimpleVarOp(t int, offset uint16) []byte {
	switch t {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_istore, byte(offset)}
		}
	case ast.VARIABLE_TYPE_LONG:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_lstore, byte(offset)}
		}
	case ast.VARIABLE_TYPE_FLOAT:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_fstore, byte(offset)}
		}
	case ast.VARIABLE_TYPE_DOUBLE:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_dstore, byte(offset)}
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_astore, byte(offset)}
		}
	default:
		panic("...")
	}
}

func loadSimpleVarOp(t int, offset uint16) []byte {
	switch t {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_iload, byte(offset)}
		}
	case ast.VARIABLE_TYPE_LONG:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_lload, byte(offset)}
		}
	case ast.VARIABLE_TYPE_FLOAT:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_fload, byte(offset)}
		}
	case ast.VARIABLE_TYPE_DOUBLE:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_dload, byte(offset)}
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
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
			if offset > 255 {
				panic("over 255")
			}
			return []byte{cg.OP_aload, byte(offset)}
		}
	default:
		panic("...")
	}
}

func copyOP(code *cg.AttributeCode, op ...byte) {
	for k, v := range op {
		code.Codes[code.CodeLength+uint16(k)] = v
	}
	code.CodeLength += uint16(len(op))
}

func copyOPLeftValue(class *cg.ClassHighLevel, code *cg.AttributeCode, ops []byte, classname, name, descriptor string) {
	if len(ops) == 0 {
		return
	}
	code.Codes[code.CodeLength] = ops[0]
	code.CodeLength++
	if classname != "" || name != "" || descriptor != "" {
		if classname == "" || name == "" || descriptor == "" {
			panic("....")
		}
		if ops[0] == cg.OP_putstatic || ops[0] == cg.OP_putfield {
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      classname,
				Name:       name,
				Descriptor: descriptor,
			}, code.Codes[code.CodeLength:code.CodeLength+2])
		} else if ops[0] == cg.OP_invokevirtual {
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      classname,
				Name:       name,
				Descriptor: descriptor,
			}, code.Codes[code.CodeLength:code.CodeLength+2])
		} else {
			panic("...")
		}
		code.CodeLength += 2
	}
	copyOP(code, ops[1:]...)

}
func loadInt32(code *cg.AttributeCode, class *cg.ClassHighLevel, value int32) {
	switch value {
	case -1:
		code.Codes[code.CodeLength] = cg.OP_iconst_m1
		code.CodeLength++
	case 0:
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
	case 1:
		code.Codes[code.CodeLength] = cg.OP_iconst_1
		code.CodeLength++
	case 2:
		code.Codes[code.CodeLength] = cg.OP_iconst_2
		code.CodeLength++
	case 3:
		code.Codes[code.CodeLength] = cg.OP_iconst_3
		code.CodeLength++
	case 4:
		code.Codes[code.CodeLength] = cg.OP_iconst_4
		code.CodeLength++
	case 5:
		code.Codes[code.CodeLength] = cg.OP_iconst_5
		code.CodeLength++
	default:
		if -127 >= value && value <= 128 {
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
	}

}
