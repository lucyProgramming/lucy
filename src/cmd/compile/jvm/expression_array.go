package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildArray(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	arr := e.Data.(*ast.ExpressionArrayLiteral)
	//	new array ,
	meta := ArrayMetas[e.VariableType.ArrayType.Typ]
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	loadInt32(class, code, int32(arr.Length))
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
	case ast.VARIABLE_TYPE_MAP:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(java_string_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(e.VariableType.ArrayType.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[e.VariableType.ArrayType.ArrayType.Typ]
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	maxstack = 4
	var index int32 = 0
	store := func() {
		switch e.VariableType.ArrayType.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			code.Codes[code.CodeLength] = cg.OP_bastore
		case ast.VARIABLE_TYPE_SHORT:
			code.Codes[code.CodeLength] = cg.OP_sastore
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_iastore
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lastore
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_fastore
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dastore
		case ast.VARIABLE_TYPE_MAP:
			fallthrough
		case ast.VARIABLE_TYPE_STRING:
			fallthrough
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY:
			code.Codes[code.CodeLength] = cg.OP_aastore
		}
		code.CodeLength++
	}
	for _, v := range arr.Expressions {
		if v.IsCall() && len(v.VariableTypes) > 1 {
			// stack top is array list
			stack, es := m.build(class, code, v, context)
			backPatchEs(es, code.CodeLength)
			if t := 3 + stack; t > maxstack {
				maxstack = t
			}
			m.buildStoreArrayListAutoVar(code, context)
			for k, t := range v.VariableTypes {
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				loadInt32(class, code, index) // load index
				stack := m.unPackArraylist(class, code, k, t, context)
				if t := 5 + stack; t > maxstack {
					maxstack = t
				}
				store()
				index++
			}
			continue
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, index) // load index
		stack, es := m.build(class, code, v, context)
		backPatchEs(es, code.CodeLength)
		if t := 5 + stack; t > maxstack {
			maxstack = t
		}
		store()
		index++
	}
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Method:     special_method_init,
		Descriptor: meta.constructorFuncDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
