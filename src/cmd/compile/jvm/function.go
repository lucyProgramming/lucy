package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildFunctionParameterAndReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode, f *ast.Function, context *Context, state *StackMapState) (maxstack uint16) {
	for _, v := range f.Typ.ParameterList { // insert into locals
		v.LocalValOffset = code.MaxLocals
		code.MaxLocals += jvmSize(v.Typ)
		state.appendLocals(class, v.Typ)
	}
	for _, v := range f.Typ.ParameterList {
		if v.BeenCaptured == false { // capture
			continue
		}
		stack := closure.createCloureVar(class, code, v.Typ)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		copyOP(code, loadSimpleVarOp(v.Typ.Typ, v.LocalValOffset)...)
		if t := 2 + jvmSize(v.Typ); t > maxstack {
			maxstack = t
		}
		m.storeLocalVar(class, code, v)
		v.LocalValOffset = code.MaxLocals
		code.MaxLocals++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
		state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className))
	}
	for _, v := range f.Typ.ReturnList {
		currentStack := uint16(0)
		if v.BeenCaptured { //create closure object
			v.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			stack := closure.createCloureVar(class, code, v.Typ)
			if stack > maxstack {
				maxstack = stack
			}
			// then load
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
			if 2 > maxstack {
				maxstack = 2
			}
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
			currentStack = 1
		} else {
			v.LocalValOffset = code.MaxLocals
			code.MaxLocals += jvmSize(v.Typ)
		}
		stack, es := m.MakeExpression.build(class, code, v.Expression, context, state)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			length := len(state.Stacks)
			if v.BeenCaptured {
				state.pushStack(class,
					state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className))
			}
			state.pushStack(class, v.Typ)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(len(state.Stacks) - length)
		}
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		m.storeLocalVar(class, code, v)
		if v.BeenCaptured {
			state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className))
		} else {
			state.appendLocals(class, v.Typ)
		}
	}
	return
}

func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, astClass *ast.Class, method *cg.MethodHighLevel, f *ast.Function) {
	context := &Context{}
	context.class = astClass
	context.function = f
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
	}()
	state := &StackMapState{}
	if method.IsConstruction { // construction method
		if f.ConstructionMethodCalledByUser == false {
			method.Code.Codes[method.Code.CodeLength] = cg.OP_aload_0
			method.Code.Codes[method.Code.CodeLength+1] = cg.OP_invokespecial
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      class.SuperClass,
				Method:     special_method_init,
				Descriptor: "()V",
			}, method.Code.Codes[method.Code.CodeLength+2:method.Code.CodeLength+4])
			method.Code.CodeLength += 4
			method.Code.MaxStack = 1
			// field default value
			m.mkFieldDefaultValue(class, method.Code, context, state)
		}
		method.Code.MaxLocals = 1

		v := &cg.StackMap_UninitializedThis_variable_info{}
		state.Locals = append(state.Locals, &cg.StackMap_verification_type_info{Verify: v})
	} else if f.Name == ast.MAIN_FUNCTION_NAME { // main function
		code := method.Code
		code.Codes[code.CodeLength] = cg.OP_new
		meta := ArrayMetas[ast.VARIABLE_TYPE_STRING]
		class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_STRING, 0)...)
		if 3 > code.MaxStack {
			code.MaxStack = 3
		}
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.classname,
			Method:     special_method_init,
			Descriptor: meta.constructorFuncDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, 1)...)
		{
			// String[] java style
			t := &ast.VariableType{Typ: ast.VARIABLE_TYPE_JAVA_ARRAY}
			t.ArrayType = &ast.VariableType{Typ: ast.VARIABLE_TYPE_STRING}
			state.appendLocals(class, t)
		}
		method.Code.MaxLocals = 1
	} else if method.AccessFlags&cg.ACC_METHOD_STATIC == 0 { // instance method
		method.Code.MaxLocals = 1
		state.appendLocals(class, state.newObjectVariableType(class.Name))
	}
	if LucyMethodSignatureParser.Need(&f.Typ) {
		d := &cg.AttributeLucyMethodDescriptor{}
		d.Descriptor = LucyMethodSignatureParser.Encode(f)
		method.AttributeLucyMethodDescritor = d
	}
	if f.HaveDefaultValue {
		method.AttributeDefaultParameters = FunctionDefaultValueParser.Encode(class, f)
	}
	if t := m.buildFunctionParameterAndReturnList(class, method.Code, f, context, state); t > method.Code.MaxStack {
		method.Code.MaxStack = t
	}
	if t := m.buildFunctionAutoVar(class, method.Code, f, context, state); t > method.Code.MaxStack {
		method.Code.MaxStack = t
	}
	m.buildBlock(class, method.Code, f.Block, context, state)
	return
}
func (m *MakeClass) buildFunctionAutoVar(class *cg.ClassHighLevel, code *cg.AttributeCode, f *ast.Function, context *Context, state *StackMapState) (maxstack uint16) {
	if f.AutoVarForException != nil {
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		f.AutoVarForException.Offset = code.MaxLocals
		code.MaxLocals++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForException.Offset)...)
		state.appendLocals(class,
			state.newObjectVariableType(java_throwable_class))
		maxstack = 1
	}
	if f.AutoVarForReturnBecauseOfDefer != nil {
		if len(f.Typ.ReturnList) > 1 {
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			f.AutoVarForReturnBecauseOfDefer.ForArrayList = code.MaxLocals
			code.MaxLocals++
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT,
				f.AutoVarForReturnBecauseOfDefer.ForArrayList)...)
			state.appendLocals(class, state.newObjectVariableType(java_arrylist_class))
		}
		maxstack = 1
	}
	if f.AutoVarForMultiReturn != nil {
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		f.AutoVarForMultiReturn.Offset = code.MaxLocals
		code.MaxLocals++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForMultiReturn.Offset)...)
		state.appendLocals(class, state.newObjectVariableType(java_arrylist_class))
		maxstack = 1
	}

	return
}

func (m *MakeClass) mkFieldDefaultValue(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context, state *StackMapState) {
	for _, v := range context.class.Fields {
		if v.DefaultValue == nil || v.IsStatic() {
			continue
		}
		code.Codes[code.CodeLength] = cg.OP_aload_0
		code.CodeLength++
		stack, _ := m.MakeExpression.build(class, code, v.Expression, context, state)
		if t := 1 + stack; t > code.MaxStack {
			code.MaxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_putfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      class.Name,
			Field:      v.Name,
			Descriptor: Descriptor.typeDescriptor(v.Typ),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
}
