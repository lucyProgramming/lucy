package jvm

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"os"
)

func (m *MakeClass) appendLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition, state *StackMapState) {
	if v.BeenCaptured { // capture
		t := &ast.VariableType{Typ: ast.VARIABLE_TYPE_OBJECT}
		t.Class = &ast.Class{}
		t.Class.Name = closure.getMeta(v.Typ.Typ).className
		state.Locals = append(state.Locals,
			state.newStackMapVerificationTypeInfo(class, t)...)
	} else {
		state.Locals = append(state.Locals, state.newStackMapVerificationTypeInfo(class, v.Typ)...)
	}
}
func (m *MakeClass) buildFunctionParameterAndReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode, f *ast.Function, context *Context, state *StackMapState) (maxstack uint16) {
	for _, v := range f.Typ.ParameterList { // insert into locals
		if v.BeenCaptured { // capture
			//because of stack map,capture parameter not allow
			fmt.Println(fmt.Sprintf("%s capture parameter not allow", ast.ErrMsgPrefix(v.Pos)))
			os.Exit(1)
		}
		if f.Name != ast.MAIN_FUNCTION_NAME {
			m.appendLocalVar(class, code, v, state)
		}
	}
	for _, v := range f.Typ.ReturnList {
		currentStack := uint16(0)
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
		backPatchEs(es, code.CodeLength)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if v.Typ.IsNumber() && v.Typ.Typ != v.Expression.VariableType.Typ {
			m.MakeExpression.numberTypeConverter(code, v.Expression.VariableType.Typ, v.Typ.Typ)
		}
		if t := currentStack + v.Typ.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		m.storeLocalVar(class, code, v)
		m.appendLocalVar(class, code, v, state)
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
	method.Code.MaxLocals = f.VarOffset // could  new slot when compile
	if method.IsConstruction {
		method.Code.Codes[method.Code.CodeLength] = cg.OP_aload_0
		method.Code.Codes[method.Code.CodeLength+1] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      class.SuperClass,
			Method:     special_method_init,
			Descriptor: "()V",
		}, method.Code.Codes[method.Code.CodeLength+2:method.Code.CodeLength+4])
		method.Code.CodeLength += 4
		method.Code.MaxStack = 1
	}
	state := &StackMapState{}
	// if function is main
	if f.Name == ast.MAIN_FUNCTION_NAME {
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
			t := &ast.VariableType{Typ: ast.VARIABLE_TYPE_JAVA_ARRAY}
			t.ArrayType = &ast.VariableType{Typ: ast.VARIABLE_TYPE_STRING}
			state.Locals = append(state.Locals,
				state.newStackMapVerificationTypeInfo(class, t)...)
			t = &ast.VariableType{Typ: ast.VARIABLE_TYPE_ARRAY}
			t.ArrayType = &ast.VariableType{Typ: ast.VARIABLE_TYPE_STRING}
			state.Locals = append(state.Locals,
				state.newStackMapVerificationTypeInfo(class, t)...)
		}

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
	for _, v := range f.AutoVars {
		switch v.(type) {
		case *ast.AutoVarForException:
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForException.Offset)...)
			maxstack = 1
			state.Locals = append(state.Locals,
				state.newStackMapVerificationTypeInfo(class,
					state.newObjectVariableType(java_throwable_class))...)

		case *ast.AutoVarForReturnBecauseOfDefer:
			code.Codes[code.CodeLength] = cg.OP_iconst_0
			code.CodeLength++
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT,
				f.AutoVarForReturnBecauseOfDefer.ExceptionIsNotNilWhenEnter)...)
			state.Locals = append(state.Locals,
				state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
			if len(f.Typ.ReturnList) > 1 {
				code.Codes[code.CodeLength] = cg.OP_iconst_0
				code.CodeLength++
				copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT,
					f.AutoVarForReturnBecauseOfDefer.IfReachBotton)...)
				state.Locals = append(state.Locals,
					state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
			}
			maxstack = 1
		case *ast.AutoVarForMultiReturn:
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForMultiReturn.Offset)...)
			maxstack = 1
			state.Locals = append(state.Locals,
				state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(java_arrylist_class))...)
		}
	}
	return
}

func (m *MakeClass) loadLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		return closure.loadLocalCloureVar(class, code, v)
	}
	maxstack = v.Typ.JvmSlotSize()
	switch v.Typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_iload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_iload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_iload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_iload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_iload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_LONG:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_lload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_lload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_lload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_lload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_lload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_fload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_fload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_fload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_fload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_fload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_dload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_dload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_dload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_dload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_dload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY: //[]int
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_aload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_aload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_aload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_aload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	}
	return
}

func (m *MakeClass) storeLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		closure.storeLocalCloureVar(class, code, v)
		return
	}
	maxstack = v.Typ.JvmSlotSize()
	switch v.Typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_istore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_istore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_istore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_istore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_istore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_LONG:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_lstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_lstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_lstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_lstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_lstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_fstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_fstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_fstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_fstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_fstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_dstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_dstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_dstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_dstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_dstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY: //[]int
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_astore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_astore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_astore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_astore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_astore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	}
	return
}
