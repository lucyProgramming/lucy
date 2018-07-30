package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildReturnStatement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	statementReturn *ast.StatementReturn, context *Context, state *StackMapState) (maxStack uint16) {
	if context.function.NoReturnValue() { // no return value
		if statementReturn.Defers != nil && len(statementReturn.Defers) > 0 {
			code.Codes[code.CodeLength] = cg.OP_aconst_null // expect exception on stack
			code.CodeLength++
			if 1 > maxStack {
				maxStack = 1
			}
			stack := buildPackage.buildDefersForReturn(class, code, context, state, statementReturn)
			if stack > maxStack {
				maxStack = stack
			}
		}
		code.Codes[code.CodeLength] = cg.OP_return
		code.CodeLength++
		return
	}
	if len(context.function.Type.ReturnList) == 1 {
		if len(statementReturn.Expressions) > 0 {
			maxStack = buildPackage.BuildExpression.build(class, code, statementReturn.Expressions[0], context, state)
		}
		// execute defer first
		if len(statementReturn.Defers) > 0 {
			//return value  is on stack,  store to local var
			if len(statementReturn.Expressions) > 0 { //rewrite return value
				buildPackage.storeLocalVar(class, code, context.function.Type.ReturnList[0])
			}
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			if 1 > maxStack {
				maxStack = 1
			}
			stack := buildPackage.buildDefersForReturn(class, code, context, state, statementReturn)
			if stack > maxStack {
				maxStack = stack
			}
			//restore the stack
			if len(statementReturn.Expressions) > 0 { //restore stack
				buildPackage.loadLocalVar(class, code, context.function.Type.ReturnList[0])
			}
		}
		// in this case,load local var is not under exception handle,should be ok
		if len(statementReturn.Expressions) == 0 {
			buildPackage.loadLocalVar(class, code, context.function.Type.ReturnList[0])
		}
		switch context.function.Type.ReturnList[0].Type.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_ireturn
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_lreturn
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_freturn
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_dreturn
		case ast.VariableTypeFunction:
			fallthrough
		case ast.VariableTypeJavaArray:
			fallthrough
		case ast.VariableTypeString:
			fallthrough
		case ast.VariableTypeObject:
			fallthrough
		case ast.VariableTypeMap:
			fallthrough
		case ast.VariableTypeArray:
			code.Codes[code.CodeLength] = cg.OP_areturn
		}
		code.CodeLength++
		return
	}
	//multi returns
	if len(statementReturn.Expressions) > 0 {
		if len(statementReturn.Expressions) == 1 {
			maxStack = buildPackage.BuildExpression.build(class, code, statementReturn.Expressions[0], context, state)
		} else {
			loadInt32(class, code, int32(len(context.function.Type.ReturnList)))
			code.Codes[code.CodeLength] = cg.OP_anewarray
			class.InsertClassConst(javaRootClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			maxStack = 2 // max stack is 2
			arrayListObject := state.newObjectVariableType(javaRootObjectArray)
			state.pushStack(class, arrayListObject)
			state.pushStack(class, arrayListObject)
			defer state.popStack(2)
			index := int32(0)
			for _, v := range statementReturn.Expressions {
				currentStack := uint16(1)
				if v.HaveMultiValue() {
					stack := buildPackage.BuildExpression.build(class, code, v, context, state)
					if t := currentStack + stack; t > maxStack {
						maxStack = t
					}
					autoVar := storeMultiValueAutoVar(class, code, state)
					for kk, _ := range v.MultiValues {
						currentStack := uint16(1)
						code.Codes[code.CodeLength] = cg.OP_dup // dup array list
						code.CodeLength++
						currentStack++
						stack = autoVar.unPackObject(class, code, kk)
						if t := stack + currentStack; t > maxStack {
							maxStack = t
						}
						loadInt32(class, code, index)
						if t := currentStack + 2; t > maxStack {
							maxStack = t
						}
						code.Codes[code.CodeLength] = cg.OP_swap
						code.Codes[code.CodeLength+1] = cg.OP_aastore
						code.CodeLength += 2
						index++
					}
					continue
				}
				code.Codes[code.CodeLength] = cg.OP_dup // dup array list
				code.CodeLength++
				currentStack++
				stack := buildPackage.BuildExpression.build(class, code, v, context, state)
				if t := stack + currentStack; t > maxStack {
					maxStack = t
				}
				//convert to object
				if v.Value.IsPointer() == false {
					typeConverter.packPrimitives(class, code, v.Value)
				}
				// append
				loadInt32(class, code, index)

				if t := currentStack + 2; t > maxStack {
					maxStack = t
				}
				code.Codes[code.CodeLength] = cg.OP_swap
				code.Codes[code.CodeLength+1] = cg.OP_aastore
				code.CodeLength += 2
				index++
			}
		}
	}
	if statementReturn.Defers != nil && len(statementReturn.Defers) > 0 {
		//store a simple var,should be no exception
		if len(statementReturn.Expressions) > 0 {
			copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
				context.function.AutoVariableForReturnBecauseOfDefer.Offset)...)
		}
		code.Codes[code.CodeLength] = cg.OP_aconst_null
		code.CodeLength++
		if 1 > maxStack {
			maxStack = 1
		}
		stack := buildPackage.buildDefersForReturn(class, code, context, state, statementReturn)
		if stack > maxStack {
			maxStack = stack
		}
		//restore the stack
		if len(statementReturn.Expressions) > 0 {
			copyOPs(code,
				loadLocalVariableOps(ast.VariableTypeObject,
					context.function.AutoVariableForReturnBecauseOfDefer.Offset)...)
		}
	}
	if len(statementReturn.Expressions) > 0 {
		code.Codes[code.CodeLength] = cg.OP_areturn
		code.CodeLength++
		return
	}
	stack := buildPackage.buildReturnFromFunctionReturnList(class, code, context)
	if stack > maxStack {
		maxStack = stack
	}
	return
}

