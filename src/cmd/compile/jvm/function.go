package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) mkParametersOffset(class *cg.ClassHighLevel, code *cg.AttributeCode,
	f *ast.Function, state *StackMapState) {
	for _, v := range f.Type.ParameterList { // insert into locals
		v.LocalValOffset = code.MaxLocals
		code.MaxLocals += jvmSlotSize(v.Type)
		state.appendLocals(class, v.Type)
	}
	if f.Type.VArgs != nil {
		f.Type.VArgs.LocalValOffset = code.MaxLocals
		code.MaxLocals++
		state.appendLocals(class, f.Type.VArgs.Type)
	}
}

func (buildPackage *BuildPackage) mkCapturedParameters(class *cg.ClassHighLevel, code *cg.AttributeCode,
	f *ast.Function, state *StackMapState) (maxStack uint16) {
	for _, v := range f.Type.ParameterList {
		if v.BeenCaptured == 0 { // not capture
			continue
		}
		stack := closure.createClosureVar(class, code, v.Type)
		if stack > maxStack {
			maxStack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		if t := 2 + jvmSlotSize(v.Type); t > maxStack {
			maxStack = t
		}
		copyOPs(code, loadLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
		buildPackage.storeLocalVar(class, code, v)
		v.LocalValOffset = code.MaxLocals //rewrite offset
		code.MaxLocals++
		copyOPs(code, storeLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
		state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Type.Type).className))
	}
	return
}

func (buildPackage *BuildPackage) buildFunctionParameterAndReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode,
	f *ast.Function, context *Context, state *StackMapState) (maxStack uint16) {
	buildPackage.mkParametersOffset(class, code, f, state)
	maxStack = buildPackage.mkCapturedParameters(class, code, f, state)
	if f.Type.VoidReturn() == false {
		for _, v := range f.Type.ReturnList {
			currentStack := uint16(0)
			if v.BeenCaptured > 0 { //create closure object
				v.LocalValOffset = code.MaxLocals
				code.MaxLocals++
				stack := closure.createClosureVar(class, code, v.Type)
				if stack > maxStack {
					maxStack = stack
				}
				// then load
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				if 2 > maxStack {
					maxStack = 2
				}
				copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, v.LocalValOffset)...)
				currentStack = 1
				state.pushStack(class,
					state.newObjectVariableType(closure.getMeta(v.Type.Type).className))
			} else {
				v.LocalValOffset = code.MaxLocals
				code.MaxLocals += jvmSlotSize(v.Type)
			}
			stack := buildPackage.BuildExpression.build(class, code, v.DefaultValueExpression, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			buildPackage.storeLocalVar(class, code, v)
			if v.BeenCaptured > 0 {
				state.popStack(1)
				state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Type.Type).className))
			} else {
				state.appendLocals(class, v.Type)
			}
		}
	}
	return
}

