package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.Statement, context *Context) {
	var maxstack uint16
	switch s.Typ {
	case ast.STATEMENT_TYPE_EXPRESSION:
		var es [][]byte
		maxstack, es = m.MakeExpression.build(class, code, s.Expression, context)
		backPatchEs(es, code)
	case ast.STATEMENT_TYPE_IF:
		maxstack = m.buildIfStatement(class, code, s.StatementIf, context)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		backPatchEs(s.StatementIf.BackPatchs, code)
	case ast.STATEMENT_TYPE_BLOCK:
		m.buildBlock(class, code, s.Block, context)
	case ast.STATEMENT_TYPE_FOR:
		maxstack = m.buildForStatement(class, code, s.StatementFor, context)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		backPatchEs(s.StatementFor.BackPatchs, code)
	case ast.STATEMENT_TYPE_CONTINUE:
		code.Codes[code.CodeLength] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[1:3], s.StatementFor.LoopBegin)
		code.CodeLength += 3
	case ast.STATEMENT_TYPE_BREAK:
		code.Codes[code.CodeLength] = cg.OP_goto
		if s.StatementBreak.StatementFor != nil {
			appendBackPatch(&s.StatementFor.BackPatchs, code.Codes[code.CodeLength+1:code.CodeLength+3])
		} else { // switch
			appendBackPatch(&s.StatementSwitch.BackPatchs, code.Codes[code.CodeLength+1:code.CodeLength+3])
		}
		code.CodeLength += 3
	case ast.STATEMENT_TYPE_RETURN:
		maxstack = m.buildReturnStatement(class, code, s.StatementReturn, context)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
	case ast.STATEMENT_TYPE_SWITCH:
		maxstack = m.buildSwitchStatement(class, code, s.StatementSwitch, context)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		backPatchEs(s.StatementSwitch.BackPatchs, code)
	case ast.STATEMENT_TYPE_SKIP: // skip this block
		panic("11111111")
	}
	return
}
func (m *MakeClass) buildIfStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementIF, context *Context) (maxstack uint16) {
	stack, es := m.MakeExpression.build(class, code, s.Condition, context)
	backPatchEs(es, code)
	if stack > maxstack {
		maxstack = stack
	}
	code.Codes[code.CodeLength] = cg.OP_ifeq
	falseExit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	m.buildBlock(class, code, s.Block, context)
	for _, v := range s.ElseIfList {
		backPatchEs([][]byte{falseExit}, code)
		stack, es := m.MakeExpression.build(class, code, v.Condition, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		falseExit = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		code.CodeLength += 3
		m.buildBlock(class, code, v.Block, context)
	}
	if s.ElseBlock != nil {
		backPatchEs([][]byte{falseExit}, code)
		falseExit = nil
		m.buildBlock(class, code, s.ElseBlock, context)
	}
	if falseExit != nil {
		backPatchEs([][]byte{falseExit}, code)
	}
	return
}

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	//init
	if s.Init != nil {
		stack, es := m.MakeExpression.build(class, code, s.Init, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
	}
	s.LoopBegin = code.CodeLength
	//condition
	if s.Condition != nil {
		stack, es := m.MakeExpression.build(class, code, s.Condition, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		appendBackPatch(&s.BackPatchs, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {

	}
	m.buildBlock(class, code, s.Block, context)
	if s.Post != nil {
		stack, es := m.MakeExpression.build(class, code, s.Init, context)
		backPatchEs(es, code)
		if stack > maxstack {
			maxstack = stack
		}
	}
	code.Codes[code.CodeLength] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:], s.LoopBegin)
	code.CodeLength += 3
	return
}

func (m *MakeClass) buildSwitchStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementSwitch, context *Context) (maxstack uint16) {
	return
}

func (m *MakeClass) buildReturnStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementReturn, context *Context) (maxstack uint16) {
	if len(s.Function.Typ.ReturnList) == 0 {
		code.Codes[code.CodeLength] = cg.OP_return
		code.CodeLength++
		return
	}
	if len(s.Function.Typ.ReturnList) == 1 {
		if len(s.Expressions) != 1 {
			panic("this is not happening")
		}
		var es [][]byte
		maxstack, es = m.MakeExpression.build(class, code, s.Expressions[0], context)
		backPatchEs(es, code)
		switch s.Function.Typ.ReturnList[0].Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_CHAR:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_ireturn
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lreturn
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_freturn
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dreturn
		case ast.VARIABLE_TYPE_STRING:
			fallthrough
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
			code.Codes[code.CodeLength] = cg.OP_areturn
		default:
			panic("......a")
		}
		code.CodeLength++
		return
	}
	//multi value to return
	//load array list
	m.MakeExpression.buildLoadArrayListAutoVar(class, code, context)
	// call clear
	code.Codes[code.CodeLength] = cg.OP_dup
	code.Codes[code.CodeLength+1] = cg.OP_invokevirtual
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/util/ArrayList",
		Name:       "clear",
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+2:code.CodeLength+4])
	code.CodeLength += 4
	maxstack = 2
	currentStack := uint16(1)
	for _, v := range s.Expressions {
		if (v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL) && len(v.VariableTypes) > 0 {
			code.Codes[code.CodeLength] = cg.OP_dup // dup array list
			currentStack++
			if currentStack > maxstack {
				maxstack = maxstack
			} // make the call
			stack, es := m.MakeExpression.build(class, code, v, context)
			backPatchEs(es, code)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/util/ArrayList",
				Name:       "addAll",
				Descriptor: "Ljava/util/Collection;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		code.Codes[code.CodeLength] = cg.OP_dup // dup array list
		currentStack++
		if currentStack > maxstack {
			maxstack = maxstack
		}
		stack, es := m.MakeExpression.build(class, code, v, context)
		backPatchEs(es, code)
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		//convert to object
		switch v.VariableType.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Integer",
				Name:       "valueOf",
				Descriptor: "(I)Ljava/lang/Integer;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Float",
				Name:       "valueOf",
				Descriptor: "(F)Ljava/lang/Float;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Double",
				Name:       "valueOf",
				Descriptor: "(D)Ljava/lang/Double;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Long",
				Name:       "valueOf",
				Descriptor: "(J)Ljava/lang/Long;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_OBJECT:
		case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		default:
			panic("~~~~~~~~~~~~")
		}
		// append
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/util/ArrayList",
			Name:       "add",
			Descriptor: "java/lang/Object",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		currentStack = 1
	}
	code.Codes[code.CodeLength] = cg.OP_areturn
	code.CodeLength++
	return
}
