package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type AutoVariableForRangeMap struct {
	MapObject        uint16
	KeySets          uint16
	KeySetsK, Length uint16
	K, V             uint16
}

func (makeClass *MakeClass) buildForRangeStatementForMap(class *cg.ClassHighLevel, code *cg.AttributeCode,
	s *ast.StatementFor, context *Context, state *StackMapState) (maxStack uint16) {
	maxStack, _ = makeClass.makeExpression.build(class, code, s.RangeAttr.RangeOn, context, state) // map instance on stack
	// if null skip
	{
		state.Stacks = append(state.Stacks,
			state.newStackMapVerificationTypeInfo(class, s.RangeAttr.RangeOn.ExpressionValue))
		context.MakeStackMap(code, state, code.CodeLength+7)
		context.MakeStackMap(code, state, code.CodeLength+11)
		state.popStack(1) // pop
	}
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6) // goto pop
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7) //goto for
	code.Codes[code.CodeLength+7] = cg.OP_pop
	code.CodeLength += 8
	forState := (&StackMapState{}).FromLast(state)
	defer state.addTop(forState)
	s.Exits = append(s.Exits, (&cg.Exit{}).FromCode(cg.OP_goto, code))
	//keySets
	code.Codes[code.CodeLength] = cg.OP_dup
	if 2 > maxStack {
		maxStack = 2
	}
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Method:     "keySet",
		Descriptor: "()Ljava/util/Set;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	code.Codes[code.CodeLength] = cg.OP_invokeinterface
	class.InsertInterfaceMethodrefConst(cg.CONSTANT_InterfaceMethodref_info_high_level{
		Class:      "java/util/Set",
		Method:     "toArray",
		Descriptor: "()[Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = 1
	code.Codes[code.CodeLength+4] = 0
	code.CodeLength += 5
	// get length
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 3 > maxStack {
		maxStack = 3
	}
	var autoVar AutoVariableForRangeMap
	{
		autoVar.Length = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
		t := &ast.VariableType{}
		t.Type = ast.VARIABLE_TYPE_JAVA_ARRAY
		t.ArrayType = forState.newObjectVariableType(java_root_class)
		autoVar.KeySets = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, t)
		autoVar.MapObject = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, forState.newObjectVariableType(java_hashmap_class))
		autoVar.KeySetsK = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Type: ast.VARIABLE_TYPE_INT})

	}
	code.Codes[code.CodeLength] = cg.OP_arraylength
	code.CodeLength++
	copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.Length)...)
	copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, autoVar.KeySets)...)
	copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, autoVar.MapObject)...)
	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_m1
	code.CodeLength++
	copyOP(code, storeLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.KeySetsK)...)
	//handle captured vars
	if s.Condition.Type == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.RangeAttr.IdentifierValue != nil && s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierValue.Variable.Type)
			s.RangeAttr.IdentifierValue.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOP(code,
				storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.IdentifierValue.Variable.Type.Type).className))
		}
		if s.RangeAttr.IdentifierKey != nil &&
			s.RangeAttr.IdentifierKey.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierKey.Variable.Type)
			s.RangeAttr.IdentifierKey.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOP(code,
				storeLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierKey.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.IdentifierKey.Variable.Type.Type).className))
		}
	}

	s.ContinueOPOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	blockState := (&StackMapState{}).FromLast(forState)
	code.Codes[code.CodeLength] = cg.OP_iinc
	if autoVar.K > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(autoVar.KeySetsK)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	// load k
	copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.KeySetsK)...)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	// load length
	copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.Length)...)
	if 5 > maxStack {
		maxStack = 5
	}
	exit := (&cg.Exit{}).FromCode(cg.OP_if_icmpge, code)
	if s.RangeAttr.IdentifierValue != nil || s.RangeAttr.ExpressionValue != nil {
		// load k sets
		copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, autoVar.KeySets)...)
		// swap
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		//get object for hashMap
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		// load  map object
		copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, autoVar.MapObject)...)
		// swap
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Method:     "get",
			Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if s.RangeAttr.RangeOn.ExpressionValue.Map.V.IsPointer() == false {
			typeConverter.unPackPrimitives(class, code, s.RangeAttr.RangeOn.ExpressionValue.Map.V)
		} else {
			typeConverter.castPointerTypeToRealType(class, code, s.RangeAttr.RangeOn.ExpressionValue.Map.V)
		}
		autoVar.V = code.MaxLocals
		code.MaxLocals += jvmSize(s.RangeAttr.RangeOn.ExpressionValue.Map.V)
		//store to V
		copyOP(code, storeLocalVariableOps(s.RangeAttr.RangeOn.ExpressionValue.Map.V.Type, autoVar.V)...)
		blockState.appendLocals(class, s.RangeAttr.RangeOn.ExpressionValue.Map.V)
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}

	// store to k,if need
	if s.RangeAttr.IdentifierKey != nil || s.RangeAttr.ExpressionKey != nil {
		// load k sets
		copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, autoVar.KeySets)...)
		// load k
		copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_INT, autoVar.KeySetsK)...)
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		if s.RangeAttr.RangeOn.ExpressionValue.Map.K.IsPointer() == false {
			typeConverter.unPackPrimitives(class, code, s.RangeAttr.RangeOn.ExpressionValue.Map.K)
		} else {
			typeConverter.castPointerTypeToRealType(class, code, s.RangeAttr.RangeOn.ExpressionValue.Map.K)
		}
		autoVar.K = code.MaxLocals
		code.MaxLocals += jvmSize(s.RangeAttr.RangeOn.ExpressionValue.Map.K)
		copyOP(code, storeLocalVariableOps(s.RangeAttr.RangeOn.ExpressionValue.Map.K.Type, autoVar.K)...)
		blockState.appendLocals(class, s.RangeAttr.RangeOn.ExpressionValue.Map.K)
	}

	// store k and v into user defined variable
	//store v in real v
	if s.Condition.Type == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.RangeAttr.IdentifierValue != nil {
			if s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
				copyOP(code, loadLocalVariableOps(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
				copyOP(code,
					loadLocalVariableOps(s.RangeAttr.IdentifierValue.Variable.Type.Type,
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
					loadLocalVariableOps(s.RangeAttr.IdentifierKey.Variable.Type.Type, autoVar.K)...)
				makeClass.storeLocalVar(class, code, s.RangeAttr.IdentifierKey.Variable)
			} else {
				s.RangeAttr.IdentifierKey.Variable.LocalValOffset = autoVar.K
			}
		}
	} else { // for k,v  = range xxx
		// store v
		stackLength := len(blockState.Stacks)
		stack, remainStack, op, _, className, name, descriptor :=
			makeClass.makeExpression.getLeftValue(class, code, s.RangeAttr.ExpressionValue, context, blockState)
		if stack > maxStack { // this means  current stack is 0
			maxStack = stack
		}
		copyOP(code,
			loadLocalVariableOps(s.RangeAttr.RangeOn.ExpressionValue.Map.V.Type, autoVar.V)...)
		if t := remainStack + jvmSize(s.RangeAttr.RangeOn.ExpressionValue.Map.V); t > maxStack {
			maxStack = t
		}
		copyOPLeftValueVersion(class, code, op, className, name, descriptor)
		forState.popStack(len(blockState.Stacks) - stackLength)
		if s.RangeAttr.ExpressionKey != nil {
			stackLength := len(blockState.Stacks)
			stack, remainStack, op, _, className, name, descriptor :=
				makeClass.makeExpression.getLeftValue(class, code, s.RangeAttr.ExpressionKey, context, blockState)
			if stack > maxStack { // this means  current stack is 0
				maxStack = stack
			}
			copyOP(code,
				loadLocalVariableOps(s.RangeAttr.RangeOn.ExpressionValue.Map.K.Type, autoVar.K)...)
			if t := remainStack + jvmSize(s.RangeAttr.RangeOn.ExpressionValue.Map.K); t > maxStack {
				maxStack = t
			}
			copyOPLeftValueVersion(class, code, op, className, name, descriptor)
			blockState.popStack(len(blockState.Stacks) - stackLength)
		}
	}
	// build block
	makeClass.buildBlock(class, code, s.Block, context, blockState)
	defer forState.addTop(blockState)
	if s.Block.DeadEnding == false {
		jumpTo(cg.OP_goto, code, s.ContinueOPOffset)
	}
	backfillExit([]*cg.Exit{exit}, code.CodeLength)

	{
		forState.pushStack(class,
			&ast.VariableType{Type: ast.VARIABLE_TYPE_INT})
		context.MakeStackMap(code, forState, code.CodeLength)
		forState.popStack(1)
	}
	// pop 1
	copyOP(code, []byte{cg.OP_pop}...)
	return
}
