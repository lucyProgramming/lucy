package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeClass *MakeClass) buildFunctionExpression(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	function := e.Data.(*ast.Function)
	if function.IsClosureFunction == false {
		function.Name = class.NewFunctionName(function.Name) // new a function name
		method := &cg.MethodHighLevel{}
		method.Name = function.Name
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_BRIDGE
		function.ClassMethod = method
		method.Class = class
		method.Descriptor = Descriptor.methodDescriptor(function)
		method.Code = &cg.AttributeCode{}
		makeClass.buildFunction(class, nil, method, function)
		class.AppendMethod(method)
		return
	}

	// function have captured vars
	className := makeClass.newClassName("closureFunction_" + function.Name)
	closureClass := &cg.ClassHighLevel{}
	closureClass.Name = className
	closureClass.SuperClass = ast.LUCY_ROOT_CLASS
	closureClass.AccessFlags = 0
	closureClass.Class.AttributeCompilerAuto = &cg.AttributeCompilerAuto{}
	closureClass.AccessFlags |= cg.ACC_CLASS_SYNTHETIC
	closureClass.AccessFlags |= cg.ACC_CLASS_FINAL
	makeClass.mkClassDefaultConstruction(closureClass, nil)
	makeClass.putClass(className, closureClass)

	method := &cg.MethodHighLevel{}
	method.Name = function.Name
	method.AccessFlags |= cg.ACC_METHOD_FINAL
	method.AccessFlags |= cg.ACC_METHOD_PUBLIC
	method.Descriptor = Descriptor.methodDescriptor(function)
	method.Class = closureClass
	function.ClassMethod = method

	closureClass.AppendMethod(method)
	//new a object to hold this closure function
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(closureClass.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	maxStack = 2 // maxstack is 2 right now
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      className,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	// store  to,wait for call
	function.VarOffSet = code.MaxLocals
	code.MaxLocals++
	state.appendLocals(class, state.newObjectVariableType(className))
	copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, function.VarOffSet)...)
	//set filed
	closureClass.Fields = make(map[string]*cg.FieldHighLevel)
	total := len(function.Closure.Variables) + len(function.Closure.Functions)
	i := 0
	for v, _ := range function.Closure.Variables {
		filed := &cg.FieldHighLevel{}
		filed.AccessFlags |= cg.ACC_FIELD_PUBLIC
		filed.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		filed.Name = v.Name
		meta := closure.getMeta(v.Type.Type)
		filed.Descriptor = "L" + meta.className + ";"
		closureClass.Fields[v.Name] = filed
		if i != total-1 {
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
		}
		if context.function.Closure.ClosureVariableExist(v) {
			// I Know class at 0 offset
			copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, 0)...)
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
			copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
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
		i++
	}
	for v, _ := range function.Closure.Functions {
		filed := &cg.FieldHighLevel{}
		filed.AccessFlags |= cg.ACC_FIELD_PUBLIC
		filed.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		filed.Name = v.Name
		filed.Descriptor = "L" + v.ClassMethod.Class.Name + ";"
		closureClass.Fields[v.Name] = filed
		if i != total-1 {
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
		}
		if context.function.Closure.ClosureFunctionExist(v) {
			// I Know class at 0 offset
			copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, 0)...)
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
			copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, v.VarOffSet)...)
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
		i++
	}
	method.Code = &cg.AttributeCode{}
	// build function
	makeClass.buildFunction(closureClass, nil, method, function)
	return

}
