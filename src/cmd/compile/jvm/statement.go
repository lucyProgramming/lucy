package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, block *ast.Block, s *ast.Statement,
	context *Context, state *StackMapState) (maxStack uint16) {
	//fmt.Println(s.Pos)
	switch s.Type {
	case ast.StatementTypeExpression:
		maxStack = buildPackage.BuildExpression.build(class, code, s.Expression, context, state)
	case ast.StatementTypeIf:
		s.StatementIf.Exits = []*cg.Exit{} //could compile multi times
		maxStack = buildPackage.buildIfStatement(class, code, s.StatementIf, context, state)
		if len(s.StatementIf.Exits) > 0 {
			writeExits(s.StatementIf.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.StatementTypeBlock:
		blockState := (&StackMapState{}).initFromLast(state)
		s.Block.Exits = []*cg.Exit{}
		buildPackage.buildBlock(class, code, s.Block, context, blockState)
		state.addTop(blockState)
		if len(s.Block.Exits) > 0 {
			writeExits(s.StatementIf.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.StatementTypeFor:
		s.StatementFor.Exits = []*cg.Exit{} //could compile multi times
		maxStack = buildPackage.buildForStatement(class, code, s.StatementFor, context, state)
		writeExits(s.StatementFor.Exits, code.CodeLength)
		context.MakeStackMap(code, state, code.CodeLength)
	case ast.StatementTypeContinue:
		buildPackage.buildDefers(class, code, context, s.StatementContinue.Defers, state)
		jumpTo(code, s.StatementContinue.StatementFor.ContinueCodeOffset)
	case ast.StatementTypeBreak:
		buildPackage.buildDefers(class, code, context, s.StatementBreak.Defers, state)
		exit := (&cg.Exit{}).Init(cg.OP_goto, code)
		if s.StatementBreak.StatementFor != nil {
			s.StatementBreak.StatementFor.Exits = append(s.StatementBreak.StatementFor.Exits, exit)
		} else if s.StatementBreak.StatementSwitch != nil { // switch
			s.StatementBreak.StatementSwitch.Exits = append(s.StatementBreak.StatementSwitch.Exits, exit)
		} else {
			s.StatementBreak.SwitchTemplateBlock.Exits = append(s.StatementBreak.SwitchTemplateBlock.Exits, exit)
		}
	case ast.StatementTypeReturn:
		maxStack = buildPackage.buildReturnStatement(class, code,
			s.StatementReturn, context, state)
	case ast.StatementTypeSwitch:
		s.StatementSwitch.Exits = []*cg.Exit{} //could compile multi times
		maxStack = buildPackage.buildSwitchStatement(class, code, s.StatementSwitch, context, state)
		if len(s.StatementSwitch.Exits) > 0 {
			if code.CodeLength == context.lastStackMapOffset {
				code.Codes[code.CodeLength] = cg.OP_nop
				code.CodeLength++
			}
			writeExits(s.StatementSwitch.Exits, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
	case ast.StatementTypeGoTo:
		buildPackage.buildDefers(class, code, context, s.StatementGoTo.Defers, state)
		if s.StatementGoTo.StatementLabel.CodeOffsetGenerated {
			jumpTo(code, s.StatementGoTo.StatementLabel.CodeOffset)
		} else {
			exit := (&cg.Exit{}).Init(cg.OP_goto, code)
			s.StatementGoTo.StatementLabel.Exits = append(s.StatementGoTo.StatementLabel.Exits, exit)
		}
	case ast.StatementTypeLabel:
		s.StatementLabel.CodeOffsetGenerated = true
		s.StatementLabel.CodeOffset = code.CodeLength
		if len(s.StatementLabel.Exits) > 0 {
			writeExits(s.StatementLabel.Exits, code.CodeLength) // back patch
		}
		context.MakeStackMap(code, state, code.CodeLength)
	case ast.StatementTypeDefer: // nothing to do  ,defer will do after block is compiled
		s.Defer.StartPc = code.CodeLength
		s.Defer.StackMapState = (&StackMapState{}).initFromLast(state)
	case ast.StatementTypeClass:
		oldName := s.Class.Name
		var name string
		if block.InheritedAttribute.ClassAndFunctionNames == "" {
			name = s.Class.Name
		} else {
			name = block.InheritedAttribute.ClassAndFunctionNames + "$" + s.Class.Name
		}
		s.Class.Name = buildPackage.newClassName(name)
		innerClass := &cg.InnerClass{
			InnerClass:  s.Class.Name,
			OuterClass:  class.Name,
			Name:        oldName,
			AccessFlags: 0,
		}
		class.Class.AttributeInnerClasses.Classes =
			append(class.Class.AttributeInnerClasses.Classes, innerClass)
		c := buildPackage.buildClass(s.Class)
		c.Class.AttributeInnerClasses.Classes = append(c.Class.AttributeInnerClasses.Classes, innerClass)
		buildPackage.putClass(c)
	case ast.StatementTypeNop:
		// nop
	case ast.StatementTypeTypeAlias:
		// handled at ast stage
	}
	return
}

func (buildPackage *BuildPackage) buildDefers(class *cg.ClassHighLevel,
	code *cg.AttributeCode, context *Context, ds []*ast.StatementDefer, from *StackMapState) {
	if len(ds) == 0 {
		return
	}
	code.Codes[code.CodeLength] = cg.OP_aconst_null
	code.CodeLength++
	index := len(ds) - 1
	for index >= 0 { // build defer,cannot have return statement is defer
		state := ds[index].StackMapState.(*StackMapState)
		state = (&StackMapState{}).initFromLast(state) // clone
		state.addTop(from)
		state.pushStack(class, state.newObjectVariableType(throwableClass))
		context.MakeStackMap(code, state, code.CodeLength)
		e := &cg.ExceptionTable{}
		e.StartPc = uint16(ds[index].StartPc)
		e.EndPc = uint16(code.CodeLength)
		e.HandlerPc = uint16(code.CodeLength)
		if ds[index].ExceptionClass == nil {
			e.CatchType = class.Class.InsertClassConst(ast.DefaultExceptionClass)
		} else {
			e.CatchType = class.Class.InsertClassConst(ds[index].ExceptionClass.Name) // custom class
		}
		code.Exceptions = append(code.Exceptions, e)
		//expect exception on stack
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
			context.exceptionVarOffset)...) // this code will make_node_objects stack is empty
		state.popStack(1)
		// build block
		context.Defer = ds[index]
		buildPackage.buildBlock(class, code, &ds[index].Block, context, state)
		from.addTop(state)
		context.Defer = nil
		ds[index].ResetLabels()

		//if need throw
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, context.exceptionVarOffset)...)
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
		if index != 0 {
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
		}
		index--
	}
}
