package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) getCaptureIdentifierLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxStack, remainStack uint16, op []byte,
	target *ast.Type, className, fieldName, fieldDescriptor string) {
	identifier := e.Data.(*ast.ExpressionIdentifier)
	target = identifier.Variable.Type
	op = []byte{cg.OP_putfield}
	meta := closure.getMeta(identifier.Variable.Type.Type)
	if context.function.Closure.ClosureVariableExist(identifier.Variable) { // capture var exits
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      class.Name,
			Field:      identifier.Name,
			Descriptor: "L" + meta.className + ";",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, identifier.Variable.LocalValOffset)...)
	}
	state.pushStack(class, state.newObjectVariableType(meta.className))
	maxStack = 1
	remainStack = 1
	className = meta.className
	fieldName = meta.fieldName
	fieldDescriptor = meta.fieldDescriptor
	return
}

func (buildExpression *BuildExpression) getMapLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxStack, remainStack uint16, op []byte,
	target *ast.Type, className, name, descriptor string) {
	index := e.Data.(*ast.ExpressionIndex)
	maxStack, _ = buildExpression.build(class, code, index.Expression, context, state)
	state.pushStack(class, state.newObjectVariableType(javaMapClass))
	stack, _ := buildExpression.build(class, code, index.Index, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	if index.Index.ExpressionValue.IsPointer() == false {
		typeConverter.packPrimitives(class, code, index.Index.ExpressionValue)
	}
	state.pushStack(class, state.newObjectVariableType(javaRootClass))
	remainStack = 2
	op = []byte{}
	if index.Expression.ExpressionValue.Map.Value.IsPointer() == false {
		op = append(op,
			typeConverter.packPrimitivesBytes(class, index.Expression.ExpressionValue.Map.Value)...)
	}
	bs4 := make([]byte, 4)
	bs4[0] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaMapClass,
		Method:     "put",
		Descriptor: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
	}, bs4[1:3])
	bs4[3] = cg.OP_pop
	op = append(op, bs4...)
	target = index.Expression.ExpressionValue.Map.Value
	className = javaMapClass
	return
}

