package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildRelations(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	if bin.Left.VariableType.IsNumber() { // in this case ,right must be a number type
		maxstack = 4
		stack, _ := m.build(class, code, bin.Left, context)
		if stack > maxstack {
			maxstack = stack
		}
		target := e.VariableType.Typ // this is target
		if target == ast.VARIABLE_TYPE_INT || target == ast.VARIABLE_TYPE_SHORT || target == ast.VARIABLE_TYPE_BYTE {
			target = ast.VARIABLE_TYPE_LONG
		}
		if target != bin.Left.VariableType.Typ {
			m.numberTypeConverter(code, bin.Left.VariableType.Typ, target)
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if t := 2 + stack; t > maxstack {
			maxstack = t
		}
		if target != bin.Right.VariableType.Typ {
			m.numberTypeConverter(code, bin.Right.VariableType.Typ, target)
		}
		switch target {
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lcmp
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_fcmpl
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dcmpl
		}
		code.CodeLength++
		if e.Typ == ast.EXPRESSION_TYPE_GT || e.Typ == ast.EXPRESSION_TYPE_LE { // > and <=
			code.Codes[code.CodeLength] = cg.OP_ifgt
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], code.CodeLength+7)
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], code.CodeLength+8)
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		}
		if e.Typ == ast.EXPRESSION_TYPE_LT || e.Typ == ast.EXPRESSION_TYPE_GE { // < and >=
			code.Codes[code.CodeLength] = cg.OP_iflt
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], code.CodeLength+7)
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], code.CodeLength+8)
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		}
		if e.Typ == ast.EXPRESSION_TYPE_EQ || e.Typ == ast.EXPRESSION_TYPE_NE { // == and !=
			code.Codes[code.CodeLength] = cg.OP_ifeq
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], code.CodeLength+7)
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], code.CodeLength+8)
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		}
		return
	}
	if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_BOOL ||
		bin.Right.VariableType.Typ == ast.VARIABLE_TYPE_BOOL { // bool type
		var es [][]byte
		maxstack, es = m.build(class, code, bin.Left, context)
		backPatchEs(es, code)
		stack, es := m.build(class, code, bin.Right, context)
		backPatchEs(es, code)
		if 1+stack > maxstack {
			maxstack = stack + 1
		}
		code.Codes[code.CodeLength] = cg.OP_if_icmpeq
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], code.CodeLength+7)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		} else {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_1
		}
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], code.CodeLength+8)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		}
		code.CodeLength += 8
		return
	}
	if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_NULL || bin.Right.VariableType.Typ == ast.VARIABLE_TYPE_NULL {
		var notNullExpression *ast.Expression
		if bin.Left.VariableType.Typ != ast.VARIABLE_TYPE_NULL {
			notNullExpression = bin.Left
		} else {
			notNullExpression = bin.Right
		}
		maxstack, _ = m.build(class, code, notNullExpression, context)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_ifnull
		} else { // ne
			code.Codes[code.CodeLength] = cg.OP_ifnonnull
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], code.CodeLength+7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], code.CodeLength+8)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}
	if bin.Left.VariableType.IsPointer() && bin.Right.VariableType.IsPointer() { //
		stack, _ := m.build(class, code, bin.Left, context)
		if stack > maxstack {
			maxstack = stack
		}
		stack, _ = m.build(class, code, bin.Right, context)
		if stack+1 > maxstack {
			maxstack = stack + 1
		}
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_if_acmpeq
		} else { // ne
			code.Codes[code.CodeLength] = cg.OP_if_acmpne
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], code.CodeLength+7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], code.CodeLength+8)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}
	panic("missing")
	return
}
