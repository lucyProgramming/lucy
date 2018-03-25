package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildUnary(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	if e.Typ == ast.EXPRESSION_TYPE_NEGATIVE {
		var es []*cg.JumpBackPatch
		maxstack, es = m.build(class, code, e.Data.(*ast.Expression), context)
		backPatchEs(es, code.CodeLength)
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_ineg
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_fneg
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dneg
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lneg
		}
		code.CodeLength++
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_NOT {
		var es []*cg.JumpBackPatch
		maxstack, es = m.build(class, code, e.Data.(*ast.Expression), context)
		backPatchEs(es, code.CodeLength)
		code.Codes[code.CodeLength] = cg.OP_ifne
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], uint16(code.CodeLength+7))
		code.Codes[code.CodeLength+3] = cg.OP_iconst_1
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], uint16(code.CodeLength+4))
		code.Codes[code.CodeLength+7] = cg.OP_iconst_0
		code.CodeLength += 8
	}
	return
}
