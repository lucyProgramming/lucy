package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildNew(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	if e.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY_INSTANCE {
		return m.buildNewArray(class, code, e, context)
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

func (m *MakeExpression) buildNewArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	//new
	n := e.Data.(*ast.ExpressionNew)
	var classname string
	code.Codes[code.CodeLength] = cg.OP_new
	switch e.VariableType.CombinationType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		classname = "lucy/lang/Arrayboolean"
		class.InsertClassConst(classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	case ast.VARIABLE_TYPE_BYTE:
		classname = "lucy/lang/Arraybyte"
		class.InsertClassConst(classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	case ast.VARIABLE_TYPE_SHORT:
		classname = "lucy/lang/Arrayshort"
		class.InsertClassConst(classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	case ast.VARIABLE_TYPE_INT:
		classname = "lucy/lang/Arrayint"
		class.InsertClassConst(classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	case ast.VARIABLE_TYPE_LONG:
		classname = "lucy/lang/Arraylong"
		class.InsertClassConst(classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	case ast.VARIABLE_TYPE_FLOAT:
		class.InsertClassConst("lucy/lang/Arrayfloat", code.Codes[code.CodeLength+1:code.CodeLength+3])
	case ast.VARIABLE_TYPE_DOUBLE:
		classname = "lucy/lang/Arraydouble"
		class.InsertClassConst(classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		classname = "lucy/lang/ArrayObject"
		class.InsertClassConst(classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	default:
		panic(1)
	}
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
	switch e.VariableType.CombinationType.Typ {
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
		class.InsertClassConst("lang/java/String", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
		panic(1)
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		panic(1)
	default:
		panic(1)
	}
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	switch e.VariableType.CombinationType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([ZI)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([BI)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([SI)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([II)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([JI)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([FI)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([DI)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      classname,
			Name:       specail_method_init,
			Descriptor: "([Ljava/lang/ObjectI)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
		panic(1)
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		panic(1)
	default:
		panic(1)
	}

	return

}
