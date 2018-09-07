package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildFunctionExpression(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	function := e.Data.(*ast.Function)
	defer func(function *ast.Function) {
		if e.IsStatementExpression {
			return
		}
		stack := buildPackage.packFunction2MethodHandle(class, code, function, context)
		if stack > maxStack {
			maxStack = stack
		}
	}(function)
	if function.Name == "" {
		function.Name = function.NameLiteralFunction()
	}
	if function.IsClosureFunction == false {
		function.Name = class.NewFunctionName(function.Name) // new a function name
		method := &cg.MethodHighLevel{}
		method.Name = function.Name
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_BRIDGE
		if function.Type.VArgs != nil {
			method.AccessFlags |= cg.ACC_METHOD_VARARGS
		}
		function.ClassMethod = method
		method.Class = class
		method.Descriptor = Descriptor.methodDescriptor(&function.Type)
		method.Code = &cg.AttributeCode{}
		buildPackage.buildFunction(class, nil, method, function)
		class.AppendMethod(method)
		return
	}

	// function have captured vars
	className := buildPackage.newClassName("closureFunction$" + function.Name)
	closureClass := &cg.ClassHighLevel{}
	closureClass.Name = className
	closureClass.SuperClass = ast.LucyRootClass
	closureClass.AccessFlags = 0
	closureClass.AccessFlags |= cg.ACC_CLASS_SYNTHETIC
	closureClass.AccessFlags |= cg.ACC_CLASS_FINAL
	buildPackage.mkClassDefaultConstruction(closureClass)
	buildPackage.putClass(closureClass)
	method := &cg.MethodHighLevel{}
	method.Name = function.Name
	method.AccessFlags |= cg.ACC_METHOD_FINAL
	method.AccessFlags |= cg.ACC_METHOD_BRIDGE
	if function.Type.VArgs != nil {
		method.AccessFlags |= cg.ACC_METHOD_VARARGS
	}
	method.Descriptor = Descriptor.methodDescriptor(&function.Type)
	method.Class = closureClass
	function.ClassMethod = method
	closureClass.AppendMethod(method)
	//new a object to hold this closure function
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(closureClass.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxStack = 2 // maxStack is 2 right now
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      className,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	// store  to,wait for call
	function.ClosureVariableOffSet = code.MaxLocals
	code.MaxLocals++
	state.appendLocals(class, state.newObjectVariableType(className))
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, function.ClosureVariableOffSet)...)
	//set filed
	closureClass.Fields = make(map[string]*cg.FieldHighLevel)
	total := len(function.Closure.Variables) + len(function.Closure.Functions) - 1
	for v, _ := range function.Closure.Variables {
		filed := &cg.FieldHighLevel{}
		filed.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		filed.Name = v.Name
		closureClass.Fields[v.Name] = filed
		if total != 0 {
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
		}
		meta := closure.getMeta(v.Type.Type)
		filed.Descriptor = "L" + meta.className + ";"
		if context.function.Closure.ClosureVariableExist(v) {
			// I Know class at 0 offset
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			if 3 > maxStack {
				maxStack = 3
			}
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      v.Name,
				Descriptor: filed.Descriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else { // not exits
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, v.LocalValOffset)...)
			if 3 > maxStack {
				maxStack = 3
			}
		}
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      className,
			Field:      v.Name,
			Descriptor: filed.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		total--
	}
	for v, _ := range function.Closure.Functions {
		filed := &cg.FieldHighLevel{}
		filed.AccessFlags |= cg.ACC_FIELD_PUBLIC
		filed.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		filed.Name = v.Name
		filed.Descriptor = "L" + v.ClassMethod.Class.Name + ";"
		closureClass.Fields[v.Name] = filed
		if total != 0 {
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
		}
		if context.function.Closure.ClosureFunctionExist(v) {
			// I Know class at 0 offset
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			if 3 > maxStack {
				maxStack = 3
			}
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      v.Name,
				Descriptor: filed.Descriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else { // not exits
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, v.ClosureVariableOffSet)...)
			if 3 > maxStack {
				maxStack = 3
			}
		}
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      className,
			Field:      v.Name,
			Descriptor: filed.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		total--
	}
	method.Code = &cg.AttributeCode{}
	// build function
	buildPackage.buildFunction(closureClass, nil, method, function)
	return

}
func (buildPackage *BuildPackage) packFunction2MethodHandle(class *cg.ClassHighLevel, code *cg.AttributeCode,
	function *ast.Function, context *Context) (maxStack uint16) {
	code.Codes[code.CodeLength] = cg.OP_invokestatic
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/invoke/MethodHandles",
		Method:     "lookup",
		Descriptor: "()Ljava/lang/invoke/MethodHandles$Lookup;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertClassConst(function.ClassMethod.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst(function.ClassMethod.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertMethodTypeConst(cg.CONSTANT_MethodType_info_high_level{
		Descriptor: function.ClassMethod.Descriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	if function.IsClosureFunction {
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/invoke/MethodHandles$Lookup",
			Method:     "findVirtual",
			Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/invoke/MethodHandles$Lookup",
			Method:     "findStatic",
			Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	}
	code.CodeLength += 3
	if 4 > maxStack {
		maxStack = 4
	}
	if function.IsClosureFunction {
		if context.function.Closure.ClosureFunctionExist(function) {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      function.Name,
				Descriptor: "L" + function.ClassMethod.Class.Name + ";",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, function.ClosureVariableOffSet)...)
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
