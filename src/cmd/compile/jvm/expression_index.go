package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildDot(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	maxstack = 2
	stack, _ := m.build(class, code, index.Expression, context)
	if stack > maxstack {
		maxstack = stack
	}
	switch index.Expression.VariableType.Typ {
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_CLASS:
		if index.Expression.VariableType.Typ == ast.VARIABLE_TYPE_CLASS {
			code.Codes[code.CodeLength] = cg.OP_getstatic
		} else {
			code.Codes[code.CodeLength] = cg.OP_getfield // object
		}
		class.InsertFieldRef(cg.CONSTANT_Fieldref_info_high_level{
			Class:      index.Expression.VariableType.Class.Name,
			Name:       index.Name,
			Descriptor: index.Field.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	default:
		panic(1)
	}
	return
}
func (m *MakeExpression) buildIndex(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 2
	index := e.Data.(*ast.ExpressionIndex)
	stack, _ := m.build(class, code, index.Expression, context)
	if stack > maxstack {
		maxstack = stack
	}
	stack, _ = m.build(class, code, index.Expression, context)
	if stack+2 > maxstack {
		maxstack = stack + 2
	}
	switch e.VariableType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_baload
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_saload
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_iaload
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_laload
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_faload
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_daload
	case ast.VARIABLE_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_aaload
	case ast.VARIABLE_TYPE_OBJECT:
		code.Codes[code.CodeLength] = cg.OP_aaload
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		code.Codes[code.CodeLength] = cg.OP_aaload
	}
	code.CodeLength++
	return
}
