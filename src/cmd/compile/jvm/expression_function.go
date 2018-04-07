package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildFunctionExpression(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	function := e.Data.(*ast.Function)
	if function.IsClosureFunction == false {
		function.Name = class.NewFunctionName(function.Name) // new a function name
		method := &cg.MethodHighLevel{}
		method.Name = function.Name
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		function.ClassMethod = method
		method.Class = class
		method.AttributeLucyInnerStaticMethod = &cg.AttributeLucyInnerStaticMethod{}
		method.Descriptor = Descriptor.methodDescriptor(function)
		m.buildFunction(class, method, function)
		class.AppendMethod(method)
		return
	}

	// function have captured vars
	classname := m.newClassName("closureFunction_" + function.Name)
	closureClass := &cg.ClassHighLevel{}
	closureClass.Name = classname
	closureClass.SuperClass = ast.LUCY_ROOT_CLASS
	closureClass.AccessFlags = 0
	closureClass.Class.AttributeClosureClass = &cg.AttributeClosureFunctionClass{}
	closureClass.AccessFlags |= cg.ACC_CLASS_SYNTHETIC
	closureClass.AccessFlags |= cg.ACC_CLASS_FINAL
	mkClassDefaultContruction(closureClass)
	m.putClass(classname, closureClass)

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
	maxstack = 2 // maxstack is 2 right now
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      classname,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	// store  to,wait for call
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, function.VarOffSetForClosure)...)
	//set filed
	closureClass.Fields = make(map[string]*cg.FiledHighLevel)
	total := len(function.ClosureVars.Vars) + len(function.ClosureVars.Funcs)
	i := 0
	for v, _ := range function.ClosureVars.Vars {
		filed := &cg.FiledHighLevel{}
		filed.AccessFlags |= cg.ACC_FIELD_PUBLIC
		filed.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		filed.Name = v.Name
		meta := closure.getMeta(v.Typ.Typ)
		filed.Descriptor = "L" + meta.className + ";"
		closureClass.Fields[v.Name] = filed
		if i != total-1 {
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
		}
		if context.function.ClosureVars.ClosureVariableExist(v) {
			// I Know class at 0 offset
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, 0)...)
			if 3 > maxstack {
				maxstack = 3
			}
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      v.Name,
				Descriptor: filed.Descriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else { // not exits
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
			if 3 > maxstack {
				maxstack = 3
			}
		}
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      classname,
			Field:      v.Name,
			Descriptor: filed.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		i++
	}
	for v, _ := range function.ClosureVars.Funcs {
		filed := &cg.FiledHighLevel{}
		filed.AccessFlags |= cg.ACC_FIELD_PUBLIC
		filed.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		filed.Name = v.Name
		filed.Descriptor = "L" + v.ClassMethod.Class.Name + ";"
		closureClass.Fields[v.Name] = filed
		if i != total-1 {
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
		}
		if context.function.ClosureVars.ClosureFunctionExist(v) {
			// I Know class at 0 offset
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, 0)...)
			if 3 > maxstack {
				maxstack = 3
			}
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      v.Name,
				Descriptor: filed.Descriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else { // not exits
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.VarOffSetForClosure)...)
			if 3 > maxstack {
				maxstack = 3
			}
		}
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      classname,
			Field:      v.Name,
			Descriptor: filed.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		i++
	}

	// build function
	m.buildFunction(closureClass, method, function)
	return

}