func (buildPackage *BuildPackage) buildReturnFromFunctionReturnList(class *cg.ClassHighLevel,
	code *cg.AttributeCode, context *Context) (maxStack uint16) {
	if context.function.NoReturnValue() { // when has no return,should not call this function
		return
	}
	if len(context.function.Type.ReturnList) == 1 {
		buildPackage.loadLocalVar(class, code, context.function.Type.ReturnList[0])
		maxStack = jvmSlotSize(context.function.Type.ReturnList[0].Type)
		switch context.function.Type.ReturnList[0].Type.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_ireturn
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_lreturn
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_freturn
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_dreturn
		case ast.VariableTypeFunction:
			fallthrough
		case ast.VariableTypeJavaArray:
			fallthrough
		case ast.VariableTypeString:
			fallthrough
		case ast.VariableTypeObject:
			fallthrough
		case ast.VariableTypeMap:
			fallthrough
		case ast.VariableTypeArray:
			code.Codes[code.CodeLength] = cg.OP_areturn
		}
		code.CodeLength++
		return
	}
	//multi returns
	//new a array list
	loadInt32(class, code, int32(len(context.function.Type.ReturnList)))
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst(javaRootClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxStack = 1 // max stack is
	index := int32(0)
	for _, v := range context.function.Type.ReturnList {
		currentStack := uint16(1)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		currentStack++
		buildPackage.loadLocalVar(class, code, v)
		if t := currentStack + jvmSlotSize(v.Type); t > maxStack {
			maxStack = t
		}
		if v.Type.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.Type)
		}
		loadInt32(class, code, index)
		if 4 > maxStack {
			maxStack = 4
		}
		code.Codes[code.CodeLength] = cg.OP_swap
		code.Codes[code.CodeLength+1] = cg.OP_aastore
		code.CodeLength += 2
		index++
	}
	code.Codes[code.CodeLength] = cg.OP_areturn
	code.CodeLength++
	return
}

func (buildPackage *BuildPackage) buildDefersForReturn(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context, from *StackMapState,
	statementReturn *ast.StatementReturn) (maxStack uint16) {
	if len(statementReturn.Defers) == 0 {
		return
	}
	index := len(statementReturn.Defers) - 1
	for index >= 0 { // build defer,cannot have return statement is defer
		state := statementReturn.Defers[index].StackMapState.(*StackMapState)
		state = (&StackMapState{}).FromLast(state) // clone
		state.addTop(from)
		state.pushStack(class, state.newObjectVariableType(throwableClass))
		context.MakeStackMap(code, state, code.CodeLength)
		e := &cg.ExceptionTable{}
		e.StartPc = uint16(statementReturn.Defers[index].StartPc)
		e.EndPc = uint16(code.CodeLength)
		e.HandlerPc = uint16(code.CodeLength)
		if statementReturn.Defers[index].ExceptionClass == nil {
			e.CatchType = class.Class.InsertClassConst(ast.DefaultExceptionClass)
		} else {
			e.CatchType = class.Class.InsertClassConst(statementReturn.Defers[index].ExceptionClass.Name) // custom class
		}
		code.Exceptions = append(code.Exceptions, e)
		//expect exception on stack
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
			context.function.AutoVariableForException.Offset)...) // this code will make stack is empty
		state.popStack(1)
		// build block
		context.Defer = statementReturn.Defers[index]
		buildPackage.buildBlock(class, code, &statementReturn.Defers[index].Block, context, state)
		from.addTop(state)
		context.Defer = nil
		for _, v := range statementReturn.Defers[index].Labels {
			v.Reset()
		}
		if index > 0 { // not the last defer, no need to look at
			index--
			continue
		}
		//if need throw
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, context.function.AutoVariableForException.Offset)...)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		state.pushStack(class, state.newObjectVariableType(throwableClass))
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

		//exception that have been handled
		if len(statementReturn.Expressions) > 0 && len(context.function.Type.ReturnList) > 1 {
			//load when function have multi returns if read to end
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, context.function.AutoVariableForReturnBecauseOfDefer.Offset)...)
			code.Codes[code.CodeLength] = cg.OP_ifnull
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
			code.Codes[code.CodeLength+3] = cg.OP_goto
			length := code.CodeLength + 3
			code.CodeLength += 6
			context.MakeStackMap(code, state, code.CodeLength)
			buildPackage.buildReturnFromFunctionReturnList(class, code, context)
			context.MakeStackMap(code, state, code.CodeLength)
			binary.BigEndian.PutUint16(code.Codes[length+1:length+3], uint16(code.CodeLength-length))
		}
		index--
	}
	return
}
