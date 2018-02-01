package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildSelfIncrement(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	ee := e.Data.(*ast.Expression)
	if t := ee.Data.(*ast.ExpressionIdentifer); ee.Typ == ast.EXPRESSION_TYPE_IDENTIFIER && t.Var.BeenCaptured == false &&
		(t.Var.Typ.Typ == ast.VARIABLE_TYPE_BYTE || t.Var.Typ.Typ == ast.VARIABLE_TYPE_CHAR ||
			t.Var.Typ.Typ == ast.VARIABLE_TYPE_SHORT || t.Var.Typ.Typ == ast.VARIABLE_TYPE_INT) {
		// identifer and not captured
		load := func() {
			switch t.Var.LocalValOffset {
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
				code.Codes[code.CodeLength] = cg.OP_iload
				code.Codes[code.CodeLength+1] = byte(t.Var.LocalValOffset)
				code.CodeLength += 2
			}
		}
		if t.Var.LocalValOffset > 255 {
			panic("over 255")
		}
		if e.IsStatementExpression == false { // I still need it`s value
			if e.Typ == ast.EXPRESSION_TYPE_DECREMENT || e.Typ == ast.EXPRESSION_TYPE_INCREMENT {
				load() // load to stack top
				maxstack = 1
			}
		}
		if e.Typ == ast.EXPRESSION_TYPE_PRE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_INCREMENT {
			code.Codes[code.CodeLength] = cg.OP_iinc
			code.Codes[code.CodeLength+1] = byte(t.Var.LocalValOffset)
			code.Codes[code.CodeLength+2] = 1
			code.CodeLength += 3
		} else { // --
			code.Codes[code.CodeLength] = cg.OP_iinc
			code.Codes[code.CodeLength+1] = byte(t.Var.LocalValOffset)
			code.Codes[code.CodeLength+2] = 255
			code.CodeLength += 3
		}
		if e.IsStatementExpression == false { // I still need it`s value
			if e.Typ == ast.EXPRESSION_TYPE_PRE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_PRE_DECREMENT {
				load() // load to stack top
				maxstack = 1
			}
		}
		return
	}
	maxstack, remainStack, op, target, classname, fieldname, fieldDescriptor := m.getLeftValue(class, code, ee, context)

	return
}
