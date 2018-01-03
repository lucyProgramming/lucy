package jvm

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type MakeClass struct {
	p              *ast.Package
	Classes        []*cg.ClassHighLevel
	mainclass      *cg.ClassHighLevel
	MakeExpression MakeExpression
}

func (m *MakeClass) Make(p *ast.Package) {
	m.p = p
	mainclass := &cg.ClassHighLevel{}
	m.mainclass = mainclass
	mainclass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	mainclass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainclass.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	mainclass.SuperClass = ast.JAVA_ROOT_CLASS
	mainclass.Name = strings.Title(p.Name)
	mainclass.Fields = make(map[string]*cg.FiledHighLevel)
	mainclass.Methods = make(map[string][]*cg.MethodHighLevel)
	m.mkVars()
	m.mkConsts()
	m.mkEnums()
	m.mkClass()
	m.mkFuncs()
	m.mkBlocks()
}

func (m *MakeClass) mkVars() {
	for k, v := range m.p.Block.Vars {
		f := &cg.FiledHighLevel{}
		f.AccessFlags = v.AccessFlags
		f.Descriptor = v.Typ.Descriptor()
		m.mainclass.Fields[k] = f
	}
}
func (m *MakeClass) mkConsts() {

}
func (m *MakeClass) mkBlocks() {

}
func (m *MakeClass) mkEnums() {

}
func (m *MakeClass) mkClass() {

}

func (m *MakeClass) mkFuncs() {
	for _, f := range m.p.Block.Funcs {
		if f.Isbuildin {
			continue
		}
		m.mkFunc(f, "")
	}
}

func (m *MakeClass) mkFunc(f *ast.Function, path string) {
	if f.IsGlobal || f.Typ.ClosureVars == nil || len(f.Typ.ClosureVars) == 0 {
		context := &Context{}
		method := m.buildFunction(m.mainclass, f, context, true, path)
		m.mainclass.Methods[method.Name] = []*cg.MethodHighLevel{method}
		method.AccessFlags = 0
		if method.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
		}
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.ClassHighLevel = m.mainclass
		return
	}
	context := &Context{}
	class := m.mkClosureFunctionClass()
	m.buildFunction(class, f, context, false, path)
}

func (m *MakeClass) mkClosureFunctionClass() *cg.ClassHighLevel {
	ret := &cg.ClassHighLevel{}
	ret.AccessFlags = cg.ACC_CLASS_FINAL
	return ret
}
func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, f *ast.Function, context *Context, isstatic bool, path string) *cg.MethodHighLevel {
	ret := &cg.MethodHighLevel{}
	ret.Name = mkPath(path, f.Name)
	f.Method = ret
	ret.Code.Codes = make([]byte, 65536)
	ret.Code.CodeLength = 0
	m.buildBlock(class, &ret.Code, f.Block, context, ret.Name)
	ret.Descriptor = f.Descriptor
	return ret
}

func (m *MakeClass) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context, path string) {
	for _, s := range b.Statements {
		m.buildStatement(class, code, s, context, path)
	}
	return
}

func (m *MakeClass) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.Statement, context *Context, path string) {
	var maxstack uint16
	switch s.Typ {
	case ast.STATEMENT_TYPE_EXPRESSION:
		var es [][]byte
		maxstack, es = m.MakeExpression.build(class, code, s.Expression, context)
		backPatchEs(es, code)
	case ast.STATEMENT_TYPE_IF:
		maxstack = m.buildIfStatement(class, code, s.StatementIf, context, path)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		backPatchEs(s.StatementIf.BackPatchs, code)
	case ast.STATEMENT_TYPE_BLOCK:
		m.buildBlock(class, code, s.Block, context, path)
	case ast.STATEMENT_TYPE_FOR:
		maxstack = m.buildForStatement(class, code, s.StatementFor, context, path)
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
		maxstack = m.buildReturnStatement(class, code, s.StatementReturn, context, path)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
	case ast.STATEMENT_TYPE_SWITCH:
		maxstack = m.buildSwitchStatement(class, code, s.StatementSwitch, context, path)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		backPatchEs(s.StatementSwitch.BackPatchs, code)
	case ast.STATEMENT_TYPE_SKIP: // skip this block
		panic("11111111")
	}
	return
}
func (m *MakeClass) buildIfStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementIF, context *Context, path string) (maxstack uint16) {
	stack, es := m.MakeExpression.build(class, code, s.Condition, context)
	backPatchEs(es, code)
	if stack > maxstack {
		maxstack = stack
	}
	code.Codes[code.CodeLength] = cg.OP_ifeq
	falseExit := code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	m.buildBlock(class, code, s.Block, context, path)
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
		m.buildBlock(class, code, v.Block, context, path)
	}
	if s.ElseBlock != nil {
		backPatchEs([][]byte{falseExit}, code)
		falseExit = nil
		m.buildBlock(class, code, s.ElseBlock, context, path)
	}
	if falseExit != nil {
		backPatchEs([][]byte{falseExit}, code)
	}
	return
}

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context, path string) (maxstack uint16) {
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
	m.buildBlock(class, code, s.Block, context, mkPath(path, fmt.Sprintf("for%d", s.Num)))
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

func (m *MakeClass) buildSwitchStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementSwitch, context *Context, path string) (maxstack uint16) {
	return
}
func (m *MakeClass) buildReturnStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementReturn, context *Context, path string) (maxstack uint16) {
	if len(s.Function.Typ.Returns) == 0 {
		code.Codes[code.CodeLength] = cg.OP_return
		code.CodeLength++
	} else if len(s.Function.Typ.Returns) == 1 {
		if len(s.Expressions) != 1 {
			panic("this is not happening")
		}
		stack, es := m.MakeExpression.build(class, code, s.Expressions[0], context)
		if stack > maxstack {
			maxstack = stack
		}
		backPatchEs(es, code)
		switch s.Function.Typ.Returns[0].Typ.Typ {
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
		panic("still working")
	}
	return
}
func (m *MakeClass) mkFuncClassMode(f *ast.Function, path string) *cg.MethodHighLevel {
	ret := &cg.MethodHighLevel{}
	return ret
}
