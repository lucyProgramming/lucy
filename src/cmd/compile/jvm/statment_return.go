package jvm

import (
	"encoding/binary"
	//"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildReturnStatement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	statementReturn *ast.StatementReturn, context *Context, state *StackMapState) (maxstack uint16) {
	if context.function.NoReturnValue() { // no return value
		if statementReturn.Defers != nil && len(statementReturn.Defers) > 0 {
			code.Codes[code.CodeLength] = cg.OP_aconst_null // expect exception on stack
			code.CodeLength++
			if 1 > maxstack {
				maxstack = 1
			}
			stack := m.buildDefersForReturn(class, code, context, state, statementReturn)
			if stack > maxstack {
				maxstack = stack
			}
		}
		code.Codes[code.CodeLength] = cg.OP_return
		code.CodeLength++
		return
	}
	if len(context.function.Typ.ReturnList) == 1 {
		var es []*cg.JumpBackPatch
		if len(statementReturn.Expressions) > 0 {
			maxstack, es = m.MakeExpression.build(class, code, statementReturn.Expressions[0], context, state)
			if len(es) > 0 {
				backPatchEs(es, code.CodeLength)
				state.pushStack(class, statementReturn.Expressions[0].Value)
				context.MakeStackMap(code, state, code.CodeLength)
				state.popStack(1)
			}
			if statementReturn.Expressions[0].Value.IsNumber() &&
				statementReturn.Expressions[0].Value.Typ != context.function.Typ.ReturnList[0].Typ.Typ {
				m.MakeExpression.numberTypeConverter(code,
					statementReturn.Expressions[0].Value.Typ, context.function.Typ.ReturnList[0].Typ.Typ)
			}
		} else { // load return parameter
		}
		// execute defer first
		if len(statementReturn.Defers) > 0 {
			//return value  is on stack,  store it temp var
			if len(statementReturn.Expressions) > 0 { //rewrite return value
				m.storeLocalVar(class, code, context.function.Typ.ReturnList[0])
			}
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			if 1 > maxstack {
				maxstack = 1
			}
			stack := m.buildDefersForReturn(class, code, context, state, statementReturn)
			if stack > maxstack {
				maxstack = stack
			}
			//restore the stack
			if len(statementReturn.Expressions) > 0 { //restore stack
				m.loadLocalVar(class, code, context.function.Typ.ReturnList[0])
			}
		}
		// in this case,load local var is not under exception handle,should be ok
		if len(statementReturn.Expressions) == 0 {
			m.loadLocalVar(class, code, context.function.Typ.ReturnList[0])
		}
		switch context.function.Typ.ReturnList[0].Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_ENUM:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_ireturn
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lreturn
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_freturn
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dreturn
		case ast.VARIABLE_TYPE_STRING:
			fallthrough
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_MAP:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY:
			code.Codes[code.CodeLength] = cg.OP_areturn
		}
		code.CodeLength++
		return
	}

	//multi returns
	if len(statementReturn.Expressions) > 0 {
		//new a array list
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst(java_arrylist_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup // dup on stack
		code.CodeLength += 4
		//call init
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_arrylist_class,
			Method:     special_method_init,
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxstack = 2 // max stack is 2
		{
			t := state.newStackMapVerificationTypeInfo(class,
				state.newObjectVariableType(java_arrylist_class))
			state.Stacks = append(state.Stacks, t)
			state.Stacks = append(state.Stacks, t)
			defer func() {
				state.popStack(2)
			}()
		}
		for _, v := range statementReturn.Expressions {
			currentStack := uint16(1)
			code.Codes[code.CodeLength] = cg.OP_dup // dup array list
			code.CodeLength++
			currentStack++
			if currentStack > maxstack {
				maxstack = maxstack
			}
			if v.MayHaveMultiValue() && len(v.Values) > 1 {
				if currentStack > maxstack {
					maxstack = maxstack
				} // make the call
				stack, _ := m.MakeExpression.build(class, code, v, context, nil)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				code.Codes[code.CodeLength] = cg.OP_invokevirtual
				class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
					Class:      java_arrylist_class,
					Method:     "addAll",
					Descriptor: "(Ljava/util/Collection;)Z",
				}, code.Codes[code.CodeLength+1:code.CodeLength+3])
				code.Codes[code.CodeLength+3] = cg.OP_pop
				code.CodeLength += 4
				continue
			}
			variableType := v.Value
			if v.MayHaveMultiValue() {
				variableType = v.Values[0]
			}
			stack, es := m.MakeExpression.build(class, code, v, context, state)
			if len(es) > 0 {
				backPatchEs(es, code.CodeLength)
				state.pushStack(class, v.Value)
				context.MakeStackMap(code, state, code.CodeLength)
				state.popStack(1) // must be bool expression
			}
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			//convert to object
			primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
			// append
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      java_arrylist_class,
				Method:     "add",
				Descriptor: "(Ljava/lang/Object;)Z",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = cg.OP_pop
			code.CodeLength += 4
		}
	} else {
		//nothing to do
	}
	if statementReturn.Defers != nil && len(statementReturn.Defers) > 0 {
		//store a simple var,should be no exception
		if len(statementReturn.Expressions) > 0 {
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT,
				context.function.AutoVarForReturnBecauseOfDefer.ForArrayList)...)

		}
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		if 1 > maxstack {
			maxstack = 1
		}
		stack := m.buildDefersForReturn(class, code, context, state, statementReturn)
		if stack > maxstack {
			maxstack = stack
		}
		//restore the stack
		copyOP(code,
			loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT,
				context.function.AutoVarForReturnBecauseOfDefer.ForArrayList)...)
	}
	if len(statementReturn.Expressions) > 0 {
		code.Codes[code.CodeLength] = cg.OP_areturn
		code.CodeLength++
		return
	}
	stack := m.buildReturnFromFunctionReturnList(class, code, context)
	if stack > maxstack {
		maxstack = stack
	}
	return
}

