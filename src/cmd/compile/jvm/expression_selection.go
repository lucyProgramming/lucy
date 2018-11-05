package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

func (this *BuildExpression) buildSelection(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	selection := e.Data.(*ast.ExpressionSelection)

	// check cast to super class
	if selection.Name == ast.SUPER {
		maxStack = this.build(class, code, selection.Expression, context, state)
		return
	}
	if selection.Method != nil { // pack to method handle
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      "java/lang/invoke/MethodHandles",
			Method:     "lookup",
			Descriptor: "()Ljava/lang/invoke/MethodHandles$Lookup;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertClassConst(selection.Expression.Value.Class.Name,
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst(selection.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertMethodTypeConst(cg.ConstantInfoMethodTypeHighLevel{
			Descriptor: Descriptor.methodDescriptor(&selection.Method.Function.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		if selection.Method.IsStatic() {
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
				Class:      "java/lang/invoke/MethodHandles$Lookup",
				Method:     "findStatic",
				Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		} else {
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
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
			stack := this.build(class, code, selection.Expression, context, state)
			if t := 1 + stack; t > maxStack {
				maxStack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
				Class:      "java/lang/invoke/MethodHandle",
				Method:     "bindTo",
				Descriptor: "(Ljava/lang/Object;)Ljava/lang/invoke/MethodHandle;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}

	switch selection.Expression.Value.Type {
	case ast.VariableTypePackage:
		if selection.PackageVariable != nil {
			maxStack = jvmSlotSize(e.Value)
			if selection.PackageVariable.JvmDescriptor == "" {
				selection.PackageVariable.JvmDescriptor = Descriptor.typeDescriptor(e.Value)
			}
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
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
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
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
			class.InsertMethodTypeConst(cg.ConstantInfoMethodTypeHighLevel{
				Descriptor: Descriptor.methodDescriptor(&selection.PackageFunction.Type),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
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
	case ast.VariableTypeDynamicSelector:
		if selection.Field != nil {
			if selection.Field.IsStatic() == false {
				code.Codes[code.CodeLength] = cg.OP_aload_0
				code.CodeLength++
				if 1 > maxStack {
					maxStack = 1
				}
				code.Codes[code.CodeLength] = cg.OP_getfield
				code.CodeLength++
			} else {
				code.Codes[code.CodeLength] = cg.OP_getstatic
				code.CodeLength++
			}
			if selection.Field.JvmDescriptor == "" {
				selection.Field.JvmDescriptor = Descriptor.typeDescriptor(selection.Field.Type)
			}
			class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
				Class:      selection.Expression.Value.Class.Name,
				Field:      selection.Name,
				Descriptor: selection.Field.JvmDescriptor,
			},
				code.Codes[code.CodeLength:code.CodeLength+2])
			code.CodeLength += 2
		} else {
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
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
			class.InsertMethodTypeConst(cg.ConstantInfoMethodTypeHighLevel{
				Descriptor: Descriptor.methodDescriptor(&selection.Method.Function.Type),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			if selection.Method.IsStatic() {
				class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
					Class:      "java/lang/invoke/MethodHandles$Lookup",
					Method:     "findStatic",
					Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
				}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			} else {
				class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
					Class:      "java/lang/invoke/MethodHandles$Lookup",
					Method:     "findVirtual",
					Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
				}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			}
			code.CodeLength += 3
			if 4 > maxStack {
				maxStack = 4
			}
			if selection.Method.IsStatic() == false {
				code.Codes[code.CodeLength] = cg.OP_aload_0
				code.CodeLength++
				code.Codes[code.CodeLength] = cg.OP_invokevirtual
				class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
					Class:      "java/lang/invoke/MethodHandle",
					Method:     "bindTo",
					Descriptor: "(Ljava/lang/Object;)Ljava/lang/invoke/MethodHandle;",
				}, code.Codes[code.CodeLength+1:code.CodeLength+3])
				code.CodeLength += 3
			}
		}
		return
	case ast.VariableTypeClass:
		maxStack = jvmSlotSize(e.Value)
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
			Class:      selection.Expression.Value.Class.Name,
			Field:      selection.Name,
			Descriptor: Descriptor.typeDescriptor(e.Value),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	case ast.VariableTypeObject:
		// object
		stack := this.build(class, code, selection.Expression, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
			Class:      selection.Expression.Value.Class.Name,
			Field:      selection.Name,
			Descriptor: Descriptor.typeDescriptor(e.Value),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := jvmSlotSize(e.Value); t > maxStack {
			maxStack = t
		}
		return
	}
	return

}
