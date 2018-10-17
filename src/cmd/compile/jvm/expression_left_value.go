package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) getCaptureIdentifierLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxStack, remainStack uint16, ops []byte,
	leftValueType LeftValueKind) {
	identifier := e.Data.(*ast.ExpressionIdentifier)
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
	leftValueType = LeftValueKindPutField
	return
}

func (buildExpression *BuildExpression) getMapLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (
	maxStack, remainStack uint16, ops []byte,
	leftValueType LeftValueKind) {
	index := e.Data.(*ast.ExpressionIndex)
	maxStack = buildExpression.build(class, code, index.Expression, context, state)
	state.pushStack(class, state.newObjectVariableType(mapClass))
	stack := buildExpression.build(class, code, index.Index, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	if index.Index.Value.IsPointer() == false {
		typeConverter.packPrimitives(class, code, index.Index.Value)
	}
	state.pushStack(class, state.newObjectVariableType(javaRootClass))
	remainStack = 2
	ops = []byte{}
	if index.Expression.Value.Map.V.IsPointer() == false {
		ops = append(ops,
			typeConverter.packPrimitivesBytes(class, index.Expression.Value.Map.V)...)
	}
	bs4 := make([]byte, 4)
	bs4[0] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      mapClass,
		Method:     "put",
		Descriptor: "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
	}, bs4[1:3])
	bs4[3] = cg.OP_pop
	ops = append(ops, bs4...)
	leftValueType = LeftValueKindMap
	return
}

func (buildExpression *BuildExpression) getLeftValue(
	class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (
	maxStack, remainStack uint16, ops []byte,
	leftValueType LeftValueKind) {
	switch e.Type {
	case ast.ExpressionTypeIdentifier:
		identifier := e.Data.(*ast.ExpressionIdentifier)
		if identifier.Name == ast.NoNameIdentifier {
			panic("this is not happening")
		}
		if identifier.Variable.IsGlobal {
			ops = make([]byte, 3)
			leftValueType = LeftValueKindPutStatic
			ops[0] = cg.OP_putstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      buildExpression.BuildPackage.mainClass.Name,
				Field:      identifier.Name,
				Descriptor: Descriptor.typeDescriptor(identifier.Variable.Type),
			}, ops[1:3])
			return
		}
		if identifier.Variable.BeenCapturedAsLeftValue > 0 {
			return buildExpression.getCaptureIdentifierLeftValue(class, code, e, context, state)
		}
		leftValueType = LeftValueKindLocalVar
		ops = storeLocalVariableOps(identifier.Variable.Type.Type, identifier.Variable.LocalValOffset)
	case ast.ExpressionTypeIndex:
		index := e.Data.(*ast.ExpressionIndex)
		if index.Expression.Value.Type == ast.VariableTypeArray {
			maxStack = buildExpression.build(class, code, index.Expression, context, state)
			state.pushStack(class, index.Expression.Value)
			stack := buildExpression.build(class, code, index.Index, context, state)
			if t := stack + 1; t > maxStack {
				maxStack = t
			}
			meta := ArrayMetas[index.Expression.Value.Array.Type]
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
			leftValueType = LeftValueKindLucyArray
			remainStack = 2 // [arrayref ,index]
		} else if index.Expression.Value.Type == ast.VariableTypeMap { // map
			return buildExpression.getMapLeftValue(class, code, e, context, state)
		} else { // java array
			maxStack = buildExpression.build(class, code, index.Expression, context, state)
			state.pushStack(class, index.Expression.Value)
			stack := buildExpression.build(class, code, index.Index, context, state)
			if t := stack + 1; t > maxStack {
				maxStack = t
			}
			leftValueType = LeftValueKindArray
			remainStack = 2 // [objectref ,index]
			state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
			switch e.Value.Type {
			case ast.VariableTypeBool:
				ops = []byte{cg.OP_bastore}
			case ast.VariableTypeByte:
				ops = []byte{cg.OP_bastore}
			case ast.VariableTypeShort:
				ops = []byte{cg.OP_sastore}
			case ast.VariableTypeChar:
				ops = []byte{cg.OP_castore}
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
			return
		}
	case ast.ExpressionTypeSelection:
		selection := e.Data.(*ast.ExpressionSelection)
		switch selection.Expression.Value.Type {
		case ast.VariableTypePackage:
			ops = make([]byte, 3)
			ops[0] = cg.OP_putstatic
			if selection.PackageVariable.JvmDescriptor == "" {
				selection.PackageVariable.JvmDescriptor = Descriptor.typeDescriptor(selection.PackageVariable.Type)
			}
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      selection.Expression.Value.Package.Name + "/main",
				Field:      selection.PackageVariable.Name,
				Descriptor: selection.PackageVariable.JvmDescriptor,
			}, ops[1:3])
			maxStack = 0
			leftValueType = LeftValueKindPutStatic
			remainStack = 0
		case ast.VariableTypeDynamicSelector:
			ops = make([]byte, 3)
			if selection.Field.IsStatic() {
				ops[0] = cg.OP_putstatic
				leftValueType = LeftValueKindPutStatic
			} else {
				code.Codes[code.CodeLength] = cg.OP_aload_0
				code.CodeLength++
				state.pushStack(class, state.newObjectVariableType(selection.Expression.Value.Class.Name))
				ops[0] = cg.OP_putfield
				remainStack = 1
				maxStack = 1
				leftValueType = LeftValueKindPutField
			}
			if selection.Field.JvmDescriptor == "" {
				selection.Field.JvmDescriptor = Descriptor.typeDescriptor(selection.Field.Type)
			}
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      selection.Expression.Value.Class.Name,
				Field:      selection.Name,
				Descriptor: selection.Field.JvmDescriptor,
			}, ops[1:3])
		default:
			ops = make([]byte, 3)
			if selection.Field.JvmDescriptor == "" {
				selection.Field.JvmDescriptor = Descriptor.typeDescriptor(selection.Field.Type)
			}
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      selection.Expression.Value.Class.Name,
				Field:      selection.Name,
				Descriptor: selection.Field.JvmDescriptor,
			}, ops[1:3])
			if selection.Field.IsStatic() {
				leftValueType = LeftValueKindPutStatic
				ops[0] = cg.OP_putstatic
			} else {
				leftValueType = LeftValueKindPutField
				ops[0] = cg.OP_putfield
				maxStack = buildExpression.build(class, code, selection.Expression, context, state)
				remainStack = 1
				state.pushStack(class, selection.Expression.Value)
			}
		}
	}
	return
}
