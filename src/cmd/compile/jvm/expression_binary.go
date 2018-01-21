package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildArithmetic(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	if e.Typ == ast.EXPRESSION_TYPE_OR || e.Typ == ast.EXPRESSION_TYPE_AND {
		stack, _ := m.build(class, code, bin.Left, context)
		if stack > maxstack {
			maxstack = stack
		}
		size := bin.Left.VariableType.JvmSlotSize()
		stack, _ = m.build(class, code, bin.Right, context)
		if stack+size > maxstack {
			maxstack = stack + size
		}
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			fallthrough
		case ast.VARIABLE_TYPE_FLOAT:
			if e.Typ == ast.EXPRESSION_TYPE_AND {
				code.Codes[code.CodeLength] = cg.OP_iand
			} else {
				code.Codes[code.CodeLength] = cg.OP_ior
			}
		case ast.VARIABLE_TYPE_DOUBLE:
			fallthrough
		case ast.VARIABLE_TYPE_LONG:
			if e.Typ == ast.EXPRESSION_TYPE_AND {
				code.Codes[code.CodeLength] = cg.OP_land
			} else {
				code.Codes[code.CodeLength] = cg.OP_lor
			}
		default:
			panic("~~~~~~~~~~~~")
		}
		code.CodeLength++
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_ADD || e.Typ == ast.EXPRESSION_TYPE_SUB || e.Typ == ast.EXPRESSION_TYPE_MUL ||
		e.Typ == ast.EXPRESSION_TYPE_DIV || e.Typ == ast.EXPRESSION_TYPE_MOD {
		if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_STRING || bin.Right.VariableType.Typ == ast.VARIABLE_TYPE_STRING {
			return m.buildStrCat(class, code, bin, context)
		}
		maxstack, _ = m.build(class, code, bin.Left, context)
		//number type,no doubt
		if e.VariableType.Typ != bin.Left.Typ {
			m.numberTypeConverter(code, bin.Left.Typ, e.VariableType.Typ)
		}
		stack, _ := m.build(class, code, bin.Right, context)
		if stack+2 > maxstack { // 2 in case long or double type on stack top
			maxstack = stack + 2
		}
		if e.VariableType.Typ != bin.Right.VariableType.Typ {
			m.numberTypeConverter(code, bin.Right.Typ, e.VariableType.Typ)
		}
		switch e.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			switch e.Typ {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_iadd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_isub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_imul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_idiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_irem
			}
			code.CodeLength++
		case ast.VARIABLE_TYPE_FLOAT:
			switch e.Typ {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_fadd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_fsub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_fmul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_fdiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_frem
			}
			code.CodeLength++
		case ast.VARIABLE_TYPE_DOUBLE:
			switch e.Typ {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_dadd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_dsub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_dmul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_ddiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_drem
			}
		case ast.VARIABLE_TYPE_LONG:
			switch e.Typ {
			case ast.EXPRESSION_TYPE_ADD:
				code.Codes[code.CodeLength] = cg.OP_ladd
			case ast.EXPRESSION_TYPE_SUB:
				code.Codes[code.CodeLength] = cg.OP_lsub
			case ast.EXPRESSION_TYPE_MUL:
				code.Codes[code.CodeLength] = cg.OP_lmul
			case ast.EXPRESSION_TYPE_DIV:
				code.Codes[code.CodeLength] = cg.OP_ldiv
			case ast.EXPRESSION_TYPE_MOD:
				code.Codes[code.CodeLength] = cg.OP_lrem
			}
		default:
			panic("~~~~~~~~~~~~")
		}
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == ast.EXPRESSION_TYPE_RIGHT_SHIFT {
		stack, _ := m.build(class, code, bin.Right, context)
		if stack+2 > maxstack {
			maxstack = stack + 2
		}
		if bin.Left.VariableType.Typ == ast.VARIABLE_TYPE_LONG { // long
			if bin.Right.VariableType.Typ != ast.VARIABLE_TYPE_LONG {
				m.stackTop2Long(code, bin.Right.VariableType.Typ)
			}
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
				code.Codes[code.CodeLength] = cg.OP_lshl
			} else {
				code.Codes[code.CodeLength] = cg.OP_lshr
			}
		} else { // int
			if bin.Right.VariableType.Typ != ast.VARIABLE_TYPE_INT {
				m.stackTop2Int(code, bin.Right.VariableType.Typ)
			}
			if e.Typ == ast.EXPRESSION_TYPE_LEFT_SHIFT {
				code.Codes[code.CodeLength] = cg.OP_ishl
			} else {
				code.Codes[code.CodeLength] = cg.OP_ishr
			}
		}
		code.CodeLength++
		return
	}
	panic(12)
	return
}

func (m *MakeExpression) buildStrCat(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.ExpressionBinary, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2 // current stack is 2
	stack, _ := m.build(class, code, e.Left, context)
	maxstack += stack
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/lang/StringBuilder",
		Name:  `<init>`,
		Type:  "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	stack, _ = m.build(class, code, e.Right, context)
	maxstack += stack
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/lang/StringBuilder",
		Name:  `append`,
		Type:  "(Ljava/lang/String;)java/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/lang/StringBuilder",
		Name:  `toString`,
		Type:  "()java/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	return maxstack
}

func (m *MakeExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16, exits [][]byte) {
	maxstack = 1
	exits = [][]byte{}
	bin := e.Data.(*ast.ExpressionBinary)
	var stack uint16
	var exits1 [][]byte
	stack, exits1 = m.build(class, code, bin.Left, context)
	if stack > maxstack {
		maxstack = stack
	}
	exits = append(exits, exits1...)
	if e.Typ == ast.EXPRESSION_TYPE_LOGICAL_OR {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifne
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		stack, exits1 = m.build(class, code, bin.Left, context)
		if stack+1 > maxstack {
			maxstack = stack + 1
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifne
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
	} else { //and
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifeq
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		stack, exits1 = m.build(class, code, bin.Left, context)
		if stack+1 > maxstack {
			maxstack = stack + 1
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_ifeq
		exits = append(exits, code.Codes[code.CodeLength+2:code.CodeLength+4])
	}
	return
}
