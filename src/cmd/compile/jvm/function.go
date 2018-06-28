package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeClass *MakeClass) buildFunctionParameterAndReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode, f *ast.Function, context *Context, state *StackMapState) (maxStack uint16) {
	for _, v := range f.Type.ParameterList { // insert into locals
		v.LocalValOffset = code.MaxLocals
		code.MaxLocals += jvmSlotSize(v.Type)
		state.appendLocals(class, v.Type)
	}
	for _, v := range f.Type.ParameterList {
		if v.BeenCaptured == false { // capture
			continue
		}
		stack := closure.createClosureVar(class, code, v.Type)
		if stack > maxStack {
			maxStack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		copyOPs(code, loadLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
		if t := 2 + jvmSlotSize(v.Type); t > maxStack {
			maxStack = t
		}
		makeClass.storeLocalVar(class, code, v)
		v.LocalValOffset = code.MaxLocals //rewrite offset
		code.MaxLocals++
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, v.LocalValOffset)...)
		state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Type.Type).className))
	}
	for _, v := range f.Type.ReturnList {
		currentStack := uint16(0)
		if v.BeenCaptured { //create closure object
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
		stack, es := makeClass.makeExpression.build(class, code, v.Expression, context, state)
		if len(es) > 0 {
			fillOffsetForExits(es, code.CodeLength)
			state.pushStack(class, v.Type)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		makeClass.storeLocalVar(class, code, v)
		if v.BeenCaptured {
			state.popStack(1)
			state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Type.Type).className))
		} else {
			state.appendLocals(class, v.Type)
		}
	}
	return
}

func (makeClass *MakeClass) buildFunction(class *cg.ClassHighLevel, astClass *ast.Class, method *cg.MethodHighLevel, f *ast.Function) {
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
	if method.IsConstruction { // construction method
		method.Code.MaxLocals = 1
		t := &cg.StackMapVerificationTypeInfo{}
		t.Verify = &cg.StackMapUninitializedThisVariableInfo{}
		state.Locals = append(state.Locals, t)
	} else if f.Name == ast.MainFunctionName { // main function
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
	} else if method.AccessFlags&cg.ACC_METHOD_STATIC == 0 { // instance method
		method.Code.MaxLocals = 1
		state.appendLocals(class, state.newObjectVariableType(class.Name))
	}
	if LucyMethodSignatureParser.Need(&f.Type) {
		d := &cg.AttributeLucyMethodDescriptor{}
		d.Descriptor = LucyMethodSignatureParser.Encode(&f.Type)
		method.AttributeLucyMethodDescriptor = d
	}
	if f.HaveDefaultValue {
		method.AttributeDefaultParameters = FunctionDefaultValueParser.Encode(class, f)
	}
	if t := makeClass.buildFunctionParameterAndReturnList(class, method.Code, f, context, state); t > method.Code.MaxStack {
		method.Code.MaxStack = t
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
	if f.NoReturnValue() == false {
		method.AttributeLucyReturnListNames = &cg.AttributeMethodParameters{}
		for _, v := range f.Type.ReturnList {
			p := &cg.MethodParameter{}
			p.Name = v.Name
			p.AccessFlags = cg.METHOD_PARAMETER_TYPE_ACC_MANDATED
			method.AttributeLucyReturnListNames.Parameters =
				append(method.AttributeLucyReturnListNames.Parameters, p)
		}
	}

	if t := makeClass.buildFunctionAutoVar(class, method.Code, f, context, state); t > method.Code.MaxStack {
		method.Code.MaxStack = t
	}
	makeClass.buildBlock(class, method.Code, &f.Block, context, state)
	return
}
func (makeClass *MakeClass) buildFunctionAutoVar(class *cg.ClassHighLevel, code *cg.AttributeCode,
	f *ast.Function, context *Context, state *StackMapState) (maxStack uint16) {
	if f.AutoVariableForException != nil {
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		f.AutoVariableForException.Offset = code.MaxLocals
		code.MaxLocals++
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, f.AutoVariableForException.Offset)...)
		state.appendLocals(class,
			state.newObjectVariableType(throwableClass))
		maxStack = 1
	}
	if f.AutoVariableForReturnBecauseOfDefer != nil {
		if len(f.Type.ReturnList) > 1 {
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			f.AutoVariableForReturnBecauseOfDefer.ForArrayList = code.MaxLocals
			code.MaxLocals++
			copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
				f.AutoVariableForReturnBecauseOfDefer.ForArrayList)...)
			state.appendLocals(class, state.newObjectVariableType(javaRootObjectArray))
		}
		maxStack = 1
	}
	if f.AutoVariableForMultiReturn != nil {
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		f.AutoVariableForMultiReturn.Offset = code.MaxLocals
		code.MaxLocals++
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, f.AutoVariableForMultiReturn.Offset)...)
		state.appendLocals(class, state.newObjectVariableType(javaRootObjectArray))
		maxStack = 1
	}

	return
}

func (makeClass *MakeClass) mkNonStaticFieldDefaultValue(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context, state *StackMapState) {
	for _, v := range context.class.Fields {
		if v.IsStatic() || v.Expression == nil {
			continue
		}
		code.Codes[code.CodeLength] = cg.OP_aload_0
		code.CodeLength++
		state.pushStack(class, state.newObjectVariableType(class.Name))
		stack, es := makeClass.makeExpression.build(class, code, v.Expression, context, state)
		if len(es) > 0 {
			state.pushStack(class, v.Expression.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
			fillOffsetForExits(es, code.CodeLength)
		}
		if t := 1 + stack; t > code.MaxStack {
			code.MaxStack = t
		}
		state.popStack(1)
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      class.Name,
			Field:      v.Name,
			Descriptor: JvmDescriptor.typeDescriptor(v.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
}
