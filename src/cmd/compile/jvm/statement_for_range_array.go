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

func (makeClass *MakeClass) buildForRangeStatementForArray(class *cg.ClassHighLevel,
	code *cg.AttributeCode, s *ast.StatementFor, context *Context, state *StackMapState) (maxStack uint16) {
	//build array expression
	maxStack, _ = makeClass.makeExpression.build(class, code, s.RangeAttr.RangeOn, context, state) // array on stack

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
	s.Exits = append(s.Exits, (&cg.Exit{}).FromCode(cg.OP_goto, code))
	forState := (&StackMapState{}).FromLast(state)
	defer func() {
		state.addTop(forState) // add top
	}()
	var autoVar AutoVariableForRangeArray
	{
		// else
		if s.RangeAttr.RangeOn.Value.Type == ast.VARIABLE_TYPE_ARRAY {
			t := &ast.VariableType{}
			t.Type = ast.VARIABLE_TYPE_JAVA_ARRAY
			t.ArrayType = s.RangeAttr.RangeOn.Value.ArrayType
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
		forState.appendLocals(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
		//end
		autoVar.End = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
		// K
		autoVar.K = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
	}

	if s.RangeAttr.RangeOn.Value.Type == ast.VARIABLE_TYPE_ARRAY {
		//get elements
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxStack {
			maxStack = 2
		}
		meta := ArrayMetas[s.RangeAttr.RangeOn.Value.ArrayType.Type]
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "elements",
			Descriptor: meta.elementsFieldDescriptor,
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		if s.RangeAttr.RangeOn.Value.ArrayType.IsPointer() &&
			s.RangeAttr.RangeOn.Value.ArrayType.Type != ast.VARIABLE_TYPE_STRING {
			code.Codes[code.CodeLength] = cg.OP_checkcast
			t := &ast.VariableType{}
			t.Type = ast.VARIABLE_TYPE_JAVA_ARRAY
			t.ArrayType = s.RangeAttr.RangeOn.Value.ArrayType
			class.InsertClassConst(Descriptor.typeDescriptor(t), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}

		copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_JAVA_ARRAY, autoVar.Elements)...)
		//get start
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "start",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.Start)...)
		//get end
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "end",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.End)...)
	} else { // java_array
		//get length
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxStack {
			maxStack = 2
		}
		code.Codes[code.CodeLength+1] = cg.OP_arraylength
		code.CodeLength += 2
		copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.End)...)
		copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_JAVA_ARRAY, autoVar.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
		copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.Start)...)
	}

	// k set to  -1
	code.Codes[code.CodeLength] = cg.OP_iconst_m1
	code.CodeLength++
	copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.K)...)

	//handle captured vars
	if s.Condition.Type == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.RangeAttr.IdentifierValue != nil && s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierValue.Variable.Type)
			s.RangeAttr.IdentifierValue.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOP(code,
				storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.RangeOn.Value.ArrayType.Type).className))
		}
		if s.RangeAttr.IdentifierKey != nil && s.RangeAttr.IdentifierKey.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierKey.Variable.Type)
			s.RangeAttr.IdentifierKey.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOP(code,
				storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierKey.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(ast.VARIABLE_TYPE_INT).className))
		}
	}

	s.ContinueOPOffset = code.CodeLength
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
	copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.Start)...)
	// load k
	copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.K)...)
	// mk index
	code.Codes[code.CodeLength] = cg.OP_iadd
	code.Codes[code.CodeLength+1] = cg.OP_dup
	code.CodeLength += 2
	// load end
	copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.End)...)
	if 3 > maxStack {
		maxStack = 3
	}
	/*
		k + start >= end,break loop,pop index on stack
		check if need to break
	*/
	rangeend := (&cg.Exit{}).FromCode(cg.OP_if_icmpge, code)
	//load elements
	if s.RangeAttr.IdentifierValue != nil || s.RangeAttr.ExpressionValue != nil {
		copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, autoVar.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		// load value
		switch s.RangeAttr.RangeOn.Value.ArrayType.Type {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			code.Codes[code.CodeLength] = cg.OP_baload
		case ast.VARIABLE_TYPE_SHORT:
			code.Codes[code.CodeLength] = cg.OP_saload
		case ast.VARIABLE_TYPE_ENUM:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_iaload
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_laload
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_faload
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_daload
		case ast.VARIABLE_TYPE_STRING:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VARIABLE_TYPE_OBJECT:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VARIABLE_TYPE_MAP:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VARIABLE_TYPE_ARRAY:
			code.Codes[code.CodeLength] = cg.OP_aaload
		case ast.VARIABLE_TYPE_JAVA_ARRAY:
			code.Codes[code.CodeLength] = cg.OP_aaload
		}
		code.CodeLength++
		// v
		autoVar.V = code.MaxLocals
		code.MaxLocals += jvmSize(s.RangeAttr.RangeOn.Value.ArrayType)
		//store to v tmp
		copyOP(code,
			storeLocalVariableOps(s.RangeAttr.RangeOn.Value.ArrayType.Type,
				autoVar.V)...)

		blockState.appendLocals(class, s.RangeAttr.RangeOn.Value.ArrayType)
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++ // pop  k on stack
	}
	//current stack is 0
	if s.Condition.Type == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.RangeAttr.IdentifierValue != nil {
			if s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
				copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
				copyOP(code,
					loadLocalVariableOps(s.RangeAttr.RangeOn.Value.ArrayType.Type,
						autoVar.V)...)
				makeClass.storeLocalVar(class, code, s.RangeAttr.IdentifierValue.Variable)
			} else {
				s.RangeAttr.IdentifierValue.Variable.LocalValOffset = autoVar.V
			}
		}

		if s.RangeAttr.IdentifierKey != nil {
			if s.RangeAttr.IdentifierKey.Variable.BeenCaptured {
				copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierKey.Variable.LocalValOffset)...)
				copyOP(code,
					loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.K)...)
				makeClass.storeLocalVar(class, code, s.RangeAttr.IdentifierKey.Variable)
			} else {
				s.RangeAttr.IdentifierKey.Variable.LocalValOffset = autoVar.K
			}
		}
	} else { // for k,v = range arr
		// store v
		//get ops,make ops ready
		stackLength := len(blockState.Stacks)
		stack, remainStack, ops, target, classname, name, descriptor := makeClass.makeExpression.getLeftValue(class,
			code, s.RangeAttr.ExpressionValue, context, blockState)
		if stack > maxStack {
			maxStack = stack
		}
		//load v
		copyOP(code, loadLocalVariableOps(s.RangeAttr.RangeOn.Value.ArrayType.Type,
			autoVar.V)...)
		if t := remainStack + jvmSize(s.RangeAttr.RangeOn.Value.ArrayType); t > maxStack {
			maxStack = t
		}
		if t := remainStack + jvmSize(target); t > maxStack {
			maxStack = t
		}
		copyOPLeftValueVersion(class, code, ops, classname, name, descriptor)
		blockState.popStack(len(blockState.Stacks) - stackLength)
		if s.RangeAttr.ExpressionKey != nil { // set to k
			stackLength := len(blockState.Stacks)
			stack, remainStack, ops, _, classname, name, descriptor := makeClass.makeExpression.getLeftValue(class,
				code, s.RangeAttr.ExpressionKey, context, blockState)
			if stack > maxStack {
				maxStack = stack
			}
			if t := remainStack + 1; t > maxStack {
				maxStack = t
			}
			// load k
			copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.K)...)
			copyOPLeftValueVersion(class, code, ops, classname, name, descriptor)
			blockState.popStack(len(blockState.Stacks) - stackLength)
		}
	}

	// build block
	makeClass.buildBlock(class, code, s.Block, context, blockState)
	defer forState.addTop(blockState)
	if s.Block.DeadEnding == false {
		jumpTo(cg.OP_goto, code, s.ContinueOPOffset)
	}

	//pop index on stack
	backfillExit([]*cg.Exit{rangeend}, code.CodeLength) // jump to here
	forState.pushStack(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
	context.MakeStackMap(code, forState, code.CodeLength)
	forState.popStack(1)
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	return
}
