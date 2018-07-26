package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

func (buildExpression *BuildExpression) buildSelection(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	selection := e.Data.(*ast.ExpressionSelection)
	if selection.Expression.Value.Type == ast.VariableTypePackage {
		if selection.PackageVariable != nil {
			maxStack = jvmSlotSize(e.Value)
			if selection.PackageVariable.JvmDescriptor == "" {
				selection.PackageVariable.JvmDescriptor = Descriptor.typeDescriptor(e.Value)
			}
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      selection.Expression.Value.Package.Name + "/main",
				Field:      selection.PackageVariable.Name,
				Descriptor: selection.PackageVariable.JvmDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			return
		}
		if selection.PackageEnumName != nil {
			loadInt32(class, code, selection.PackageEnumName.Value)
			maxStack = 1
			return
		}
		if selection.PackageFunction != nil { // pack to method handle
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/invoke/MethodHandles",
				Method:     "lookup",
				Descriptor: "()Ljava/lang/invoke/MethodHandles$Lookup;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertClassConst(selection.Expression.Value.Package.Name+"/main",
				code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertStringConst(selection.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertMethodTypeConst(cg.CONSTANT_MethodType_info_high_level{
				Descriptor: Descriptor.methodDescriptor(&selection.PackageFunction.Type),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/invoke/MethodHandles$Lookup",
				Method:     "findStatic",
				Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			if 4 > maxStack {
				maxStack = 4
			}
			return
		}
	}
	// check cast to super class
	if selection.Name == ast.SUPER {
		// no  need to cast to father
		//		if selection.Expression.ExpressionValue.Type == ast.VariableTypeObject {
		//			maxStack, _ = buildExpression.build(class, code, selection.Expression, context, state)
		//			code.Codes[code.CodeLength] = cg.OP_checkcast
		//			class.InsertClassConst(e.ExpressionValue.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		//			code.CodeLength += 3
		//		}
		return
	}

	if selection.Method != nil { // pack to method handle
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/invoke/MethodHandles",
			Method:     "lookup",
			Descriptor: "()Ljava/lang/invoke/MethodHandles$Lookup;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertClassConst(selection.Expression.Value.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst(selection.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertMethodTypeConst(cg.CONSTANT_MethodType_info_high_level{
			Descriptor: Descriptor.methodDescriptor(&selection.Method.Function.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		if selection.Method.IsStatic() {
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/invoke/MethodHandles$Lookup",
				Method:     "findStatic",
				Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		} else {
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/invoke/MethodHandles$Lookup",
				Method:     "findVirtual",
				Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		}
		code.CodeLength += 3
		if 4 > maxStack {
			maxStack = 4
		}
		if selection.Expression.Value.Type == ast.VariableTypeObject {
			stack := buildExpression.build(class, code, selection.Expression, context, state)
			if stack > maxStack {
				maxStack = stack
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/invoke/MethodHandle",
				Method:     "bindTo",
				Descriptor: "(Ljava/lang/Object;)Ljava/lang/invoke/MethodHandle;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}
	if selection.Expression.Value.Type == ast.VariableTypeClass {
		maxStack = jvmSlotSize(e.Value)
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      selection.Expression.Value.Class.Name,
			Field:      selection.Name,
			Descriptor: Descriptor.typeDescriptor(e.Value),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	// object
	maxStack = buildExpression.build(class, code, selection.Expression, context, state)
	if t := jvmSlotSize(e.Value); t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      selection.Expression.Value.Class.Name,
		Field:      selection.Name,
		Descriptor: Descriptor.typeDescriptor(e.Value),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
