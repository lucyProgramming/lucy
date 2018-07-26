package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type AutoVariableForRangeArray struct {
	Elements      uint16
	Start, End, K uint16
	V             uint16
}

func (buildPackage *BuildPackage) buildForRangeStatementForArray(class *cg.ClassHighLevel,
	code *cg.AttributeCode, s *ast.StatementFor, context *Context, state *StackMapState) (maxStack uint16) {
	//build array expression
	maxStack = buildPackage.BuildExpression.build(class, code, s.RangeAttr.RangeOn, context, state) // array on stack

	// if null skip
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6) // goto pop
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7) //goto for
	code.Codes[code.CodeLength+7] = cg.OP_pop
	state.pushStack(class, s.RangeAttr.RangeOn.Value)
	context.MakeStackMap(code, state, code.CodeLength+7)
	context.MakeStackMap(code, state, code.CodeLength+11)
	state.popStack(1)
	code.CodeLength += 8
	// code.CodeLength += 3
	s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
	forState := (&StackMapState{}).FromLast(state)
	defer func() {
		state.addTop(forState) // add top
	}()
	var autoVar AutoVariableForRangeArray
	{
		// else
		if s.RangeAttr.RangeOn.Value.Type == ast.VariableTypeArray {
			t := &ast.Type{}
			t.Type = ast.VariableTypeJavaArray
			t.Array = s.RangeAttr.RangeOn.Value.Array
			autoVar.Elements = code.MaxLocals
			code.MaxLocals++
			forState.appendLocals(class, t)
		} else {
			autoVar.Elements = code.MaxLocals
			code.MaxLocals++
			forState.appendLocals(class, s.RangeAttr.RangeOn.Value)
		}
		// start
		autoVar.Start = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})
		//end
		autoVar.End = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})
		// K
		autoVar.K = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})
	}

	if s.RangeAttr.RangeOn.Value.Type == ast.VariableTypeArray {
		//get elements
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxStack {
			maxStack = 2
		}
		meta := ArrayMetas[s.RangeAttr.RangeOn.Value.Array.Type]
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "elements",
			Descriptor: meta.elementsFieldDescriptor,
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		if s.RangeAttr.RangeOn.Value.Array.IsPointer() &&
			s.RangeAttr.RangeOn.Value.Array.Type != ast.VariableTypeString {
			code.Codes[code.CodeLength] = cg.OP_checkcast
			t := &ast.Type{}
			t.Type = ast.VariableTypeJavaArray
			t.Array = s.RangeAttr.RangeOn.Value.Array
			class.InsertClassConst(Descriptor.typeDescriptor(t), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}

		copyOPs(code, storeLocalVariableOps(ast.VariableTypeJavaArray, autoVar.Elements)...)
		//get start
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "start",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.Start)...)
		//get end
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "end",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.End)...)
	} else { // java_array
		//get length
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxStack {
			maxStack = 2
		}
		code.Codes[code.CodeLength+1] = cg.OP_arraylength
		code.CodeLength += 2
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.End)...)
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeJavaArray, autoVar.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.Start)...)
	}

	// k set to  -1
	code.Codes[code.CodeLength] = cg.OP_iconst_m1
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)

	//handle captured vars
	if s.Condition.Type == ast.ExpressionTypeColonAssign {
		if s.RangeAttr.IdentifierValue != nil && s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierValue.Variable.Type)
			s.RangeAttr.IdentifierValue.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.RangeOn.Value.Array.Type).className))
		}
		if s.RangeAttr.IdentifierKey != nil && s.RangeAttr.IdentifierKey.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierKey.Variable.Type)
			s.RangeAttr.IdentifierKey.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, s.RangeAttr.IdentifierKey.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(ast.VariableTypeInt).className))
		}
	}

	s.ContinueCodeOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	blockState := (&StackMapState{}).FromLast(forState)
	code.Codes[code.CodeLength] = cg.OP_iinc
	if autoVar.K > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(autoVar.K)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	// load start
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.Start)...)
	// load k
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)
	// mk index
	code.Codes[code.CodeLength] = cg.OP_iadd
	code.Codes[code.CodeLength+1] = cg.OP_dup
	code.CodeLength += 2
	// load end
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.End)...)
	if 3 > maxStack {
		maxStack = 3
	}
	/*
		k + start >= end,break loop,pop index on stack
		check if need to break
	*/
	rangeEnd := (&cg.Exit{}).Init(cg.OP_if_icmpge, code)
	//load elements
	if s.RangeAttr.IdentifierValue != nil || s.RangeAttr.ExpressionValue != nil {
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, autoVar.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		// load value
		switch s.RangeAttr.RangeOn.Value.Array.Type {
		case ast.VariableTypeBool:
			fallthrough
		case ast.VariableTypeByte:
			code.Codes[code.CodeLength] = cg.OP_baload
		case ast.VariableTypeShort:
			code.Codes[code.CodeLength] = cg.OP_saload
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
		case ast.VariableTypeString:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VariableTypeObject:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VariableTypeMap:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VariableTypeArray:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VariableTypeJavaArray:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VariableTypeFunction:
			code.Codes[code.CodeLength] = cg.OP_aaload
		}
		code.CodeLength++
		// v
		autoVar.V = code.MaxLocals
		code.MaxLocals += jvmSlotSize(s.RangeAttr.RangeOn.Value.Array)
		//store to v tmp
		copyOPs(code,
			storeLocalVariableOps(s.RangeAttr.RangeOn.Value.Array.Type,
				autoVar.V)...)

		blockState.appendLocals(class, s.RangeAttr.RangeOn.Value.Array)
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++ // pop  k on stack
	}
	//current stack is 0
	if s.Condition.Type == ast.ExpressionTypeColonAssign {
		if s.RangeAttr.IdentifierValue != nil {
			if s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
				copyOPs(code,
					loadLocalVariableOps(s.RangeAttr.RangeOn.Value.Array.Type,
						autoVar.V)...)
				buildPackage.storeLocalVar(class, code, s.RangeAttr.IdentifierValue.Variable)
			} else {
				s.RangeAttr.IdentifierValue.Variable.LocalValOffset = autoVar.V
			}
		}

		if s.RangeAttr.IdentifierKey != nil {
			if s.RangeAttr.IdentifierKey.Variable.BeenCaptured {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, s.RangeAttr.IdentifierKey.Variable.LocalValOffset)...)
				copyOPs(code,
					loadLocalVariableOps(ast.VariableTypeInt, autoVar.K)...)
				buildPackage.storeLocalVar(class, code, s.RangeAttr.IdentifierKey.Variable)
			} else {
				s.RangeAttr.IdentifierKey.Variable.LocalValOffset = autoVar.K
			}
		}
	} else { // for k,v = range arr
		// store v
		//get ops,make ops ready
		stackLength := len(blockState.Stacks)
		stack, remainStack, ops, target, _ := buildPackage.BuildExpression.getLeftValue(class,
			code, s.RangeAttr.ExpressionValue, context, blockState)
		if stack > maxStack {
			maxStack = stack
		}
		//load v
		copyOPs(code, loadLocalVariableOps(s.RangeAttr.RangeOn.Value.Array.Type,
			autoVar.V)...)
		if t := remainStack + jvmSlotSize(s.RangeAttr.RangeOn.Value.Array); t > maxStack {
			maxStack = t
		}
		if t := remainStack + jvmSlotSize(target); t > maxStack {
			maxStack = t
		}
		copyOPs(code, ops...)
		blockState.popStack(len(blockState.Stacks) - stackLength)
		if s.RangeAttr.ExpressionKey != nil { // set to k
			stackLength := len(blockState.Stacks)
			stack, remainStack, ops, _, _ := buildPackage.BuildExpression.getLeftValue(class,
				code, s.RangeAttr.ExpressionKey, context, blockState)
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
	buildPackage.buildBlock(class, code, s.Block, context, blockState)
	defer forState.addTop(blockState)
	if s.Block.WillNotExecuteToEnd == false {
		jumpTo(cg.OP_goto, code, s.ContinueCodeOffset)
	}

	//pop index on stack
	writeExits([]*cg.Exit{rangeEnd}, code.CodeLength) // jump to here
	forState.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
	context.MakeStackMap(code, forState, code.CodeLength)
	forState.popStack(1)
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	return
}
