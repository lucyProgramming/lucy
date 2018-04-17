package jvm

import (
	"encoding/binary"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context, state *StackMapState) {
	if len(b.Defers) > 0 { // should be more defers when compile
		context.Defers = append(context.Defers, b.Defers...)
	}
	for _, s := range b.Statements {
		maxstack := m.buildStatement(class, code, b, s, context, state)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		if len(state.Stacks) > 0 {
			panic(fmt.Sprintf("stack is not empty:%d", len(state.Stacks)))
		}
	}
	if len(b.Defers) > 0 {
		//slice out
		context.Defers = context.Defers[0 : len(context.Defers)-len(b.Defers)]
		//execute form begin to  end , no exceptions,push a null exception for defer to catch
		if b.IsFunctionTopBlock == false {
			code.Codes[code.CodeLength] = cg.OP_aconst_null // let defer to catch
			code.CodeLength++
		}
	}
	if b.IsFunctionTopBlock == false {
		m.buildDefers(class, code, state, context, b.Defers, true, nil)
	}
	return
}

func (m *MakeClass) buildDefers(class *cg.ClassHighLevel, code *cg.AttributeCode, state *StackMapState, context *Context, ds []*ast.Defer, needExceptionTable bool, r *ast.StatementReturn) {
	if len(ds) == 0 {
		return
	}
	var endPc, startPc int
	if needExceptionTable {
		endPc = code.CodeLength
		startPc = ds[len(ds)-1].StartPc
	}
	index := len(ds) - 1
	for index >= 0 { // build defer,cannot have return statement is defer
		// insert exceptions
		if needExceptionTable {
			{
				t := &ast.VariableType{}
				t.Typ = ast.VARIABLE_TYPE_OBJECT
				t.Class = &ast.Class{}
				t.Class.Name = "java/lang/Throwable"
				state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, t)...)
			}
			code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(state, code.CodeLength))
			state.popStack(1)
			e := &cg.ExceptionTable{}
			e.StartPc = uint16(startPc)
			e.Endpc = uint16(endPc)
			e.HandlerPc = uint16(code.CodeLength)
			e.CatchType = class.Class.InsertClassConst(java_throwable_class) //runtime
			code.Exceptions = append(code.Exceptions, e)
			startPc = code.CodeLength
			if index == len(ds)-1 && r != nil && context.function.NoReturnValue() == false {
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				if 2 > code.MaxStack {
					code.MaxStack = 2
				}
				code.Codes[code.CodeLength] = cg.OP_ifnonnull
				binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
				code.Codes[code.CodeLength+3] = cg.OP_goto
				op := storeSimpleVarOp(ast.VARIABLE_TYPE_INT, context.function.AutoVarForReturnBecauseOfDefer.ExceptionIsNotNilWhenEnter)
				binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 4+uint16(len(op)))
				code.Codes[code.CodeLength+6] = cg.OP_iconst_1
				code.CodeLength += 7
				copyOP(code, op...)
			}
			//expect exception on stack
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...) // this code will make stack is empty
		}
		m.buildBlock(class, code, &ds[index].Block, context, (&StackMapState{}).FromLast(state))
		// load to stack
		// this code maxStack is 2
		if 2 > code.MaxStack {
			code.MaxStack = 2
		}
		if needExceptionTable == false {
			index--
			continue
		}
		if index != 0 {
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...)
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
			code.Codes[code.CodeLength] = cg.OP_ifnonnull
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
			code.Codes[code.CodeLength+3] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 4)
			code.Codes[code.CodeLength+6] = cg.OP_athrow
			code.CodeLength += 7
			endPc = code.CodeLength
			index--
			continue
		}
		// index is 0,last defer block
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_ifnonnull
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
		code.Codes[code.CodeLength+3] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 4) // goto pop
		code.Codes[code.CodeLength+6] = cg.OP_athrow
		code.Codes[code.CodeLength+7] = cg.OP_pop // pop exception on stack
		code.CodeLength += 8
		if r == nil || context.function.NoReturnValue() || len(r.Expressions) == 0 {
			endPc = code.CodeLength
			index--
			continue
		}
		// load if enter defers there is a exception
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT,
			context.function.AutoVarForReturnBecauseOfDefer.ExceptionIsNotNilWhenEnter)...)
		code.Codes[code.CodeLength] = cg.OP_ifne
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
		code.Codes[code.CodeLength+3] = cg.OP_goto
		noExceptionExitCodeLength := code.CodeLength + 3
		code.CodeLength += 6
		//expection that have been handled
		if len(context.function.Typ.ReturnList) == 1 {
			t := m.buildReturnFromFunctionReturnList(class, code, context)
			if t > code.MaxStack {
				code.MaxStack = t
			}
		} else {
			//load when function have multi returns if read to end
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, context.function.AutoVarForReturnBecauseOfDefer.IfReachBotton)...)
			code.Codes[code.CodeLength] = cg.OP_ifeq
			codeLength := code.CodeLength
			code.CodeLength += 3
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForMultiReturn.Offset)...)
			code.Codes[code.CodeLength] = cg.OP_areturn
			code.CodeLength++
			binary.BigEndian.PutUint16(code.Codes[codeLength+1:codeLength+3], uint16(code.CodeLength-codeLength))
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 0)
		}
		binary.BigEndian.PutUint16(code.Codes[noExceptionExitCodeLength+1:noExceptionExitCodeLength+3], uint16(code.CodeLength-noExceptionExitCodeLength)) // exit is here
		endPc = code.CodeLength
		index--
	}
}
