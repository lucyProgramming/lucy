package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) appendLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode,
	v *ast.VariableDefinition, state *StackMapState) {
	if v.BeenCaptured { // capture
		t := &ast.VariableType{Typ: ast.VARIABLE_TYPE_OBJECT}
		t.Class = &ast.Class{}
		t.Class.Name = closure.getMeta(v.Typ.Typ).className
		v.LocalValOffset = state.appendLocals(class, code, t)
	} else {
		v.LocalValOffset = state.appendLocals(class, code, v.Typ)
	}
}
func (m *MakeClass) buildFunctionParameterAndReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode, f *ast.Function, context *Context, state *StackMapState) (maxstack uint16) {
	for _, v := range f.Typ.ParameterList { // insert into locals
		if v.BeenCaptured { // capture
			//because of stack map,capture parameter not allow
			panic("...")
		}
		if f.Name != ast.MAIN_FUNCTION_NAME {
			m.appendLocalVar(class, code, v, state)
		}
	}
	for _, v := range f.Typ.ReturnList {
		currentStack := uint16(0)
		m.appendLocalVar(class, code, v, state)
		if v.BeenCaptured { //create closure object
			stack := closure.createCloureVar(class, code, v)
			if stack > maxstack {
				maxstack = stack
			}
			// then load
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
			currentStack = 1
		}
		stack, es := m.MakeExpression.build(class, code, v.Expression, context, state)
		if len(es) > 0 {
			length := len(state.Stacks)
			if v.BeenCaptured {
				state.Stacks = append(state.Stacks,
					state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className))...)
			}
			state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, v.Typ)...)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(len(state.Stacks) - length)
			backPatchEs(es, code.CodeLength)
		}
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if v.Typ.IsNumber() && v.Typ.Typ != v.Expression.Value.Typ {
			m.MakeExpression.numberTypeConverter(code, v.Expression.Value.Typ, v.Typ.Typ)
		}
		if t := currentStack + jvmSize(v.Typ); t > maxstack {
			maxstack = t
		}
		m.storeLocalVar(class, code, v)
	}
	return
}

func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, method *cg.MethodHighLevel, f *ast.Function) {
	context := &Context{}
	context.function = f
	context.method = method
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
	}()
	state := &StackMapState{}
	if method.IsConstruction { // construction method
		method.Code.Codes[method.Code.CodeLength] = cg.OP_aload_0
		method.Code.Codes[method.Code.CodeLength+1] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      class.SuperClass,
			Method:     special_method_init,
			Descriptor: "()V",
		}, method.Code.Codes[method.Code.CodeLength+2:method.Code.CodeLength+4])
		method.Code.CodeLength += 4
		method.Code.MaxStack = 1
		method.Code.MaxLocals = 1
		state.appendLocals(class, method.Code, state.newObjectVariableType(class.Name))
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
			state.appendLocals(class, code, t)
			// []string lucy style
			t = &ast.VariableType{Typ: ast.VARIABLE_TYPE_ARRAY}
			t.ArrayType = &ast.VariableType{Typ: ast.VARIABLE_TYPE_STRING}
			state.appendLocals(class, code, t)
		}
		f.Typ.ParameterList[0].LocalValOffset = 1 // main(args []string) args offset is 1
	} else if method.AccessFlags&cg.ACC_METHOD_STATIC == 0 { // instance method
		state.appendLocals(class, method.Code,
			state.newObjectVariableType(class.Name))
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
		f.AutoVarForException.Offset = state.appendLocals(class, code,
			state.newObjectVariableType(java_throwable_class))
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForException.Offset)...)
		maxstack = 1
	}
	if f.AutoVarForReturnBecauseOfDefer != nil {
		//code.Codes[code.CodeLength] = cg.OP_iconst_0
		//code.CodeLength++
		//f.AutoVarForReturnBecauseOfDefer.ExceptionIsNotNilWhenEnter =
		//	state.appendLocals(class, code, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
		//copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT,
		//	f.AutoVarForReturnBecauseOfDefer.ExceptionIsNotNilWhenEnter)...)
		if len(f.Typ.ReturnList) > 1 {
			//if reach botton
			code.Codes[code.CodeLength] = cg.OP_iconst_0
			code.CodeLength++
			f.AutoVarForReturnBecauseOfDefer.IfReachBotton =
				state.appendLocals(class, code, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT,
				f.AutoVarForReturnBecauseOfDefer.IfReachBotton)...)

			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			f.AutoVarForReturnBecauseOfDefer.ForArrayList =
				state.appendLocals(class, code, state.newObjectVariableType(java_arrylist_class))
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT,
				f.AutoVarForReturnBecauseOfDefer.ForArrayList)...)
		}
		maxstack = 1
	}
	if f.AutoVarForMultiReturn != nil {
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		f.AutoVarForMultiReturn.Offset =
			state.appendLocals(class, code, state.newObjectVariableType(java_arrylist_class))
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForMultiReturn.Offset)...)
		maxstack = 1
	}

	return
}
