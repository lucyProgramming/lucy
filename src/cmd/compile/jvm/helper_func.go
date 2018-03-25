package jvm

import (
	"encoding/binary"
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	make a default construction
*/
func mkClassDefaultContruction(class *cg.ClassHighLevel) {
	method := &cg.MethodHighLevel{}
	method.Name = special_method_init
	method.Descriptor = "()V"
	method.AccessFlags |= cg.ACC_METHOD_PUBLIC
	length := 5
	method.Code.Codes = make([]byte, length)
	method.Code.Codes[0] = cg.OP_aload_0
	method.Code.Codes[1] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      class.SuperClass,
		Name:       special_method_init,
		Descriptor: "()V",
	}, method.Code.Codes[2:4])
	method.Code.Codes[4] = cg.OP_return
	method.Code.MaxStack = 1
	method.Code.MaxLocals = 1
	method.Code.CodeLength = length
	class.AppendMethod(method)
}

func backPatchEs(es []*cg.JumpBackPatch, t int) {
	for _, e := range es {
		offset := int16(t - int(e.CurrentCodeLength))
		e.Bs[0] = byte(offset >> 8)
		e.Bs[1] = byte(offset)
	}
}

func jumpto(op byte, code *cg.AttributeCode, to int) {
	b := (&cg.JumpBackPatch{}).FromCode(op, code)
	backPatchEs([]*cg.JumpBackPatch{b}, to)
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
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		fallthrough
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY:
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
		panic(fmt.Sprintf("typ:%v", t))
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
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY:
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
		code.Codes[code.CodeLength+k] = v
	}
	code.CodeLength += len(op)
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
func postion(class *cg.ClassHighLevel, code *cg.AttributeCode, s string) {
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst(s, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_pop
	code.CodeLength += 4

}
func loadInt32(class *cg.ClassHighLevel, code *cg.AttributeCode, value int32) {
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

func checkStackTopIfNagetiveThrowIndexOutOfRangeException(class *cg.ClassHighLevel, code *cg.AttributeCode) (maxstack uint16) {
	maxstack = 2
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_iflt
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
	code.Codes[code.CodeLength+3] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 12)
	code.Codes[code.CodeLength+6] = cg.OP_pop
	code.Codes[code.CodeLength+7] = cg.OP_new
	class.InsertClassConst(java_index_out_of_range_exception_class, code.Codes[code.CodeLength+8:code.CodeLength+10])
	code.Codes[code.CodeLength+10] = cg.OP_dup
	code.Codes[code.CodeLength+11] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_index_out_of_range_exception_class,
		Name:       special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+12:code.CodeLength+14])
	code.Codes[code.CodeLength+14] = cg.OP_athrow
	code.CodeLength += 15
	return
}

//func IfStackTopStringIsNullThenLoad(class *cg.ClassHighLevel, code *cg.AttributeCode, s string) (incrment uint16) {
//	incrment = 1
//	code.Codes[code.CodeLength] = cg.OP_dup
//	code.CodeLength++
//	code.Codes[code.CodeLength] = cg.OP_ifnull
//	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
//	code.Codes[code.CodeLength+3] = cg.OP_goto
//	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 7)
//	code.Codes[code.CodeLength+6] = cg.OP_pop
//	code.Codes[code.CodeLength+7] = cg.OP_ldc_w
//	class.InsertStringConst("", code.Codes[code.CodeLength+8:code.CodeLength+10])
//	code.CodeLength += 10
//	return
//}