func (buildExpression *BuildExpression) getLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (
	maxStack, remainStack uint16, op []byte,
	target *ast.Type, className, name, descriptor string) {
	switch e.Type {
	case ast.ExpressionTypeIdentifier:
		identifier := e.Data.(*ast.ExpressionIdentifier)
		if identifier.Variable.IsGlobal {
			op = []byte{cg.OP_putstatic}
			target = identifier.Variable.Type
			className = buildExpression.BuildPackage.mainClass.Name
			name = identifier.Name
			descriptor = JvmDescriptor.typeDescriptor(identifier.Variable.Type)
			return
		}
		if identifier.Variable.BeenCaptured {
			return buildExpression.getCaptureIdentifierLeftValue(class, code, e, context, state)
		}
		if identifier.Name == ast.NoNameIdentifier {
			panic("this is not happening")
		}
		switch identifier.Variable.Type.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			if identifier.Variable.LocalValOffset == 0 {
				op = []byte{cg.OP_istore_0}
			} else if identifier.Variable.LocalValOffset == 1 {
				op = []byte{cg.OP_istore_1}
			} else if identifier.Variable.LocalValOffset == 2 {
				op = []byte{cg.OP_istore_2}
			} else if identifier.Variable.LocalValOffset == 3 {
				op = []byte{cg.OP_istore_3}
			} else if identifier.Variable.LocalValOffset <= 255 {
				op = []byte{cg.OP_istore, byte(identifier.Variable.LocalValOffset)}
			} else {
				panic("over 255")
			}
		case ast.VariableTypeFloat:
			if identifier.Variable.LocalValOffset == 0 {
				op = []byte{cg.OP_fstore_0}
			} else if identifier.Variable.LocalValOffset == 1 {
				op = []byte{cg.OP_fstore_1}
			} else if identifier.Variable.LocalValOffset == 2 {
				op = []byte{cg.OP_fstore_2}
			} else if identifier.Variable.LocalValOffset == 3 {
				op = []byte{cg.OP_fstore_3}
			} else if identifier.Variable.LocalValOffset <= 255 {
				op = []byte{cg.OP_fstore, byte(identifier.Variable.LocalValOffset)}
			} else {
				panic("over 255")
			}
		case ast.VariableTypeDouble:
			if identifier.Variable.LocalValOffset == 0 {
				op = []byte{cg.OP_dstore_0}
			} else if identifier.Variable.LocalValOffset == 1 {
				op = []byte{cg.OP_dstore_1}
			} else if identifier.Variable.LocalValOffset == 2 {
				op = []byte{cg.OP_dstore_2}
			} else if identifier.Variable.LocalValOffset == 3 {
				op = []byte{cg.OP_dstore_3}
			} else if identifier.Variable.LocalValOffset <= 255 {
				op = []byte{cg.OP_dstore, byte(identifier.Variable.LocalValOffset)}
			} else {
				panic("over 255")
			}
		case ast.VariableTypeLong:
			if identifier.Variable.LocalValOffset == 0 {
				op = []byte{cg.OP_lstore_0}
			} else if identifier.Variable.LocalValOffset == 1 {
				op = []byte{cg.OP_lstore_1}
			} else if identifier.Variable.LocalValOffset == 2 {
				op = []byte{cg.OP_lstore_2}
			} else if identifier.Variable.LocalValOffset == 3 {
				op = []byte{cg.OP_lstore_3}
			} else if identifier.Variable.LocalValOffset <= 255 {
				op = []byte{cg.OP_lstore, byte(identifier.Variable.LocalValOffset)}
			} else {
				panic("over 255")
			}
		default: // must be a object type
			if identifier.Variable.LocalValOffset == 0 {
				op = []byte{cg.OP_astore_0}
			} else if identifier.Variable.LocalValOffset == 1 {
				op = []byte{cg.OP_astore_1}
			} else if identifier.Variable.LocalValOffset == 2 {
				op = []byte{cg.OP_astore_2}
			} else if identifier.Variable.LocalValOffset == 3 {
				op = []byte{cg.OP_astore_3}
			} else if identifier.Variable.LocalValOffset <= 255 {
				op = []byte{cg.OP_astore, byte(identifier.Variable.LocalValOffset)}
			} else {
				panic("over 255")
			}
		}
		target = identifier.Variable.Type
	case ast.ExpressionTypeIndex:
		index := e.Data.(*ast.ExpressionIndex)
		if index.Expression.ExpressionValue.Type == ast.VariableTypeArray {
			meta := ArrayMetas[index.Expression.ExpressionValue.Array.Type]
			maxStack, _ = buildExpression.build(class, code, index.Expression, context, state)
			state.pushStack(class, index.Expression.ExpressionValue)
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Field:      "end",
				Descriptor: "I",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_swap
			code.CodeLength++
			code.Codes[code.CodeLength] = cg.OP_dup_x1
			code.CodeLength++
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Field:      "start",
				Descriptor: "I",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
			state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})

			stack, _ := buildExpression.build(class, code, index.Index, context, state)
			if t := stack + 3; t > maxStack {
				maxStack = t
			}
			code.Codes[code.CodeLength] = cg.OP_iadd
			code.CodeLength++
			code.Codes[code.CodeLength] = cg.OP_dup_x1
			code.CodeLength++
			{
				state.popStack(3)
				state.pushStack(class, state.newObjectVariableType(meta.className))
				state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
				context.MakeStackMap(code, state, code.CodeLength+6)
				context.MakeStackMap(code, state, code.CodeLength+16)
				state.popStack(2)
			}
			code.Codes[code.CodeLength] = cg.OP_if_icmple
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
			code.Codes[code.CodeLength+3] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 13)
			code.Codes[code.CodeLength+6] = cg.OP_pop // incase stack over flow
			code.Codes[code.CodeLength+7] = cg.OP_pop
			code.Codes[code.CodeLength+8] = cg.OP_new
			class.InsertClassConst(javaIndexOutOfRangeExceptionClass, code.Codes[code.CodeLength+9:code.CodeLength+11])
			code.Codes[code.CodeLength+11] = cg.OP_dup
			code.Codes[code.CodeLength+12] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      javaIndexOutOfRangeExceptionClass,
				Method:     specialMethodInit,
				Descriptor: "()V",
			}, code.Codes[code.CodeLength+13:code.CodeLength+15])
			code.Codes[code.CodeLength+15] = cg.OP_athrow
			// index not out of range
			code.Codes[code.CodeLength+16] = cg.OP_swap
			code.Codes[code.CodeLength+17] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      meta.className,
				Field:      "elements",
				Descriptor: meta.elementsFieldDescriptor,
			}, code.Codes[code.CodeLength+18:code.CodeLength+20])
			code.CodeLength += 20
			code.Codes[code.CodeLength] = cg.OP_swap
			code.CodeLength++
			{
				t := &ast.Type{}
				t.Type = ast.VariableTypeJavaArray
				t.Array = index.Expression.ExpressionValue.Array
				state.pushStack(class, t)
				state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
			}
			switch e.ExpressionValue.Type {
			case ast.VariableTypeBool:
				op = []byte{cg.OP_bastore}
			case ast.VariableTypeByte:
				op = []byte{cg.OP_bastore}
			case ast.VariableTypeShort:
				op = []byte{cg.OP_sastore}
			case ast.VariableTypeEnum:
				fallthrough
			case ast.VariableTypeInt:
				op = []byte{cg.OP_iastore}
			case ast.VariableTypeLong:
				op = []byte{cg.OP_lastore}
			case ast.VariableTypeFloat:
				op = []byte{cg.OP_fastore}
			case ast.VariableTypeDouble:
				op = []byte{cg.OP_dastore}
			case ast.VariableTypeFunction:
				fallthrough
			case ast.VariableTypeString:
				fallthrough
			case ast.VariableTypeObject:
				fallthrough
			case ast.VariableTypeMap:
				fallthrough
			case ast.VariableTypeArray:
				fallthrough
			case ast.VariableTypeJavaArray:
				op = []byte{cg.OP_aastore}
			}
			remainStack = 2 // [arrayref ,index]
			target = e.ExpressionValue
		} else if index.Expression.ExpressionValue.Type == ast.VariableTypeMap { // map
			return buildExpression.getMapLeftValue(class, code, e, context, state)
		} else { // java array
			maxStack, _ = buildExpression.build(class, code, index.Expression, context, state)
			state.pushStack(class, index.Expression.ExpressionValue)
			stack, _ := buildExpression.build(class, code, index.Index, context, state)
			if t := stack + 1; t > maxStack {
				maxStack = t
			}
			target = e.ExpressionValue
			remainStack = 2 // [objectref ,index]
			state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
			switch e.ExpressionValue.Type {
			case ast.VariableTypeBool:
				op = []byte{cg.OP_bastore}
			case ast.VariableTypeByte:
				op = []byte{cg.OP_bastore}
			case ast.VariableTypeShort:
				op = []byte{cg.OP_sastore}
			case ast.VariableTypeEnum:
				fallthrough
			case ast.VariableTypeInt:
				op = []byte{cg.OP_iastore}
			case ast.VariableTypeLong:
				op = []byte{cg.OP_lastore}
			case ast.VariableTypeFloat:
				op = []byte{cg.OP_fastore}
			case ast.VariableTypeDouble:
				op = []byte{cg.OP_dastore}
			case ast.VariableTypeFunction:
				fallthrough
			case ast.VariableTypeString:
				fallthrough
			case ast.VariableTypeObject:
				fallthrough
			case ast.VariableTypeMap:
				fallthrough
			case ast.VariableTypeArray:
				fallthrough
			case ast.VariableTypeJavaArray:
				op = []byte{cg.OP_aastore}
			}
			return
		}
	case ast.ExpressionTypeSelection:
		selection := e.Data.(*ast.ExpressionSelection)
		if selection.Expression.ExpressionValue.Type == ast.VariableTypePackage {
			op = []byte{cg.OP_putstatic}
			target = selection.PackageVariable.Type
			className = selection.Expression.ExpressionValue.Package.Name + "/main"
			name = selection.PackageVariable.Name
			if selection.PackageVariable.JvmDescriptor == "" {
				selection.PackageVariable.JvmDescriptor = JvmDescriptor.typeDescriptor(e.ExpressionValue)
			}
			descriptor = selection.PackageVariable.JvmDescriptor
			maxStack = 0
			remainStack = 0
		} else {
			className = selection.Expression.ExpressionValue.Class.Name
			target = selection.Field.Variable.Type
			name = selection.Name
			if selection.Field.LoadFromOutSide {
				descriptor = selection.Field.JvmDescriptor
			} else {
				descriptor = JvmDescriptor.typeDescriptor(target)
			}
			if selection.Field.IsStatic() {
				op = []byte{cg.OP_putstatic}
			} else {
				maxStack, _ = buildExpression.build(class, code, selection.Expression, context, state)
				remainStack = 1
				state.pushStack(class, selection.Expression.ExpressionValue)
				op = []byte{cg.OP_putfield}
			}
		}
	}
	return
}
