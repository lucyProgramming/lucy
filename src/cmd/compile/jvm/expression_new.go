package jvm

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildNew(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	if e.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY_INSTANCE {
		return m.buildNewArray(class, code, e, context)
	}
	if e.VariableType.Typ == ast.VARIABLE_TYPE_MAP {
		return m.buildNewMap(class, code, e, context)
	}
	n := e.Data.(*ast.ExpressionNew)
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(n.Typ.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2
	stackneed := maxstack
	size := uint16(0)
	for _, v := range n.Args {
		if v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || ast.EXPRESSION_TYPE_METHOD_CALL == v.Typ {
			panic(1)
		}
		size = e.VariableType.JvmSlotSize()
		stack, es := m.build(class, code, v, context)
		if stackneed+stack > maxstack {
			maxstack = stackneed + stack
		}
		stackneed += size
		backPatchEs(es, code.CodeLength)
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	methodref := cg.CONSTANT_Methodref_info_high_level{
		Class:      n.Typ.Class.Name,
		Name:       n.Construction.Func.Name,
		Descriptor: n.Construction.Func.Descriptor,
	}
	class.InsertMethodRefConst(methodref, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
func (m *MakeExpression) buildNewMap(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	maxstack = 2
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Name:       specail_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
func (m *MakeExpression) buildNewArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	//new
	n := e.Data.(*ast.ExpressionNew)
	meta := ArrayMetas[e.VariableType.ArrayType.Typ]
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@", e.VariableType.ArrayType.Typ == ast.VARIABLE_TYPE_ARRAY)
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxstack = 2
	// call init
	stack, _ := m.build(class, code, n.Args[0], context) // must be a interger
	maxstack += stack
	code.Codes[code.CodeLength] = cg.OP_dup // dup top
	code.Codes[code.CodeLength+1] = cg.OP_iconst_2
	code.Codes[code.CodeLength+2] = cg.OP_imul
	if 5 > maxstack {
		maxstack = 5
	}
	code.CodeLength += 3
	switch e.VariableType.ArrayType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_BOOLEAN
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_BYTE
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_SHORT
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_INT
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_LONG
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_FLOAT
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ATYPE_T_DOUBLE
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(java_string_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		meta := ArrayMetas[e.VariableType.ArrayType.ArrayType.Typ]
		class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	default:
		panic("sssssssssssss")
	}
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Name:       specail_method_init,
		Descriptor: meta.initFuncDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return

}
