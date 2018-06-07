package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type AutoVarForRangeMap struct {
	MapObject        uint16
	KeySets          uint16
	KeySetsK, Length uint16
	K, V             uint16
}

func (m *MakeClass) buildForRangeStatementForMap(class *cg.ClassHighLevel, code *cg.AttributeCode,
	s *ast.StatementFor, context *Context, state *StackMapState) (maxstack uint16) {
	maxstack, _ = m.MakeExpression.build(class, code, s.RangeAttr.RangeOn, context, state) // map instance on stack
	// if null skip
	{
		state.Stacks = append(state.Stacks,
			state.newStackMapVerificationTypeInfo(class, s.RangeAttr.RangeOn.Value))
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
	s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	//keySets
	code.Codes[code.CodeLength] = cg.OP_dup
	if 2 > maxstack {
		maxstack = 2
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
	if 3 > maxstack {
		maxstack = 3
	}
	var autoVar AutoVarForRangeMap
	{
		autoVar.Length = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
		t := &ast.VariableType{}
		t.Typ = ast.VARIABLE_TYPE_JAVA_ARRAY
		t.ArrayType = forState.newObjectVariableType(java_root_class)
		autoVar.KeySets = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, t)
		autoVar.MapObject = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, forState.newObjectVariableType(java_hashmap_class))
		autoVar.KeySetsK = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})

	}
	code.Codes[code.CodeLength] = cg.OP_arraylength
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.Length)...)
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, autoVar.KeySets)...)
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, autoVar.MapObject)...)
	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_m1
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.KeySetsK)...)

	//handle captured vars
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.RangeAttr.IdentifierV != nil && s.RangeAttr.IdentifierV.Var.BeenCaptured {
			closure.createCloureVar(class, code, s.RangeAttr.IdentifierV.Var.Typ)
			s.RangeAttr.IdentifierV.Var.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOP(code,
				storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierV.Var.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.IdentifierV.Var.Typ.Typ).className))
		}
		if s.RangeAttr.IdentifierK != nil &&
			s.RangeAttr.IdentifierK.Var.BeenCaptured {
			closure.createCloureVar(class, code, s.RangeAttr.IdentifierK.Var.Typ)
			s.RangeAttr.IdentifierK.Var.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOP(code,
				storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierK.Var.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.IdentifierK.Var.Typ.Typ).className))
		}
	}

	context.MakeStackMap(code, forState, code.CodeLength)
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
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.KeySetsK)...)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	// load length
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.Length)...)
	if 5 > maxstack {
		maxstack = 5
	}
	exit := (&cg.JumpBackPatch{}).FromCode(cg.OP_if_icmpge, code)
	if s.RangeAttr.IdentifierV != nil || s.RangeAttr.ExpressionV != nil {
		// load k sets
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, autoVar.KeySets)...)
		// swap
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		//get object for hashMap
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		// load  map object
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, autoVar.MapObject)...)
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
		if s.RangeAttr.RangeOn.Value.Map.V.IsPointer() == false {
			typeConverter.getFromObject(class, code, s.RangeAttr.RangeOn.Value.Map.V)
		} else {
			typeConverter.castPointerTypeToRealType(class, code, s.RangeAttr.RangeOn.Value.Map.V)
		}
		autoVar.V = code.MaxLocals
		code.MaxLocals += jvmSize(s.RangeAttr.RangeOn.Value.Map.V)
		//store to V
		copyOP(code, storeSimpleVarOp(s.RangeAttr.RangeOn.Value.Map.V.Typ, autoVar.V)...)
		blockState.appendLocals(class, s.RangeAttr.RangeOn.Value.Map.V)
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}

	// store to k,if need
	if s.RangeAttr.IdentifierK != nil || s.RangeAttr.ExpressionK != nil {
		// load k sets
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, autoVar.KeySets)...)
		// load k
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.KeySetsK)...)
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		if s.RangeAttr.RangeOn.Value.Map.K.IsPointer() == false {
			typeConverter.getFromObject(class, code, s.RangeAttr.RangeOn.Value.Map.K)
		} else {
			typeConverter.castPointerTypeToRealType(class, code, s.RangeAttr.RangeOn.Value.Map.K)
		}
		autoVar.K = code.MaxLocals
		code.MaxLocals += jvmSize(s.RangeAttr.RangeOn.Value.Map.K)
		copyOP(code, storeSimpleVarOp(s.RangeAttr.RangeOn.Value.Map.K.Typ, autoVar.K)...)
		blockState.appendLocals(class, s.RangeAttr.RangeOn.Value.Map.K)
	}

	// store k and v into user defined variable
	//store v in real v
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.RangeAttr.IdentifierV != nil {
			if s.RangeAttr.IdentifierV.Var.BeenCaptured {
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierV.Var.LocalValOffset)...)
				copyOP(code,
					loadSimpleVarOp(s.RangeAttr.IdentifierV.Var.Typ.Typ,
						autoVar.V)...)
				m.storeLocalVar(class, code, s.RangeAttr.IdentifierV.Var)
			} else {
				s.RangeAttr.IdentifierV.Var.LocalValOffset = autoVar.V
			}
		}
		if s.RangeAttr.IdentifierK != nil {
			if s.RangeAttr.IdentifierK.Var.BeenCaptured {
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.RangeAttr.IdentifierK.Var.LocalValOffset)...)
				copyOP(code,
					loadSimpleVarOp(s.RangeAttr.IdentifierK.Var.Typ.Typ, autoVar.K)...)
				m.storeLocalVar(class, code, s.RangeAttr.IdentifierK.Var)
			} else {
				s.RangeAttr.IdentifierK.Var.LocalValOffset = autoVar.K
			}
		}
	} else { // for k,v  = range xxx
		// store v
		stackLength := len(blockState.Stacks)
		stack, remainStack, op, _, classname, name, descriptor :=
			m.MakeExpression.getLeftValue(class, code, s.RangeAttr.ExpressionV, context, blockState)
		if stack > maxstack { // this means  current stack is 0
			maxstack = stack
		}
		copyOP(code,
			loadSimpleVarOp(s.RangeAttr.RangeOn.Value.Map.V.Typ, autoVar.V)...)
		if t := remainStack + jvmSize(s.RangeAttr.RangeOn.Value.Map.V); t > maxstack {
			maxstack = t
		}
		copyOPLeftValue(class, code, op, classname, name, descriptor)
		forState.popStack(len(blockState.Stacks) - stackLength)
		if s.RangeAttr.ExpressionK != nil {
			stackLength := len(blockState.Stacks)
			stack, remainStack, op, _, classname, name, descriptor :=
				m.MakeExpression.getLeftValue(class, code, s.RangeAttr.ExpressionK, context, blockState)
			if stack > maxstack { // this means  current stack is 0
				maxstack = stack
			}
			copyOP(code,
				loadSimpleVarOp(s.RangeAttr.RangeOn.Value.Map.K.Typ, autoVar.K)...)
			if t := remainStack + jvmSize(s.RangeAttr.RangeOn.Value.Map.K); t > maxstack {
				maxstack = t
			}
			if classname == java_hashmap_class {
				// put in object
				typeConverter.putPrimitiveInObject(class, code,
					s.RangeAttr.RangeOn.Value.Map.K)
			}
			copyOPLeftValue(class, code, op, classname, name, descriptor)
			blockState.popStack(len(blockState.Stacks) - stackLength)
		}
	}
	// build block

	m.buildBlock(class, code, s.Block, context, blockState)
	defer forState.addTop(blockState)
	if s.Block.DeadEnding == false {
		jumpTo(cg.OP_goto, code, s.ContinueOPOffset)
	}
	backPatchEs([]*cg.JumpBackPatch{exit}, code.CodeLength)

	{
		forState.pushStack(class,
			&ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
		context.MakeStackMap(code, forState, code.CodeLength)
		forState.popStack(1)
	}
	// pop 1
	copyOP(code, []byte{cg.OP_pop}...)
	return
}