func (m *MakeClass) buildReturnFromFunctionReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context) (maxstack uint16) {
	if context.function.NoReturnValue() { // when has no return,should not call this function
		return
	}
	if len(context.function.Typ.ReturnList) == 1 {
		m.loadLocalVar(class, code, context.function.Typ.ReturnList[0])
		maxstack = jvmSize(context.function.Typ.ReturnList[0].Typ)
		switch context.function.Typ.ReturnList[0].Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_ENUM:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_ireturn
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lreturn
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_freturn
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dreturn
		case ast.VARIABLE_TYPE_STRING:
			fallthrough
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_MAP:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY:
			code.Codes[code.CodeLength] = cg.OP_areturn
		}
		code.CodeLength++
		return
	}
	//multi returns
	//new a array list
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(java_arrylist_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup // dup on stack
	code.CodeLength += 4
	//call init
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 2 // max stack is
	for _, v := range context.function.Typ.ReturnList {
		currentStack := uint16(1)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack++
		if currentStack > maxstack {
			maxstack = currentStack
		}
		m.loadLocalVar(class, code, v)
		if t := currentStack + jvmSize(v.Typ); t > maxstack {
			maxstack = t
		}
		if v.Typ.IsPointer() == false {
			primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, v.Typ)
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_arrylist_class,
			Method:     "add",
			Descriptor: "(Ljava/lang/Object;)Z",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_pop
		code.CodeLength += 4
	}
	code.Codes[code.CodeLength] = cg.OP_areturn
	code.CodeLength++
	return
}

func (m *MakeClass) buildDefersForReturn(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context, ss *StackMapState,
	statementReturn *ast.StatementReturn) (maxstack uint16) {
	if len(statementReturn.Defers) == 0 {
		return
	}
	index := len(statementReturn.Defers) - 1
	for index >= 0 { // build defer,cannot have return statement is defer
		state := statementReturn.Defers[index].StackMapState.(*StackMapState)
		state.addTop(ss)
		state.pushStack(class, state.newObjectVariableType(java_throwable_class))
		context.MakeStackMap(code, state, code.CodeLength)
		e := &cg.ExceptionTable{}
		e.StartPc = uint16(statementReturn.Defers[index].StartPc)
		e.Endpc = uint16(code.CodeLength)
		e.HandlerPc = uint16(code.CodeLength)
		if statementReturn.Defers[index].ExceptionClass == nil {
			e.CatchType = class.Class.InsertClassConst(ast.DEFAULT_EXCEPTION_CLASS)
		} else {
			e.CatchType = class.Class.InsertClassConst(statementReturn.Defers[index].ExceptionClass.Name) // custom class
		}
		code.Exceptions = append(code.Exceptions, e)
		//expect exception on stack
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT,
			context.function.AutoVarForException.Offset)...) // this code will make stack is empty
		state.popStack(1)
		// build block
		context.Defer = statementReturn.Defers[index]
		m.buildBlock(class, code, &statementReturn.Defers[index].Block, context, state)
		context.Defer = nil
		//if need throw
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		state.pushStack(class, state.newObjectVariableType(java_throwable_class))
		context.MakeStackMap(code, state, code.CodeLength+6)
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.popStack(1)
		code.Codes[code.CodeLength] = cg.OP_ifnonnull
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
		code.Codes[code.CodeLength+3] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 4) // goto pop
		code.Codes[code.CodeLength+6] = cg.OP_athrow
		code.Codes[code.CodeLength+7] = cg.OP_pop // pop exception on stack
		code.CodeLength += 8
		if index > 0 { // not the last defer, no need to look at
			index--
			continue
		}
		//expection that have been handled
		if len(statementReturn.Expressions) > 0 {
			if len(context.function.Typ.ReturnList) == 1 {
				t := m.buildReturnFromFunctionReturnList(class, code, context)
				if t > code.MaxStack {
					code.MaxStack = t
				}
			} else {
				//load when function have multi returns if read to end
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForReturnBecauseOfDefer.ForArrayList)...)
				code.Codes[code.CodeLength] = cg.OP_ifnull
				binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
				code.Codes[code.CodeLength+3] = cg.OP_goto
				length := code.CodeLength + 3
				code.CodeLength += 6
				context.MakeStackMap(code, state, code.CodeLength)
				m.buildReturnFromFunctionReturnList(class, code, context)
				context.MakeStackMap(code, state, code.CodeLength)
				binary.BigEndian.PutUint16(code.Codes[length+1:length+3], uint16(code.CodeLength-length))
			}
		}
		index--
	}
	return
}
