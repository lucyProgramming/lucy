package jvm

import (
	"encoding/binary"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type HasDefer struct {
	has         bool
	returnExits []*cg.JumpBackPatch
}

func (m *MakeClass) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context) {
	var startPc, endPc int
	hasDefer := &HasDefer{
		has: len(b.Defers) > 0,
	}
	if hasDefer.has {
		startPc = code.CodeLength
		context.function.MkAutoVarForReturnBecauseOfDefer()
		loadInt32(class, code, 0)
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, context.function.AutoVarForReturnBecauseOfDefer.Returnd)...)
	}
	for _, s := range b.Statements {
		maxstack := m.buildStatement(class, code, s, context, hasDefer)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
	}
	if len(b.Defers) != 0 {
		endPc = code.CodeLength
		//execute form begin to  end , no exceptions,push a null exception for defer to catch
		code.Codes[code.CodeLength] = cg.OP_aconst_null // let defer to catch
		code.CodeLength++
	}
	if hasDefer.has && len(hasDefer.returnExits) > 0 { // return should goto here and execute
		backPatchEs(hasDefer.returnExits, code.CodeLength)
	}
	index := len(b.Defers) - 1
	for index >= 0 { // build defer,cannot have return statement is defer
		// insert exceptions
		e := &cg.ExceptionTable{}
		e.StartPc = uint16(startPc)
		e.Endpc = uint16(endPc)
		e.HandlerPc = uint16(code.CodeLength)
		e.CatchType = class.Class.InsertClassConst("java/lang/Throwable")
		code.Exceptions = append(code.Exceptions, e)
		startPc = code.CodeLength
		hasDefer := &HasDefer{}
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForException.Offset)...)
		for _, s := range b.Defers[index].Block.Statements {
			maxstack := m.buildStatement(class, code, s, context, hasDefer)
			if maxstack > code.MaxStack {
				code.MaxStack = maxstack
			}
		}
		// load to stack
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, b.InheritedAttribute.Function.AutoVarForException.Offset)...)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_ifnonnull
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		code.Codes[code.CodeLength+3] = cg.OP_pop // pop exception on stack
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		code.Codes[code.CodeLength+7] = cg.OP_athrow
		code.CodeLength += 8
		endPc = code.CodeLength
		index--
	}
	//after defer executed ,check is need to return
	if hasDefer.has && len(hasDefer.returnExits) > 0 { // could be return
		//load return to stack
		f := func() {
			//this is  real return,after defer executed
			if context.function.HaveNoReturnValue() {
				code.Codes[code.CodeLength] = cg.OP_return
				code.CodeLength++
			} else {
				panic("......")
			}
		}
		if b.IsFunctionTopBlock {
			f()
		} else {
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, context.function.AutoVarForReturnBecauseOfDefer.Returnd)...)
			code.Codes[code.CodeLength] = cg.OP_ifne
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6) // goto return
			code.Codes[code.CodeLength+3] = cg.OP_goto
			offset := code.CodeLength + 3
			exit := code.Codes[code.CodeLength+4 : code.CodeLength+6]
			code.CodeLength += 6
			f()
			binary.BigEndian.PutUint16(exit, uint16(code.CodeLength-offset))
		}
	}

	return
}
