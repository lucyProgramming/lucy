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
	method.Code.Codes = make([]byte, 5)
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
	method.Code.CodeLength = uint16(len(method.Code.Codes))
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
	code.Codes[code.CodeLength] = cg.OP_goto
	b := &cg.JumpBackPatch{}
	b.CurrentCodeLength = code.CodeLength
	b.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	backPatchEs([]*cg.JumpBackPatch{b}, to)
	code.CodeLength += 3
}

/*

 */
func storeSimpleVarOp(t *ast.VariableType, offset uint16) []byte {
	switch t.Typ {
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
			return []byte{cg.OP_iastore, byte(offset)}
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

func copyOP(code *cg.AttributeCode, op ...byte) {
	for k, v := range op {
		code.Codes[code.CodeLength+uint16(k)] = v
	}
	code.CodeLength += uint16(len(op))
}
