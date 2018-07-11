package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) getCaptureIdentifierLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxStack, remainStack uint16, ops []byte,
	target *ast.Type, leftValueType int) {
	identifier := e.Data.(*ast.ExpressionIdentifier)
	target = identifier.Variable.Type
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
	ops = make([]byte, 3)
	ops[0] = cg.OP_putfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      meta.className,
		Field:      meta.fieldName,
		Descriptor: meta.fieldDescription,
	}, ops[1:3])
	leftValueType = LeftValueTypePutField
	return
}

func (buildExpression *BuildExpression) getMapLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxStack, remainStack uint16, ops []byte,
	target *ast.Type, leftValueType int) {
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
	ops = []byte{}
	if index.Expression.ExpressionValue.Map.V.IsPointer() == false {
		ops = append(ops,
			typeConverter.packPrimitivesBytes(class, index.Expression.ExpressionValue.Map.V)...)
	}
	bs4 := make([]byte, 4)
	bs4[0] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaMapClass,
		Method:     "put",
		Descriptor: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
	}, bs4[1:3])
	bs4[3] = cg.OP_pop
	ops = append(ops, bs4...)
	target = index.Expression.ExpressionValue.Map.V
	leftValueType = LeftValueTypeMap
	return
}

func (buildExpression *BuildExpression) getLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (
	maxStack, remainStack uint16, ops []byte,
	target *ast.Type, leftValueType int) {
	switch e.Type {
	case ast.ExpressionTypeIdentifier:
		identifier := e.Data.(*ast.ExpressionIdentifier)
		if identifier.Variable.IsGlobal {
			ops = make([]byte, 3)
			leftValueType = LeftValueTypePutStatic
			ops[0] = cg.OP_putstatic
			target = identifier.Variable.Type
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      buildExpression.BuildPackage.mainClass.Name,
				Field:      identifier.Name,
				Descriptor: Descriptor.typeDescriptor(identifier.Variable.Type),
			}, ops[1:3])
			return
		}
		if identifier.Variable.BeenCaptured {
			return buildExpression.getCaptureIdentifierLeftValue(class, code, e, context, state)
		}
		if identifier.Name == ast.NoNameIdentifier {
			panic("this is not happening")
		}
		leftValueType = LeftValueTypeLocalVar
		ops = storeLocalVariableOps(identifier.Variable.Type.Type, identifier.Variable.LocalValOffset)
		target = identifier.Variable.Type
	case ast.ExpressionTypeIndex:
		index := e.Data.(*ast.ExpressionIndex)
		if index.Expression.ExpressionValue.Type == ast.VariableTypeArray {
			maxStack, _ = buildExpression.build(class, code, index.Expression, context, state)
			state.pushStack(class, index.Expression.ExpressionValue)
			stack, _ := buildExpression.build(class, code, index.Index, context, state)
			if t := stack + 1; t > maxStack {
				maxStack = t
			}
			meta := ArrayMetas[index.Expression.ExpressionValue.Array.Type]
			ops = make([]byte, 3)
			ops[0] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.className,
				Method:     "set",
				Descriptor: meta.setMethodDescription,
			}, ops[1:3])
			state.pushStack(class, &ast.Type{
				Type: ast.VariableTypeInt,
			})
			leftValueType = LeftValueTypeLucyArray
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
				ops = []byte{cg.OP_bastore}
			case ast.VariableTypeByte:
				ops = []byte{cg.OP_bastore}
			case ast.VariableTypeShort:
				ops = []byte{cg.OP_sastore}
			case ast.VariableTypeEnum:
				fallthrough
			case ast.VariableTypeInt:
				ops = []byte{cg.OP_iastore}
			case ast.VariableTypeLong:
				ops = []byte{cg.OP_lastore}
			case ast.VariableTypeFloat:
				ops = []byte{cg.OP_fastore}
			case ast.VariableTypeDouble:
				ops = []byte{cg.OP_dastore}
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
				ops = []byte{cg.OP_aastore}
			}
			leftValueType = LeftValueTypeArray
			return
		}
	case ast.ExpressionTypeSelection:
		selection := e.Data.(*ast.ExpressionSelection)
		if selection.Expression.ExpressionValue.Type == ast.VariableTypePackage {
			ops = make([]byte, 3)
			ops[0] = cg.OP_putstatic
			target = selection.PackageVariable.Type
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      selection.Expression.ExpressionValue.Package.Name + "/main",
				Field:      selection.PackageVariable.Name,
				Descriptor: selection.PackageVariable.JvmDescriptor,
			}, ops[1:3])
			maxStack = 0
			leftValueType = LeftValueTypePutStatic
			remainStack = 0
		} else {
			target = selection.Field.Variable.Type
			ops = make([]byte, 3)
			if selection.Field.JvmDescriptor == "" {
				selection.Field.JvmDescriptor = Descriptor.typeDescriptor(target)
			}
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      selection.Expression.ExpressionValue.Class.Name,
				Field:      selection.Name,
				Descriptor: selection.Field.JvmDescriptor,
			}, ops[1:3])
			if selection.Field.IsStatic() {
				leftValueType = LeftValueTypePutStatic
				ops[0] = cg.OP_putstatic
			} else {
				leftValueType = LeftValueTypePutField
				ops[0] = cg.OP_putfield
				maxStack, _ = buildExpression.build(class, code, selection.Expression, context, state)
				remainStack = 1
				state.pushStack(class, selection.Expression.ExpressionValue)
			}
		}
	}
	return
}
