package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type AutoVariableForRangeJavaArray struct {
	Elements uint16
	End      uint16
	K, V     uint16
}

func (this *BuildPackage) buildForRangeStatementForJavaArray(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	s *ast.StatementFor,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	//build array expression
	maxStack = this.BuildExpression.build(class, code, s.RangeAttr.RangeOn, context, state) // array on stack
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	noNullExit := (&cg.Exit{}).Init(cg.OP_ifnonnull, code)
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
	writeExits([]*cg.Exit{noNullExit}, code.CodeLength)
	state.pushStack(class, s.RangeAttr.RangeOn.Value)
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	forState := (&StackMapState{}).initFromLast(state)
	defer state.addTop(forState) // add top
	attr := s.RangeAttr
	var autoVar AutoVariableForRangeJavaArray
	{
		autoVar.Elements = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, attr.RangeOn.Value)
		// K
		autoVar.K = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})
		//end
		autoVar.End = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})

	}

	//get length
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.CodeLength++
	if 2 > maxStack {
		maxStack = 2
	}
	code.Codes[code.CodeLength] = cg.OP_arraylength
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.End)...)
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeJavaArray, autoVar.Elements)...)

	code.Codes[code.CodeLength] = cg.OP_iconst_m1
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)

	//handle captured vars
	if s.Condition.Type == ast.ExpressionTypeVarAssign {
		if attr.IdentifierValue != nil &&
			attr.IdentifierValue.Variable.BeenCapturedAsLeftValue > 0 {
			closure.createClosureVar(class, code, attr.IdentifierValue.Variable.Type)
			attr.IdentifierValue.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, attr.IdentifierValue.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(attr.RangeOn.Value.Array.Type).className))
		}
		if attr.IdentifierKey != nil &&
			attr.IdentifierKey.Variable.BeenCapturedAsLeftValue > 0 {
			closure.createClosureVar(class, code, attr.IdentifierKey.Variable.Type)
			attr.IdentifierKey.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, attr.IdentifierKey.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(ast.VariableTypeInt).className))
		}
	}
	s.ContinueCodeOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	blockState := (&StackMapState{}).initFromLast(forState)

	code.Codes[code.CodeLength] = cg.OP_iinc
	if autoVar.K > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(autoVar.K)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3

	// load  k
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)

	// load end
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.End)...)
	if 2 > maxStack {
		maxStack = 2
	}
	s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_if_icmpge, code))
	//load elements
	if attr.IdentifierValue != nil || attr.ExpressionValue != nil {
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, autoVar.Elements)...)
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)
		if 2 > maxStack {
			maxStack = 2
		}
		// load value
		switch attr.RangeOn.Value.Array.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			code.Codes[code.CodeLength] = cg.OP_baload
		case ast.VariableTypeShort:
			code.Codes[code.CodeLength] = cg.OP_saload
		case ast.VariableTypeChar:
			code.Codes[code.CodeLength] = cg.OP_caload
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_iaload
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_laload
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_faload
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_daload
		default:
			code.Codes[code.CodeLength] = cg.OP_aaload
		}
		code.CodeLength++
		// v
		autoVar.V = code.MaxLocals
		code.MaxLocals += jvmSlotSize(attr.RangeOn.Value.Array)
		//store to v tmp
		copyOPs(code,
			storeLocalVariableOps(attr.RangeOn.Value.Array.Type,
				autoVar.V)...)

		blockState.appendLocals(class, attr.RangeOn.Value.Array)
	}
	//current stack is 0
	if s.Condition.Type == ast.ExpressionTypeVarAssign {
		if attr.IdentifierValue != nil {
			if attr.IdentifierValue.Variable.BeenCapturedAsLeftValue > 0 {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, attr.IdentifierValue.Variable.LocalValOffset)...)
				copyOPs(code,
					loadLocalVariableOps(attr.RangeOn.Value.Array.Type,
						autoVar.V)...)
				this.storeLocalVar(class, code, attr.IdentifierValue.Variable)
			} else {
				attr.IdentifierValue.Variable.LocalValOffset = autoVar.V
			}
		}
		if attr.IdentifierKey != nil {
			if attr.IdentifierKey.Variable.BeenCapturedAsLeftValue > 0 {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject,
					attr.IdentifierKey.Variable.LocalValOffset)...)
				copyOPs(code,
					loadLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)
				this.storeLocalVar(class, code, attr.IdentifierKey.Variable)
			} else {
				attr.IdentifierKey.Variable.LocalValOffset = autoVar.K
			}
		}
	} else { // for k,v = range arr
		// store v
		//get ops,make_node_objects ops ready
		if attr.ExpressionValue != nil {
			stackLength := len(blockState.Stacks)
			stack, remainStack, ops, _ := this.BuildExpression.getLeftValue(class,
				code, attr.ExpressionValue, context, blockState)
			if stack > maxStack {
				maxStack = stack
			}
			//load v
			copyOPs(code, loadLocalVariableOps(attr.RangeOn.Value.Array.Type,
				autoVar.V)...)
			if t := remainStack + jvmSlotSize(attr.RangeOn.Value.Array); t > maxStack {
				maxStack = t
			}
			copyOPs(code, ops...)
			blockState.popStack(len(blockState.Stacks) - stackLength)
		}
		if attr.ExpressionKey != nil { // set to k
			stackLength := len(blockState.Stacks)
			stack, remainStack, ops, _ := this.BuildExpression.getLeftValue(class,
				code, attr.ExpressionKey, context, blockState)
			if stack > maxStack {
				maxStack = stack
			}
			if t := remainStack + 1; t > maxStack {
				maxStack = t
			}
			// load k
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)
			copyOPs(code, ops...)
			blockState.popStack(len(blockState.Stacks) - stackLength)
		}
	}

	// build block
	this.buildBlock(class, code, s.Block, context, blockState)
	forState.addTop(blockState)
	if s.Block.NotExecuteToLastStatement == false {
		jumpTo(code, s.ContinueCodeOffset)
	}

	return
}