func (buildPackage *BuildPackage) buildFunction(class *cg.ClassHighLevel, astClass *ast.Class, method *cg.MethodHighLevel,
	f *ast.Function) {
	context := &Context{}
	context.lastStackMapOffset = -1
	context.class = astClass
	context.function = f
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
	}()
	state := &StackMapState{}
	if method.AccessFlags&cg.ACC_METHOD_STATIC == 0 {
		if method.IsConstruction { // construction method
			method.Code.MaxLocals = 1
			t := &cg.StackMapVerificationTypeInfo{}
			t.Verify = &cg.StackMapUninitializedThisVariableInfo{}
			state.Locals = append(state.Locals, t)
			buildPackage.mkParametersOffset(class, method.Code, f, state)
			stack := buildPackage.BuildExpression.build(class, method.Code, f.CallFatherConstructionExpression,
				context, state)
			if stack > method.Code.MaxStack {
				method.Code.MaxStack = stack
			}
			state.Locals[0] = state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(class.Name))
			buildPackage.mkFieldDefaultValue(class, method.Code, astClass, context, state)
			buildPackage.mkCapturedParameters(class, method.Code, f, state)
		} else {
			method.Code.MaxLocals = 1
			state.appendLocals(class, state.newObjectVariableType(class.Name))
		}
	}
	if f.IsGlobalMain() { // main function
		code := method.Code
		code.Codes[code.CodeLength] = cg.OP_new
		meta := ArrayMetas[ast.VariableTypeString]
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeString, 0)...)
		if 3 > code.MaxStack {
			code.MaxStack = 3
		}
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     specialMethodInit,
			Descriptor: meta.constructorFuncDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, 1)...)
		{
			// String[] java style
			t := &ast.Type{Type: ast.VariableTypeJavaArray}
			t.Array = &ast.Type{Type: ast.VariableTypeString}
			state.appendLocals(class, t)
		}
		method.Code.MaxLocals = 1
	}
	if LucyMethodSignatureParser.Need(&f.Type) {
		d := &cg.AttributeLucyMethodDescriptor{}
		d.Descriptor = LucyMethodSignatureParser.Encode(&f.Type)
		method.AttributeLucyMethodDescriptor = d
	}
	if f.HaveDefaultValue {
		method.AttributeDefaultParameters = DefaultValueParser.Encode(class, f)
	}
	if method.IsConstruction == false {
		if t := buildPackage.buildFunctionParameterAndReturnList(class, method.Code, f, context, state); t > method.Code.MaxStack {
			method.Code.MaxStack = t
		}
	}
	{
		method.AttributeMethodParameters = &cg.AttributeMethodParameters{}
		for _, v := range f.Type.ParameterList {
			p := &cg.MethodParameter{}
			p.Name = v.Name
			p.AccessFlags = cg.METHOD_PARAMETER_TYPE_ACC_MANDATED
			method.AttributeMethodParameters.Parameters = append(method.AttributeMethodParameters.Parameters, p)
		}
	}
	if f.Type.VoidReturn() == false {
		method.AttributeLucyReturnListNames = &cg.AttributeMethodParameters{}
		for _, v := range f.Type.ReturnList {
			p := &cg.MethodParameter{}
			p.Name = v.Name
			p.AccessFlags = cg.METHOD_PARAMETER_TYPE_ACC_MANDATED
			method.AttributeLucyReturnListNames.Parameters =
				append(method.AttributeLucyReturnListNames.Parameters, p)
		}
	}
	if len(f.Type.ReturnList) > 1 {
		if t := buildPackage.buildFunctionMultiReturnOffset(class, method.Code,
			f, context, state); t > method.Code.MaxStack {
			method.Code.MaxStack = t
		}
	}
	if f.HasDefer {
		context.exceptionVarOffset = method.Code.MaxLocals
		method.Code.MaxLocals++
		method.Code.Codes[method.Code.CodeLength] = cg.OP_aconst_null
		method.Code.CodeLength++
		copyOPs(method.Code, storeLocalVariableOps(ast.VariableTypeObject, context.exceptionVarOffset)...)
		state.appendLocals(class, state.newObjectVariableType(ast.JavaThrowableClass))
	}
	buildPackage.buildBlock(class, method.Code, &f.Block, context, state)
	return
}
func (buildPackage *BuildPackage) buildFunctionMultiReturnOffset(class *cg.ClassHighLevel, code *cg.AttributeCode,
	f *ast.Function, context *Context, state *StackMapState) (maxStack uint16) {
	code.Codes[code.CodeLength] = cg.OP_aconst_null
	code.CodeLength++
	context.multiValueVarOffset = code.MaxLocals
	code.MaxLocals++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
		context.multiValueVarOffset)...)
	state.appendLocals(class, state.newObjectVariableType(javaRootObjectArray))
	maxStack = 1
	return
}

func (buildPackage *BuildPackage) mkFieldDefaultValue(class *cg.ClassHighLevel, code *cg.AttributeCode,
	astClass *ast.Class, context *Context, state *StackMapState) {
	for _, v := range astClass.Fields {
		if v.IsStatic() || v.DefaultValueExpression == nil {
			continue
		}
		code.Codes[code.CodeLength] = cg.OP_aload_0
		code.CodeLength++
		state.pushStack(class, state.newObjectVariableType(class.Name))
		stack := buildPackage.BuildExpression.build(class, code, v.DefaultValueExpression, context, state)
		if t := 1 + stack; t > code.MaxStack {
			code.MaxStack = t
		}
		state.popStack(1)
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      class.Name,
			Field:      v.Name,
			Descriptor: Descriptor.typeDescriptor(v.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
}
