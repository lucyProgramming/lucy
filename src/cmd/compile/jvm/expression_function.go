package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildFunctionExpression(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
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
		function.Name = class.NewMethodName(function.Name) // new a function name
		method := &cg.MethodHighLevel{}
		method.Name = function.Name
		method.AccessFlags |= cg.AccMethodFinal
		method.AccessFlags |= cg.AccMethodPrivate
		method.AccessFlags |= cg.AccMethodStatic
		method.AccessFlags |= cg.AccMethodBridge
		if function.Type.VArgs != nil {
			method.AccessFlags |= cg.AccMethodVarargs
		}
		function.Entrance = method
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
	closureClass.AccessFlags |= cg.AccClassSynthetic
	closureClass.AccessFlags |= cg.AccClassFinal
	buildPackage.mkClassDefaultConstruction(closureClass)
	buildPackage.putClass(closureClass)
	method := &cg.MethodHighLevel{}
	method.Name = function.Name
	method.AccessFlags |= cg.AccMethodFinal
	method.AccessFlags |= cg.AccMethodBridge
	if function.Type.VArgs != nil {
		method.AccessFlags |= cg.AccMethodVarargs
	}
	method.Descriptor = Descriptor.methodDescriptor(&function.Type)
	method.Class = closureClass
	function.Entrance = method
	closureClass.AppendMethod(method)
	//new a object to hold this closure function
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(closureClass.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxStack = 2 // maxStack is 2 right now
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
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
		field := &cg.FieldHighLevel{}
		field.AccessFlags |= cg.AccFieldSynthetic
		field.Name = v.Name
		closureClass.Fields[v.Name] = field
		if total != 0 {
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
		}
		if v.BeenCapturedAsLeftValue > 0 {
			meta := closure.getMeta(v.Type.Type)
			field.Descriptor = "L" + meta.className + ";"
			if context.function.Closure.ClosureVariableExist(v) {
				// I Know class object at offset 0
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
				if 3 > maxStack {
					maxStack = 3
				}
				code.Codes[code.CodeLength] = cg.OP_getfield
				class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
					Class:      class.Name,
					Field:      v.Name,
					Descriptor: field.Descriptor,
				}, code.Codes[code.CodeLength+1:code.CodeLength+3])
				code.CodeLength += 3
			} else { // not exits
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, v.LocalValOffset)...)
				if 3 > maxStack {
					maxStack = 3
				}
			}
		} else {
			field.Descriptor = Descriptor.typeDescriptor(v.Type)
			if context.function.Closure.ClosureVariableExist(v) {
				// I Know class object at offset 0
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
				if 3 > maxStack {
					maxStack = 3
				}
				code.Codes[code.CodeLength] = cg.OP_getfield
				class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
					Class:      class.Name,
					Field:      v.Name,
					Descriptor: field.Descriptor,
				}, code.Codes[code.CodeLength+1:code.CodeLength+3])
				code.CodeLength += 3
			} else { // not exits
				copyOPs(code, loadLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
				if 3 > maxStack {
					maxStack = 3
				}
			}
		}
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
			Class:      className,
			Field:      v.Name,
			Descriptor: field.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		total--
	}
	for v, _ := range function.Closure.Functions {
		if v.IsClosureFunction == false {
			continue
		}
		filed := &cg.FieldHighLevel{}
		filed.AccessFlags |= cg.AccFieldPublic
		filed.AccessFlags |= cg.AccFieldSynthetic
		filed.Name = v.Name
		filed.Descriptor = "L" + v.Entrance.Class.Name + ";"
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
			class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
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
		class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
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
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      "java/lang/invoke/MethodHandles",
		Method:     "lookup",
		Descriptor: "()Ljava/lang/invoke/MethodHandles$Lookup;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertClassConst(function.Entrance.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst(function.Entrance.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertMethodTypeConst(cg.ConstantInfoMethodTypeHighLevel{
		Descriptor: function.Entrance.Descriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	if function.IsClosureFunction {
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      "java/lang/invoke/MethodHandles$Lookup",
			Method:     "findVirtual",
			Descriptor: "(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/MethodHandle;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
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
			class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
				Class:      class.Name,
				Field:      function.Name,
				Descriptor: "L" + function.Entrance.Class.Name + ";",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, function.ClosureVariableOffSet)...)
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
