package jvm

//import (
//	"github.com/756445638/lucy/src/cmd/compile/ast"
//	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
//)

//func (m *MakeExpression) buildConst(class *cg.ClassHighLevel, code *cg.AttributeCode, c *ast.Const) (maxstack uint16) {
//	maxstack = c.Typ.JvmSlotSize()
//	switch c.Typ.Typ {
//	case ast.VARIABLE_TYPE_BOOL:
//		if c.Data.(bool) {
//			code.Codes[code.CodeLength] = cg.OP_iconst_1
//		} else {
//			code.Codes[code.CodeLength] = cg.OP_iconst_0
//		}
//		code.CodeLength++
//	case ast.VARIABLE_TYPE_BYTE:
//		switch c.Data.(byte) {
//		case 0:
//			code.Codes[code.CodeLength] = cg.OP_iconst_0
//			code.CodeLength++
//		case 1:
//			code.Codes[code.CodeLength] = cg.OP_iconst_1
//			code.CodeLength++
//		case 2:
//			code.Codes[code.CodeLength] = cg.OP_iconst_2
//			code.CodeLength++
//		case 3:
//			code.Codes[code.CodeLength] = cg.OP_iconst_3
//			code.CodeLength++
//		case 4:
//			code.Codes[code.CodeLength] = cg.OP_iconst_4
//			code.CodeLength++
//		case 5:
//			code.Codes[code.CodeLength] = cg.OP_iconst_5
//			code.CodeLength++
//		default:
//			code.Codes[code.CodeLength] = cg.OP_bipush
//			code.Codes[code.CodeLength+1] = c.Data.(byte)
//			code.CodeLength += 2
//		}

//	case ast.VARIABLE_TYPE_SHORT:

//	case ast.VARIABLE_TYPE_INT:
//		switch c.Data.(int32) {
//		case 0:
//			code.Codes[code.CodeLength] = cg.OP_iconst_0
//			code.CodeLength++
//		case 1:
//			code.Codes[code.CodeLength] = cg.OP_iconst_1
//			code.CodeLength++
//		case 2:
//			code.Codes[code.CodeLength] = cg.OP_iconst_2
//			code.CodeLength++
//		case 3:
//			code.Codes[code.CodeLength] = cg.OP_iconst_3
//			code.CodeLength++
//		case 4:
//			code.Codes[code.CodeLength] = cg.OP_iconst_4
//			code.CodeLength++
//		case 5:
//			code.Codes[code.CodeLength] = cg.OP_iconst_5
//			code.CodeLength++
//		default:
//			code.Codes[code.CodeLength] = cg.OP_ldc_w

//		}
//	case ast.VARIABLE_TYPE_LONG:
//	case ast.VARIABLE_TYPE_FLOAT:
//	case ast.VARIABLE_TYPE_DOUBLE:
//	case ast.VARIABLE_TYPE_STRING:
//	}
//	return
//}
