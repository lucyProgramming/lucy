package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildDot(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	if index.Expression.VariableType.Typ == ast.VARIABLE_TYPE_CLASS {
		maxstack = e.VariableType.JvmSlotSize()
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      index.Expression.VariableType.Class.Name,
			Name:       index.Name,
			Descriptor: m.MakeClass.Descriptor.typeDescriptor(e.VariableType),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	maxstack, _ = m.build(class, code, index.Expression, context)
	if t := e.VariableType.JvmSlotSize(); t > maxstack {
		maxstack = t
	}
	code.Codes[code.CodeLength] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      index.Expression.VariableType.Class.Name,
		Name:       index.Name,
		Descriptor: m.MakeClass.Descriptor.typeDescriptor(e.VariableType),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}

func (m *MakeExpression) buildMapIndex(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	maxstack, _ = m.build(class, code, index.Expression, context)
	currentStack := uint16(1)
	//build index
	if index.Expression.VariableType.Map.K.IsPointer() == false {
		currentStack += m.prepareStackForMapAssignWhenValueIsNotPointer(class, code, index.Expression.VariableType.Map.K)
	}
	stack, _ := m.build(class, code, index.Index, context)
	if t := currentStack + stack; t > maxstack {
		maxstack = t
	}
	if index.Expression.VariableType.Map.K.IsInteger() && index.Expression.VariableType.Map.K.Typ != index.Index.VariableType.Typ {
		m.numberTypeConverter(code, index.Index.VariableType.Typ, index.Expression.VariableType.Map.K.Typ)
	}
	if t := currentStack + index.Expression.VariableType.Map.K.JvmSlotSize(); t > maxstack {
		maxstack = t
	}
	if index.Expression.VariableType.Map.K.IsPointer() == false {
		switch index.Expression.VariableType.Map.K.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Integer",
				Name:       specail_method_init,
				Descriptor: "(I)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Long",
				Name:       specail_method_init,
				Descriptor: "(J)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Float",
				Name:       specail_method_init,
				Descriptor: "(F)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Double",
				Name:       specail_method_init,
				Descriptor: "(D)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Name:       "get",
		Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// if null
	code.Codes[code.CodeLength] = cg.OP_dup
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6)
	code.Codes[code.CodeLength+4] = cg.OP_goto
	//	code.Codes[code.CodeLength+5 : code.CodeLength+7]
	code.CodeLength += 8

	if index.Expression.VariableType.Map.V.IsPointer() == false {
		switch index.Expression.VariableType.Map.V.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength] = cg.OP_iconst_0
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength] = cg.OP_lconst_0
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength] = cg.OP_fconst_0
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_pop
			code.Codes[code.CodeLength] = cg.OP_dconst_0
		}
	}
	nullexit := (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)

	if index.Expression.VariableType.Map.V.IsPointer() == false {
		switch index.Expression.VariableType.Map.V.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst("java/lang/Integer", code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Integer",
				Name:       "intValue",
				Descriptor: "()I",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst("java/lang/Long", code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Long",
				Name:       "longValue",
				Descriptor: "()J",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst("java/lang/Float", code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Float",
				Name:       "floatValue",
				Descriptor: "()F",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst("java/lang/Double", code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Double",
				Name:       "doubleValue",
				Descriptor: "()D",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
	backPatchEs([]*cg.JumpBackPatch{nullexit}, code.CodeLength)
	return
}
func (m *MakeExpression) buildIndex(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	if index.Expression.VariableType.Typ == ast.VARIABLE_TYPE_MAP {
		return m.buildMapIndex(class, code, e, context)
	}

	maxstack, _ = m.build(class, code, index.Expression, context)
	stack, _ := m.build(class, code, index.Index, context)
	if t := stack + 1; t > maxstack {
		maxstack = t
	}
	meta := ArrayMetas[e.VariableType.Typ]
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Name:       "get",
		Descriptor: meta.getDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
