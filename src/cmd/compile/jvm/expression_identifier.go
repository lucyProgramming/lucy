package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildCapturedIdentifer(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifer)
	cv := context.function.ClosureVars.ClosureVarsExist(identifier.Var)
	if cv == nil { // this var been captured,is declare in this function
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, identifier.Var.LocalValOffset)...)
		if 1 > maxstack {
			maxstack = 1
		}
		if t := identifier.Var.Typ.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		switch identifier.Var.Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			meta := closure.ClosureObjectMetas[CLOSURE_INT_CLASS]
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Name:       meta.fieldName,
				Descriptor: meta.fieldDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_LONG:
			meta := closure.ClosureObjectMetas[CLOSURE_LONG_CLASS]
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Name:       meta.fieldName,
				Descriptor: meta.fieldDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_FLOAT:
			meta := closure.ClosureObjectMetas[CLOSURE_FLOAT_CLASS]
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Name:       meta.fieldName,
				Descriptor: meta.fieldDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_DOUBLE:
			meta := closure.ClosureObjectMetas[CLOSURE_DOUBLE_CLASS]
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Name:       meta.fieldName,
				Descriptor: meta.fieldDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_STRING:
			meta := closure.ClosureObjectMetas[CLOSURE_STRING_CLASS]
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Name:       meta.fieldName,
				Descriptor: meta.fieldDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_MAP:
			meta := closure.ClosureObjectMetas[CLOSURE_OBJECT_CLASS]
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Name:       meta.fieldName,
				Descriptor: meta.fieldDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_OBJECT:
		case ast.VARIABLE_TYPE_ARRAY:
		case ast.VARIABLE_TYPE_JAVA_ARRAY:
		}
		return
	}

	return
}

func (m *MakeExpression) buildIdentifer(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifer)
	if identifier.Var.IsGlobal { //fetch global var
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      context.mainclass.Name,
			Name:       identifier.Name,
			Descriptor: Descriptor.typeDescriptor(identifier.Var.Typ),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxstack = identifier.Var.Typ.JvmSlotSize()
		return
	}
	if identifier.Var.BeenCaptured {
		return m.buildCapturedIdentifer(class, code, e, context)
	}

	switch identifier.Var.Typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_iload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_iload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_iload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_iload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_iload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("local int var out of range")
		}
		maxstack = 1
	case ast.VARIABLE_TYPE_FLOAT:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_fload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_fload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_fload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_fload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_fload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("local float var out of range")
		}
		maxstack = 1
	case ast.VARIABLE_TYPE_DOUBLE:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_dload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_dload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_dload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_dload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_dload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("local double var out of range")
		}
		maxstack = 2
	case ast.VARIABLE_TYPE_LONG:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_lload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_lload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_lload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_lload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_lload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("local double var out of range")
		}
		maxstack = 2
	default: // object types
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_aload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_aload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_aload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_aload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("local object var out of range")
		}
		maxstack = 1
	}
	return
}
