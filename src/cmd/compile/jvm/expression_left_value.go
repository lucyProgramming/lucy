package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) getCaptureIdentiferLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxstack, remainStack uint16, op []byte,
	target *ast.VariableType, classname, fieldname, fieldDescriptor string) {
	identifier := e.Data.(*ast.ExpressionIdentifer)
	target = identifier.Var.Typ
	op = []byte{cg.OP_putfield}
	meta := closure.getMeta(identifier.Var.Typ.Typ)
	if context.function.ClosureVars.ClosureVariableExist(identifier.Var) { // capture var exits
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, 0)...)
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      class.Name,
			Field:      identifier.Name,
			Descriptor: "L" + meta.className + ";",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, identifier.Var.LocalValOffset)...)
	}
	state.Stacks = append(state.Stacks,
		state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(meta.className))...)
	maxstack = 1
	remainStack = 1
	classname = meta.className
	fieldname = meta.fieldName
	fieldDescriptor = meta.fieldDescriptor
	return
}

func (m *MakeExpression) getMapLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxstack, remainStack uint16, op []byte,
	target *ast.VariableType, classname, name, descriptor string) {
	index := e.Data.(*ast.ExpressionIndex)
	maxstack, _ = m.build(class, code, index.Expression, context, nil)
	stack, _ := m.build(class, code, index.Index, context, nil)
	if t := 1 + stack; t > maxstack {
		maxstack = t
	}
	state.Stacks = append(state.Stacks,
		state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(java_hashmap_class))...)
	state.Stacks = append(state.Stacks,
		state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(java_root_class))...)
	primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, index.Index.Value)
	remainStack = 2
	op = []byte{cg.OP_invokevirtual, cg.OP_pop}
	target = index.Expression.Value.Map.V
	classname = java_hashmap_class
	name = "put"
	descriptor = "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"
	return
}

func (m *MakeExpression) getLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (
	maxstack, remainStack uint16, op []byte,
	target *ast.VariableType, classname, name, descriptor string) {
	switch e.Typ {
	case ast.EXPRESSION_TYPE_IDENTIFIER:
		identifier := e.Data.(*ast.ExpressionIdentifer)
		if identifier.Var.IsGlobal {
			op = []byte{cg.OP_putstatic}
			target = identifier.Var.Typ
			classname = m.MakeClass.mainclass.Name
			name = identifier.Name
			descriptor = Descriptor.typeDescriptor(identifier.Var.Typ)
			return
		}
		if identifier.Var.BeenCaptured {
			return m.getCaptureIdentiferLeftValue(class, code, e, context, state)
		}
		if identifier.Name == ast.NO_NAME_IDENTIFIER {
			return //
		}
		switch identifier.Var.Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_ENUM:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_istore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_istore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_istore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_istore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_istore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local int var offset > 255")
			}
		case ast.VARIABLE_TYPE_FLOAT:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_fstore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_fstore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_fstore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_fstore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_fstore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		case ast.VARIABLE_TYPE_DOUBLE:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_dstore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_dstore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_dstore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_dstore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_dstore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		case ast.VARIABLE_TYPE_LONG:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_lstore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_lstore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_lstore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_lstore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_lstore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		default: // must be a object type
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_astore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_astore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_astore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_astore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_astore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		}
		target = identifier.Var.Typ
	case ast.EXPRESSION_TYPE_INDEX:
		index := e.Data.(*ast.ExpressionIndex)
		if index.Expression.Value.Typ == ast.VARIABLE_TYPE_ARRAY {
			maxstack, _ = m.build(class, code, index.Expression, context, state)
			stack, _ := m.build(class, code, index.Index, context, state)
			if t := stack + 1; t > maxstack {
				maxstack = t
			}
			meta := ArrayMetas[e.Value.Typ]
			classname = meta.classname
			name = "set"
			descriptor = meta.setDescriptor
			target = e.Value
			remainStack = 2 // [objectref ,index]
			state.Stacks = append(state.Stacks,
				state.newStackMapVerificationTypeInfo(class, index.Expression.Value)...)
			state.Stacks = append(state.Stacks,
				state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
			op = []byte{cg.OP_invokevirtual}
		} else if index.Expression.Value.Typ == ast.VARIABLE_TYPE_MAP { // map
			return m.getMapLeftValue(class, code, e, context, state)
		} else { // java array
			maxstack, _ = m.build(class, code, index.Expression, context, state)
			stack, _ := m.build(class, code, index.Index, context, state)
			if t := stack + 1; t > maxstack {
				maxstack = t
			}
			target = e.Value
			remainStack = 2 // [objectref ,index]
			state.Stacks = append(state.Stacks,
				state.newStackMapVerificationTypeInfo(class, index.Expression.Value)...)
			state.Stacks = append(state.Stacks,
				state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
			switch e.Value.Typ {
			case ast.VARIABLE_TYPE_BOOL:
				fallthrough
			case ast.VARIABLE_TYPE_BYTE:
				fallthrough
			case ast.VARIABLE_TYPE_SHORT:
				fallthrough
			case ast.VARIABLE_TYPE_ENUM:
				fallthrough
			case ast.VARIABLE_TYPE_INT:
				op = []byte{cg.OP_iastore}
			case ast.VARIABLE_TYPE_LONG:
				op = []byte{cg.OP_lastore}
			case ast.VARIABLE_TYPE_FLOAT:
				op = []byte{cg.OP_fastore}
			case ast.VARIABLE_TYPE_DOUBLE:
				op = []byte{cg.OP_dastore}
			case ast.VARIABLE_TYPE_STRING:
				fallthrough
			case ast.VARIABLE_TYPE_OBJECT:
				fallthrough
			case ast.VARIABLE_TYPE_MAP:
				fallthrough
			case ast.VARIABLE_TYPE_ARRAY:
				fallthrough
			case ast.VARIABLE_TYPE_JAVA_ARRAY:
				op = []byte{cg.OP_aastore}
			}
			return
		}
	case ast.EXPRESSION_TYPE_DOT:
		dot := e.Data.(*ast.ExpressionDot)
		if dot.Expression.Value.Typ == ast.VARIABLE_TYPE_PACKAGE {
			op = []byte{cg.OP_putstatic}
			target = dot.PackageVariableDefinition.Typ
			classname = dot.Expression.Value.Package.Name + "/main"
			name = dot.PackageVariableDefinition.Name
			descriptor = dot.PackageVariableDefinition.Descriptor
			maxstack = 0
			remainStack = 0
		} else {
			classname = dot.Expression.Value.Class.Name
			target = dot.Field.VariableDefinition.Typ
			name = dot.Name
			if dot.Field.LoadFromOutSide {
				descriptor = dot.Field.Descriptor
			} else {
				descriptor = Descriptor.typeDescriptor(target)
			}
			if dot.Field.IsStatic() {
				op = []byte{cg.OP_putstatic}
			} else {
				op = []byte{cg.OP_putfield}
				maxstack, _ = m.build(class, code, dot.Expression, context, state)
				remainStack = 1
				state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, dot.Expression.Value)...)
			}
		}
	}
	return
}
