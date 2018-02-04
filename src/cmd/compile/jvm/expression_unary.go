package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
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
		code.Codes[code.CodeLength] = cg.OP_ifne                                      // length 1
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], code.CodeLength+7) // length 2
		code.Codes[code.CodeLength+3] = cg.OP_iconst_1                                // length 1
		code.Codes[code.CodeLength+4] = cg.OP_goto                                    // length 1
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:], code.CodeLength+4) // length 2
		code.Codes[code.CodeLength+7] = cg.OP_iconst_0                                // length 1
		code.CodeLength += 8
	}
	return
}
