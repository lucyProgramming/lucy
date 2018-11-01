package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type AutoVariableForRangeMap struct {
	MapObject               uint16
	KeySets                 uint16
	KeySetsK, KeySetsLength uint16
	K, V                    uint16
}

func (buildPackage *BuildPackage) buildForRangeStatementForMap(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	s *ast.StatementFor,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	attr := s.RangeAttr
	maxStack = buildPackage.BuildExpression.build(
		class,
		code,
		attr.RangeOn,
		context,
		state) // map instance on stack
	// if null skip
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	noNullExit := (&cg.Exit{}).Init(cg.OP_ifnonnull, code)
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	s.Exits = append(s.Exits, (&cg.Exit{}).Init(cg.OP_goto, code))
	writeExits([]*cg.Exit{noNullExit}, code.CodeLength)
	state.pushStack(class, attr.RangeOn.Value)
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	forState := (&StackMapState{}).initFromLast(state)
	defer state.addTop(forState) // add top
	//keySets
	code.Codes[code.CodeLength] = cg.OP_dup
	if 2 > maxStack {
		maxStack = 2
	}
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      mapClass,
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
		autoVar.KeySetsLength = code.MaxLocals
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
		forState.appendLocals(class, forState.newObjectVariableType(mapClass))
		autoVar.KeySetsK = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.Type{Type: ast.VariableTypeInt})

	}
	code.Codes[code.CodeLength] = cg.OP_arraylength
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsLength)...)
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, autoVar.KeySets)...)
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, autoVar.MapObject)...)
	// k set to -1
	code.Codes[code.CodeLength] = cg.OP_iconst_m1
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsK)...)
	//handle captured vars
	if s.Condition.Type == ast.ExpressionTypeVarAssign {
		if attr.IdentifierValue != nil && attr.IdentifierValue.Variable.BeenCapturedAsLeftValue > 0 {
			closure.createClosureVar(class, code, attr.IdentifierValue.Variable.Type)
			attr.IdentifierValue.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, attr.IdentifierValue.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(attr.IdentifierValue.Variable.Type.Type).className))
		}
		if attr.IdentifierKey != nil &&
			attr.IdentifierKey.Variable.BeenCapturedAsLeftValue > 0 {
			closure.createClosureVar(class, code, attr.IdentifierKey.Variable.Type)
			attr.IdentifierKey.Variable.LocalValOffset = code.MaxLocals
			code.MaxLocals++
			copyOPs(code,
				storeLocalVariableOps(ast.VariableTypeObject, attr.IdentifierKey.Variable.LocalValOffset)...)
			forState.appendLocals(class,
				forState.newObjectVariableType(closure.getMeta(attr.IdentifierKey.Variable.Type.Type).className))
		}
	}

	s.ContinueCodeOffset = code.CodeLength
	context.MakeStackMap(code, forState, code.CodeLength)
	blockState := (&StackMapState{}).initFromLast(forState)
	code.Codes[code.CodeLength] = cg.OP_iinc
	if autoVar.K > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(autoVar.KeySetsK)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	// load k
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsK)...)

	// load length
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsLength)...)
	if 2 > maxStack {
		maxStack = 2
	}
	exit := (&cg.Exit{}).Init(cg.OP_if_icmpge, code)
	if attr.IdentifierValue != nil || attr.ExpressionValue != nil {
		// load k sets
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, autoVar.KeySets)...)

		// swap
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsK)...)
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
			Class:      mapClass,
			Method:     "get",
			Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if attr.RangeOn.Value.Map.V.IsPointer() == false {
			typeConverter.unPackPrimitives(class, code, attr.RangeOn.Value.Map.V)
		} else {
			typeConverter.castPointer(class, code, attr.RangeOn.Value.Map.V)
		}
		autoVar.V = code.MaxLocals
		code.MaxLocals += jvmSlotSize(attr.RangeOn.Value.Map.V)
		//store to V
		copyOPs(code, storeLocalVariableOps(attr.RangeOn.Value.Map.V.Type, autoVar.V)...)
		blockState.appendLocals(class, attr.RangeOn.Value.Map.V)
	}

	// store to k,if need
	if attr.IdentifierKey != nil || attr.ExpressionKey != nil {
		// load k sets
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, autoVar.KeySets)...)
		// load k
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, autoVar.KeySetsK)...)
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
		if attr.RangeOn.Value.Map.K.IsPointer() == false {
			typeConverter.unPackPrimitives(class, code, attr.RangeOn.Value.Map.K)
		} else {
			typeConverter.castPointer(class, code, attr.RangeOn.Value.Map.K)
		}
		autoVar.K = code.MaxLocals
		code.MaxLocals += jvmSlotSize(attr.RangeOn.Value.Map.K)
		copyOPs(code, storeLocalVariableOps(attr.RangeOn.Value.Map.K.Type, autoVar.K)...)
		blockState.appendLocals(class, attr.RangeOn.Value.Map.K)
	}

	// store k and v into user defined variable
	//store v in real v
	if s.Condition.Type == ast.ExpressionTypeVarAssign {
		if attr.IdentifierValue != nil {
			if attr.IdentifierValue.Variable.BeenCapturedAsLeftValue > 0 {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, attr.IdentifierValue.Variable.LocalValOffset)...)
				copyOPs(code,
					loadLocalVariableOps(attr.IdentifierValue.Variable.Type.Type,
						autoVar.V)...)
				buildPackage.storeLocalVar(class, code, attr.IdentifierValue.Variable)
			} else {
				attr.IdentifierValue.Variable.LocalValOffset = autoVar.V
			}
		}
		if attr.IdentifierKey != nil {
			if attr.IdentifierKey.Variable.BeenCapturedAsLeftValue > 0 {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, attr.IdentifierKey.Variable.LocalValOffset)...)
				copyOPs(code,
					loadLocalVariableOps(attr.IdentifierKey.Variable.Type.Type, autoVar.K)...)
				buildPackage.storeLocalVar(class, code, attr.IdentifierKey.Variable)
			} else {
				attr.IdentifierKey.Variable.LocalValOffset = autoVar.K
			}
		}
	} else { // for k,v  = range xxx
		// store v
		if attr.ExpressionValue != nil {
			stackLength := len(blockState.Stacks)
			stack, remainStack, op, _ :=
				buildPackage.BuildExpression.getLeftValue(class, code, attr.ExpressionValue, context, blockState)
			if stack > maxStack { // this means  current stack is 0
				maxStack = stack
			}
			copyOPs(code,
				loadLocalVariableOps(attr.RangeOn.Value.Map.V.Type, autoVar.V)...)
			if t := remainStack + jvmSlotSize(attr.RangeOn.Value.Map.V); t > maxStack {
				maxStack = t
			}
			copyOPs(code, op...)
			forState.popStack(len(blockState.Stacks) - stackLength)
		}
		if attr.ExpressionKey != nil {
			stackLength := len(blockState.Stacks)
			stack, remainStack, op, _ :=
				buildPackage.BuildExpression.getLeftValue(class, code, attr.ExpressionKey, context, blockState)
			if stack > maxStack { // this means  current stack is 0
				maxStack = stack
			}
			copyOPs(code,
				loadLocalVariableOps(attr.RangeOn.Value.Map.K.Type, autoVar.K)...)
			if t := remainStack + jvmSlotSize(attr.RangeOn.Value.Map.K); t > maxStack {
				maxStack = t
			}
			copyOPs(code, op...)
			blockState.popStack(len(blockState.Stacks) - stackLength)
		}
	}
	// build block
	buildPackage.buildBlock(class, code, s.Block, context, blockState)
	forState.addTop(blockState)
	if s.Block.NotExecuteToLastStatement == false {
		jumpTo(code, s.ContinueCodeOffset)
	}
	writeExits([]*cg.Exit{exit}, code.CodeLength)
	context.MakeStackMap(code, forState, code.CodeLength)
	return
}
