package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildMapIndex(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	maxstack, _ = m.build(class, code, index.Expression, context, nil)
	currentStack := uint16(1)
	//build index
	stack, _ := m.build(class, code, index.Index, context, nil)
	if t := currentStack + stack; t > maxstack {
		maxstack = t
	}
	currentStack = 2 // mapref kref
	if index.Expression.VariableType.Map.K.IsPointer() == false {
		primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, index.Expression.VariableType.Map.K)
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Method:     "get",
		Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if index.Expression.VariableType.Map.V.IsPointer() {
		primitiveObjectConverter.castPointerTypeToRealType(class, code, index.Expression.VariableType.Map.V)
	} else {
		code.Codes[code.CodeLength] = cg.OP_dup // incrment the stack
		code.CodeLength++
		if t := 1 + currentStack; t > maxstack {
			maxstack = t
		}
		code.Codes[code.CodeLength] = cg.OP_ifnonnull
		codeLength := code.CodeLength
		code.CodeLength += 3
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
		code.Codes[code.CodeLength] = cg.OP_goto
		codeLength2 := code.CodeLength
		code.CodeLength += 3
		// no null goto here
		{
			t := &ast.VariableType{}
			t.Typ = ast.VARIABLE_TYPE_OBJECT
			t.Class = &ast.Class{}
			t.Class.Name = java_root_class
			state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, t)...)
			context.MakeStackMap(code, state, code.CodeLength)
		}

		binary.BigEndian.PutUint16(code.Codes[codeLength+1:codeLength+3], uint16(code.CodeLength-codeLength))
		primitiveObjectConverter.getFromObject(class, code, index.Expression.VariableType.Map.V)
		{
			state.popStack(1) // pop java_root_class ref
			state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, e.VariableType)...)
			context.MakeStackMap(code, state, code.CodeLength)
		}
		binary.BigEndian.PutUint16(code.Codes[codeLength2+1:codeLength2+3], uint16(code.CodeLength-codeLength2))
	}
	return
}

func (m *MakeExpression) buildIndex(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	index := e.Data.(*ast.ExpressionIndex)
	if index.Expression.VariableType.Typ == ast.VARIABLE_TYPE_MAP {
		return m.buildMapIndex(class, code, e, context, state)
	}
	maxstack, _ = m.build(class, code, index.Expression, context, nil)
	stack, _ := m.build(class, code, index.Index, context, nil)
	if t := stack + 1; t > maxstack {
		maxstack = t
	}
	meta := ArrayMetas[e.VariableType.Typ]
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Method:     "get",
		Descriptor: meta.getDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if e.VariableType.IsPointer() && e.VariableType.Typ != ast.VARIABLE_TYPE_STRING {
		primitiveObjectConverter.castPointerTypeToRealType(class, code, e.VariableType)
	}
	return
}
