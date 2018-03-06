package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildSelfIncrement(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	ee := e.Data.(*ast.Expression)
	if t := ee.Data.(*ast.ExpressionIdentifer); ee.Typ == ast.EXPRESSION_TYPE_IDENTIFIER &&
		t.Var.BeenCaptured == false &&
		t.Var.Typ.Typ == ast.VARIABLE_TYPE_INT {
		// identifer  and not captured and type`s int
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
			if e.IsIncrement() {
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
			if e.IsIncrement() == false {
				load() // load to stack top
				maxstack = 1
			}
		}
		return
	}
	maxstack, remainStack, op, _, classname, fieldname, fieldDescriptor := m.getLeftValue(class, code, ee, context)
	//left value must can be used as right value,
	stack, _ := m.build(class, code, ee, context) // load it`s value
	if t := stack + remainStack; t > maxstack {
		maxstack = t
	}
	currentStack := ee.VariableType.JvmSlotSize() + remainStack
	if currentStack > maxstack {
		maxstack = currentStack
	}
	if e.IsStatementExpression == false {
		if e.Typ == ast.EXPRESSION_TYPE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_DECREMENT {
			currentStack += m.controlStack2FitAssign(code, op, classname, e.VariableType)
			if currentStack > maxstack {
				maxstack = currentStack
			}
		}
	}
	switch e.VariableType.Typ {
	case ast.VARIABLE_TYPE_BYTE:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.Codes[code.CodeLength+2] = cg.OP_i2b
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_SHORT:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.Codes[code.CodeLength+2] = cg.OP_i2s
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_INT:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		code.Codes[code.CodeLength+1] = cg.OP_i2l
		if t := currentStack + 2; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+2] = cg.OP_ladd
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_FLOAT:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		code.Codes[code.CodeLength+1] = cg.OP_i2f
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+2] = cg.OP_fadd
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_DOUBLE:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		code.Codes[code.CodeLength+1] = cg.OP_i2d
		if t := currentStack + 2; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+2] = cg.OP_dadd
		code.CodeLength += 3
	}
	if e.IsStatementExpression == false {
		if e.Typ == ast.EXPRESSION_TYPE_PRE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_PRE_DECREMENT {
			currentStack += m.controlStack2FitAssign(code, op, classname, e.VariableType)
			if currentStack > maxstack {
				maxstack = currentStack
			}
		}
	}
	//copy op
	copyOP(code, op...)
	if classname != "" {
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      classname,
			Name:       fieldname,
			Descriptor: fieldDescriptor,
		}, code.Codes[code.CodeLength:code.CodeLength+2])
		code.CodeLength += 2
	}
	return
}
