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
	stack, _ := m.build(class, code, index.Index, context)
	if t := currentStack + stack; t > maxstack {
		maxstack = t
	}
	if index.Expression.VariableType.Map.K.IsPointer() == false {
		PrimitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, index.Expression.VariableType.Map.K)
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
	if index.Expression.VariableType.Map.V.IsPointer() == false {
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 8)
	} else if index.Expression.VariableType.Map.V.Typ == ast.VARIABLE_TYPE_STRING {
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 10)
	} else {
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 6)
	}
	code.CodeLength += 7

	switch index.Expression.VariableType.Map.V.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_pop
		code.Codes[code.CodeLength+1] = cg.OP_iconst_0
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_pop
		code.Codes[code.CodeLength+1] = cg.OP_lconst_0
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_pop
		code.Codes[code.CodeLength+1] = cg.OP_fconst_0
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_pop
		code.Codes[code.CodeLength+1] = cg.OP_dconst_0
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_pop
		code.Codes[code.CodeLength+1] = cg.OP_ldc_w
		class.InsertStringConst("", code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
	}
	nullexit := (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)
	if index.Expression.VariableType.Map.V.IsPointer() == false {
		PrimitiveObjectConverter.getFromObject(class, code, index.Expression.VariableType.Map.V)
	} else {
		PrimitiveObjectConverter.castPointerTypeToRealType(class, code, index.Expression.VariableType.Map.V)
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
