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
	} else if len(s.Function.Typ.ReturnList) == 1 {
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
			fallthrough
		case ast.VARIABLE_TYPE_FUNCTION:
			panic("1111111")
		default:
			panic("......a")
		}
		code.CodeLength++
	} else {

	}
	return
}
