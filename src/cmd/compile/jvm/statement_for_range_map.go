package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForRangeStatementForMap(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context, state *StackMapState) (maxstack uint16) {
	maxstack, _ = m.MakeExpression.build(class, code, s.StatmentForRangeAttr.Expression, context, state) // map instance on stack
	// if null skip
	{
		state.Stacks = append(state.Stacks,
			state.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.Expression.VariableType)...)
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
			context.MakeStackMap(state, code.CodeLength+7))
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
			context.MakeStackMap(state, code.CodeLength+11))
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
	code.Codes[code.CodeLength] = cg.OP_arraylength
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsKLength)...)
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySets)...)
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeMap.MapObject)...)
	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsK)...)
	{
		forState.Locals = append(forState.Locals,
			forState.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
		t := &ast.VariableType{}
		t.Typ = ast.VARIABLE_TYPE_JAVA_ARRAY
		t.ArrayType = forState.newObjectVariableType(java_root_class)
		forState.Locals = append(forState.Locals, forState.newStackMapVerificationTypeInfo(class, t)...)
		t2 := forState.newObjectVariableType(java_hashmap_class)
		forState.Locals = append(forState.Locals, forState.newStackMapVerificationTypeInfo(class, t2)...)
		forState.Locals = append(forState.Locals,
			forState.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
	}

	//continue offset start from here
	loopBeginsAt := code.CodeLength
	code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(forState, code.CodeLength))
	// load  map object
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeMap.MapObject)...)
	// load k sets
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySets)...)
	// load k
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsK)...)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	// load length
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsKLength)...)
	if 5 > maxstack {
		maxstack = 5
	}
	exit := (&cg.JumpBackPatch{}).FromCode(cg.OP_if_icmpge, code)
	//get object for hashMap
	code.Codes[code.CodeLength] = cg.OP_aaload
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_hashmap_class,
		Method:     "get",
		Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if s.StatmentForRangeAttr.Expression.VariableType.Map.V.IsPointer() == false {
		primitiveObjectConverter.getFromObject(class, code, s.StatmentForRangeAttr.Expression.VariableType.Map.V)
	} else {
		primitiveObjectConverter.castPointerTypeToRealType(class, code, s.StatmentForRangeAttr.Expression.VariableType.Map.V)
	}
	//store to V
	copyOP(code, storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ, s.StatmentForRangeAttr.AutoVarForRangeMap.V)...)
	forState.Locals = append(forState.Locals, state.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.Expression.VariableType.Map.V)...)
	// store to k,if need
	if s.StatmentForRangeAttr.ModelKV {
		// load k sets
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySets)...)
		// load k
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsK)...)
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		if s.StatmentForRangeAttr.Expression.VariableType.Map.K.IsPointer() == false {
			primitiveObjectConverter.getFromObject(class, code, s.StatmentForRangeAttr.Expression.VariableType.Map.K)
		} else {
			primitiveObjectConverter.castPointerTypeToRealType(class, code, s.StatmentForRangeAttr.Expression.VariableType.Map.K)
		}
		copyOP(code, storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.K.Typ, s.StatmentForRangeAttr.AutoVarForRangeMap.K)...)
		forState.Locals = append(forState.Locals, state.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.Expression.VariableType.Map.K)...)
	}
	// store k and v into user defined variable
	//store v in real v
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.StatmentForRangeAttr.IdentifierV.Var.BeenCaptured {
			closure.createCloureVar(class, code, s.StatmentForRangeAttr.IdentifierV.Var)
			// load to stack
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
			copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ,
				s.StatmentForRangeAttr.AutoVarForRangeMap.V)...)
			closure.storeLocalCloureVar(class, code, s.StatmentForRangeAttr.IdentifierV.Var)
			{
				t := &ast.VariableType{}
				t.Typ = ast.VARIABLE_TYPE_OBJECT
				t.Class = &ast.Class{}
				t.Class.Name = closure.getMeta(s.StatmentForRangeAttr.IdentifierV.Var.Typ.Typ).className
				forState.Locals = append(forState.Locals,
					forState.newStackMapVerificationTypeInfo(class, t)...)
			}
		} else {
			// load v
			copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ,
				s.StatmentForRangeAttr.AutoVarForRangeMap.V)...)
			copyOP(code, storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ,
				s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
			forState.Locals = append(forState.Locals,
				forState.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.IdentifierV.Var.Typ)...)
		}
		if s.StatmentForRangeAttr.ModelKV {
			if s.StatmentForRangeAttr.IdentifierK.Var.BeenCaptured {
				closure.createCloureVar(class, code, s.StatmentForRangeAttr.IdentifierK.Var)
				copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.IdentifierK.Var.Typ.Typ,
					s.StatmentForRangeAttr.AutoVarForRangeMap.K)...)
				copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.K.Typ,
					s.StatmentForRangeAttr.AutoVarForRangeMap.K)...)
				closure.storeLocalCloureVar(class, code, s.StatmentForRangeAttr.IdentifierV.Var)
				{
					t := &ast.VariableType{}
					t.Typ = ast.VARIABLE_TYPE_OBJECT
					t.Class = &ast.Class{}
					t.Class.Name = closure.getMeta(ast.VARIABLE_TYPE_INT).className
					forState.Locals = append(forState.Locals,
						forState.newStackMapVerificationTypeInfo(class, t)...)
				}
			} else {
				copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.K.Typ,
					s.StatmentForRangeAttr.AutoVarForRangeMap.K)...)
				copyOP(code, storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.K.Typ,
					s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset)...)
				forState.Locals = append(forState.Locals,
					forState.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.IdentifierK.Var.Typ)...)
			}
		}
	} else { // for k,v  = range xxx
		// store v
		stack, remainStack, op, _, classname, name, descriptor := m.MakeExpression.getLeftValue(class, code, s.StatmentForRangeAttr.ExpressionV, context, state)
		if stack > maxstack { // this means  current stack is 0
			maxstack = stack
		}
		copyOP(code,
			loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ, s.StatmentForRangeAttr.AutoVarForRangeMap.V)...)
		if t := remainStack + s.StatmentForRangeAttr.Expression.VariableType.Map.V.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		copyOPLeftValue(class, code, op, classname, name, descriptor)
		if s.StatmentForRangeAttr.ModelKV {
			stack, remainStack, op, _, classname, name, descriptor := m.MakeExpression.getLeftValue(class, code, s.StatmentForRangeAttr.ExpressionK, context, state)
			if stack > maxstack { // this means  current stack is 0
				maxstack = stack
			}
			copyOP(code,
				loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.K.Typ, s.StatmentForRangeAttr.AutoVarForRangeMap.K)...)
			if t := remainStack + s.StatmentForRangeAttr.Expression.VariableType.Map.K.JvmSlotSize(); t > maxstack {
				maxstack = t
			}
			if classname == java_hashmap_class {
				// put in object
				primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, s.StatmentForRangeAttr.Expression.VariableType.Map.K)
			}
			copyOPLeftValue(class, code, op, classname, name, descriptor)
		}
	}
	// build block
	m.buildBlock(class, code, s.Block, context, forState)
	s.ContinueOPOffset = code.CodeLength
	code.Codes[code.CodeLength] = cg.OP_iinc
	if s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsK > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsK)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	jumpto(cg.OP_goto, code, loopBeginsAt)
	backPatchEs([]*cg.JumpBackPatch{exit}, code.CodeLength)

	{
		// object ref
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(java_hashmap_class))...)
		t := &ast.VariableType{}
		t.Typ = ast.VARIABLE_TYPE_JAVA_ARRAY
		t.ArrayType = state.newObjectVariableType(java_root_class)
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, t)...)
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
		code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(state, code.CodeLength))
		state.popStack(3)
	}

	// pop 3
	copyOP(code, []byte{cg.OP_pop, cg.OP_pop, cg.OP_pop}...)
	return
}
