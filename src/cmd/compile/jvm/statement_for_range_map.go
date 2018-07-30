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

func (buildPackage *BuildPackage) buildForRangeStatementForMap(class *cg.ClassHighLevel, code *cg.AttributeCode,
	s *ast.StatementFor, context *Context, state *StackMapState) (maxStack uint16) {
	maxStack = buildPackage.BuildExpression.build(class, code, s.RangeAttr.RangeOn, context, state) // map instance on stack
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
	s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
	//keySets
	code.Codes[code.CodeLength] = cg.OP_dup
	if 2 > maxStack {
		maxStack = 2
	}
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaMapClass,
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
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})
		t := &ast.Type{}
		t.Type = ast.VariableTypeJavaArray
		t.Array = forState.newObjectVariableType(javaRootClass)
		autoVar.KeySets = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, t)
		autoVar.MapObject = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, forState.newObjectVariableType(javaMapClass))
		autoVar.KeySetsK = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})

	}
	code.Codes[code.CodeLength] = cg.OP_arraylength
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.Length)...)
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, autoVar.KeySets)...)
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, autoVar.MapObject)...)
	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_m1
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsK)...)
	//handle captured vars
	if s.Condition.Type == ast.ExpressionTypeColonAssign {
		if s.RangeAttr.IdentifierValue != nil && s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierValue.Variable.Type)
			s.RangeAttr.IdentifierValue.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.IdentifierValue.Variable.Type.Type).className))
		}
		if s.RangeAttr.IdentifierKey != nil &&
			s.RangeAttr.IdentifierKey.Variable.BeenCaptured {
			closure.createClosureVar(class, code, s.RangeAttr.IdentifierKey.Variable.Type)
			s.RangeAttr.IdentifierKey.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, s.RangeAttr.IdentifierKey.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(s.RangeAttr.IdentifierKey.Variable.Type.Type).className))
		}
	}

	s.ContinueCodeOffset = code.CodeLength
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
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsK)...)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	// load length
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.Length)...)
	if 5 > maxStack {
		maxStack = 5
	}
	exit := (&cg.Exit{}).Init(cg.OP_if_icmpge, code)
	if s.RangeAttr.IdentifierValue != nil || s.RangeAttr.ExpressionValue != nil {
		// load k sets
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, autoVar.KeySets)...)
		// swap
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		//get object for hashMap
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		// load  map object
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, autoVar.MapObject)...)
		// swap
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaMapClass,
			Method:     "get",
			Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if s.RangeAttr.RangeOn.Value.Map.V.IsPointer() == false {
			typeConverter.unPackPrimitives(class, code, s.RangeAttr.RangeOn.Value.Map.V)
		} else {
			typeConverter.castPointer(class, code, s.RangeAttr.RangeOn.Value.Map.V)
		}
		autoVar.V = code.MaxLocals
		code.MaxLocals += jvmSlotSize(s.RangeAttr.RangeOn.Value.Map.V)
		//store to V
		copyOPs(code, storeLocalVariableOps(s.RangeAttr.RangeOn.Value.Map.V.Type, autoVar.V)...)
		blockState.appendLocals(class, s.RangeAttr.RangeOn.Value.Map.V)
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}

	// store to k,if need
	if s.RangeAttr.IdentifierKey != nil || s.RangeAttr.ExpressionKey != nil {
		// load k sets
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, autoVar.KeySets)...)
		// load k
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsK)...)
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		if s.RangeAttr.RangeOn.Value.Map.K.IsPointer() == false {
			typeConverter.unPackPrimitives(class, code, s.RangeAttr.RangeOn.Value.Map.K)
		} else {
			typeConverter.castPointer(class, code, s.RangeAttr.RangeOn.Value.Map.K)
		}
		autoVar.K = code.MaxLocals
		code.MaxLocals += jvmSlotSize(s.RangeAttr.RangeOn.Value.Map.K)
		copyOPs(code, storeLocalVariableOps(s.RangeAttr.RangeOn.Value.Map.K.Type, autoVar.K)...)
		blockState.appendLocals(class, s.RangeAttr.RangeOn.Value.Map.K)
	}

	// store k and v into user defined variable
	//store v in real v
	if s.Condition.Type == ast.ExpressionTypeColonAssign {
		if s.RangeAttr.IdentifierValue != nil {
			if s.RangeAttr.IdentifierValue.Variable.BeenCaptured {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, s.RangeAttr.IdentifierValue.Variable.LocalValOffset)...)
				copyOPs(code,
					loadLocalVariableOps(s.RangeAttr.IdentifierValue.Variable.Type.Type,
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
					loadLocalVariableOps(s.RangeAttr.IdentifierKey.Variable.Type.Type, autoVar.K)...)
				buildPackage.storeLocalVar(class, code, s.RangeAttr.IdentifierKey.Variable)
			} else {
				s.RangeAttr.IdentifierKey.Variable.LocalValOffset = autoVar.K
			}
		}
	} else { // for k,v  = range xxx
		// store v
		if s.RangeAttr.ExpressionValue != nil {
			stackLength := len(blockState.Stacks)
			stack, remainStack, op, _ :=
				buildPackage.BuildExpression.getLeftValue(class, code, s.RangeAttr.ExpressionValue, context, blockState)
			if stack > maxStack { // this means  current stack is 0
				maxStack = stack
			}
			copyOPs(code,
				loadLocalVariableOps(s.RangeAttr.RangeOn.Value.Map.V.Type, autoVar.V)...)
			if t := remainStack + jvmSlotSize(s.RangeAttr.RangeOn.Value.Map.V); t > maxStack {
				maxStack = t
			}
			copyOPs(code, op...)
			forState.popStack(len(blockState.Stacks) - stackLength)
		}
		if s.RangeAttr.ExpressionKey != nil {
			stackLength := len(blockState.Stacks)
			stack, remainStack, op, _ :=
				buildPackage.BuildExpression.getLeftValue(class, code, s.RangeAttr.ExpressionKey, context, blockState)
			if stack > maxStack { // this means  current stack is 0
				maxStack = stack
			}
			copyOPs(code,
				loadLocalVariableOps(s.RangeAttr.RangeOn.Value.Map.K.Type, autoVar.K)...)
			if t := remainStack + jvmSlotSize(s.RangeAttr.RangeOn.Value.Map.K); t > maxStack {
				maxStack = t
			}
			copyOPs(code, op...)
			blockState.popStack(len(blockState.Stacks) - stackLength)
		}

	}
	// build block
	buildPackage.buildBlock(class, code, s.Block, context, blockState)
	forState.addTop(blockState)
	if s.Block.WillNotExecuteToEnd == false {
		jumpTo(cg.OP_goto, code, s.ContinueCodeOffset)
	}
	writeExits([]*cg.Exit{exit}, code.CodeLength)
	{
		forState.pushStack(class,
			&ast.Type{Type: ast.VariableTypeInt})
		context.MakeStackMap(code, forState, code.CodeLength)
		forState.popStack(1)
	}
	// pop 1
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	return
}
