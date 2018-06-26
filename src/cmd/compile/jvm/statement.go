package jvm

import (
	"encoding/binary"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeClass *MakeClass) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, s *ast.Statement,
	context *Context, state *StackMapState) (maxStack uint16) {
	fmt.Println(s.Pos)
	switch s.Type {
	case ast.STATEMENT_TYPE_EXPRESSION:
		if s.Expression.Type == ast.EXPRESSION_TYPE_FUNCTION_LITERAL {
			return makeClass.buildFunctionExpression(class, code, s.Expression, context, state)
		}
		maxStack, _ = makeClass.makeExpression.build(class, code, s.Expression, context, state)
	case ast.STATEMENT_TYPE_IF:
		s.StatementIf.Exits = []*cg.Exit{} //could compile multi times
		maxStack = makeClass.buildIfStatement(class, code, s.StatementIf, context, state)
		if len(s.StatementIf.Exits) > 0 {
			fillOffsetForExits(s.StatementIf.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_BLOCK: //new
		var blockState *StackMapState
		if s.Block.HaveVariableDefinition() {
			blockState = (&StackMapState{}).FromLast(state)
		} else {
			blockState = state
		}
		makeClass.buildBlock(class, code, s.Block, context, blockState)
		state.addTop(blockState)
	case ast.STATEMENT_TYPE_FOR:
		s.StatementFor.Exits = []*cg.Exit{} //could compile multi times
		maxStack = makeClass.buildForStatement(class, code, s.StatementFor, context, state)
		if len(s.StatementFor.Exits) > 0 {
			fillOffsetForExits(s.StatementFor.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_CONTINUE:
		if len(s.StatementContinue.Defers) > 0 {
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			makeClass.buildDefers(class, code, context, s.StatementContinue.Defers, state)
		}
		jumpTo(cg.OP_goto, code, s.StatementContinue.StatementFor.ContinueCodeOffset)
	case ast.STATEMENT_TYPE_BREAK:
		if len(s.StatementBreak.Defers) > 0 {
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			makeClass.buildDefers(class, code, context, s.StatementBreak.Defers, state)
		}
		b := (&cg.Exit{}).FromCode(cg.OP_goto, code)
		if s.StatementBreak.StatementFor != nil {
			s.StatementBreak.StatementFor.Exits = append(s.StatementBreak.StatementFor.Exits, b)
		} else { // switch
			s.StatementBreak.StatementSwitch.Exits = append(s.StatementBreak.StatementSwitch.Exits, b)
		}
	case ast.STATEMENT_TYPE_RETURN:
		maxStack = makeClass.buildReturnStatement(class, code, s.StatementReturn, context, state)
	case ast.STATEMENT_TYPE_SWITCH:
		s.StatementSwitch.Exits = []*cg.Exit{} //could compile multi times
		maxStack = makeClass.buildSwitchStatement(class, code, s.StatementSwitch, context, state)
		if len(s.StatementSwitch.Exits) > 0 {
			if code.CodeLength == context.lastStackMapOffset {
				code.Codes[code.CodeLength] = cg.OP_nop
				code.CodeLength++
			}
			fillOffsetForExits(s.StatementSwitch.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.STATEMENT_TYPE_GOTO:
		if s.StatementGoTo.StatementLabel.CodeOffsetGenerated {
			jumpTo(cg.OP_goto, code, s.StatementGoTo.StatementLabel.CodeOffset)
		} else {
			b := (&cg.Exit{}).FromCode(cg.OP_goto, code)
			s.StatementGoTo.StatementLabel.Exits = append(s.StatementGoTo.StatementLabel.Exits, b)
		}
	case ast.STATEMENT_TYPE_LABEL:
		s.StatementLabel.CodeOffsetGenerated = true
		s.StatementLabel.CodeOffset = code.CodeLength
		s.StatementLabel.Exits = []*cg.Exit{} //could compile multi times
		if len(s.StatementLabel.Exits) > 0 {
			fillOffsetForExits(s.StatementLabel.Exits, code.CodeLength) // back patch
		}
		context.MakeStackMap(code, state, code.CodeLength)
	case ast.STATEMENT_TYPE_DEFER: // nothing to do  ,defer will do after block is compiled
		s.Defer.StartPc = code.CodeLength
		s.Defer.StackMapState = (&StackMapState{}).FromLast(state)
	case ast.STATEMENT_TYPE_CLASS:
		s.Class.Name = makeClass.newClassName(s.Class.Name)
		c := makeClass.buildClass(s.Class)
		makeClass.putClass(c)
	case ast.STATEMENT_TYPE_NOP:
		// nop
	}
	return
}

func (makeClass *MakeClass) buildDefers(class *cg.ClassHighLevel,
	code *cg.AttributeCode, context *Context, ds []*ast.StatementDefer, from *StackMapState) {
	index := len(ds) - 1
	for index >= 0 { // build defer,cannot have return statement is defer
		state := ds[index].StackMapState.(*StackMapState)
		state = (&StackMapState{}).FromLast(state) // clone
		state.addTop(from)
		state.pushStack(class, state.newObjectVariableType(java_throwable_class))
		context.MakeStackMap(code, state, code.CodeLength)
		e := &cg.ExceptionTable{}
		e.StartPc = uint16(ds[index].StartPc)
		e.EndPc = uint16(code.CodeLength)
		e.HandlerPc = uint16(code.CodeLength)
		if ds[index].ExceptionClass == nil {
			e.CatchType = class.Class.InsertClassConst(ast.DEFAULT_EXCEPTION_CLASS)
		} else {
			e.CatchType = class.Class.InsertClassConst(ds[index].ExceptionClass.Name) // custom class
		}
		code.Exceptions = append(code.Exceptions, e)
		//expect exception on stack
		copyOPs(code, storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT,
			context.function.AutoVariableForException.Offset)...) // this code will make stack is empty
		state.popStack(1)
		// build block
		context.Defer = ds[index]
		makeClass.buildBlock(class, code, &ds[index].Block, context, state)
		from.addTop(state)
		context.Defer = nil
		if index > 0 {
			index--
			continue
		}
		//if need throw
		copyOPs(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVariableForException.Offset)...)
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
		index--
	}
}

//func (makeClass *MakeClass) buildDefers(class *cg.ClassHighLevel,
//	code *cg.AttributeCode, context *Context, ds []*ast.StatementDefer, state *StackMapState) {
//	index := len(ds) - 1
//	for index >= 0 {
//		var ss *StackMapState
//		if ds[index].Block.HaveVariableDefinition() {
//			ss = (&StackMapState{}).FromLast(state)
//		} else {
//			ss = state
//		}
//		makeClass.buildBlock(class, code, &ds[index].Block, context, ss)
//		index--
//		state.addTop(ss)
//	}
//}
