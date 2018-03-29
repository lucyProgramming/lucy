package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForRangeStatementForMap(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	maxstack, _ = m.MakeExpression.build(class, code, s.StatmentForRangeAttr.Expression, context) // map instance on stack
	// if null skip
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6) // goto pop
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7) //goto for
	code.Codes[code.CodeLength+7] = cg.OP_pop
	code.CodeLength += 8
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
	//continue offset start from here
	loopBeginsAt := code.CodeLength
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
	}
	// store k and v into user defined variable
	//store v in real v
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.StatmentForRangeAttr.IdentifierV.Var.BeenCaptured {
			closure.createCloureVar(class, code, s.StatmentForRangeAttr.IdentifierV.Var)
			// load to stack
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
			copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ, s.StatmentForRangeAttr.AutoVarForRangeMap.V)...)
			closure.storeLocalCloureVar(class, code, s.StatmentForRangeAttr.IdentifierV.Var)
		} else {
			// load v
			copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ, s.StatmentForRangeAttr.AutoVarForRangeMap.V)...)
			copyOP(code, storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ, s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
		}
		if s.StatmentForRangeAttr.ModelKV {
			if s.StatmentForRangeAttr.IdentifierK.Var.BeenCaptured {
				closure.createCloureVar(class, code, s.StatmentForRangeAttr.IdentifierK.Var)
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset)...)
				// load k sets
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySets)...)
				// load k
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsK)...)
				code.Codes[code.CodeLength] = cg.OP_aaload
				code.CodeLength++
				primitiveObjectConverter.getFromObject(class, code, s.StatmentForRangeAttr.Expression.VariableType.Map.K)
				closure.storeLocalCloureVar(class, code, s.StatmentForRangeAttr.IdentifierV.Var)
			} else {
				// load k sets
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySets)...)
				// load k
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeMap.KeySetsK)...)
				code.Codes[code.CodeLength] = cg.OP_aaload
				code.CodeLength++
				primitiveObjectConverter.getFromObject(class, code, s.StatmentForRangeAttr.Expression.VariableType.Map.K)
				copyOP(code, storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.K.Typ, s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset)...)
			}
		}
	} else { // for k,v  = range xxx
		// store v
		stack, remainStack, op, _, classname, name, descriptor := m.MakeExpression.getLeftValue(class, code, s.StatmentForRangeAttr.ExpressionV, context)
		if stack > maxstack { // this means  current stack is 0
			maxstack = stack
		}
		copyOP(code,
			loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.Map.V.Typ, s.StatmentForRangeAttr.AutoVarForRangeMap.V)...)
		if t := remainStack + s.StatmentForRangeAttr.Expression.VariableType.Map.V.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		copyOPLeftValue(class, code, op, classname, name, descriptor)
	}
	// build block
	m.buildBlock(class, code, s.Block, context)
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

	// pop 3
	code.Codes[code.CodeLength] = cg.OP_pop
	code.Codes[code.CodeLength+1] = cg.OP_pop
	code.Codes[code.CodeLength+2] = cg.OP_pop
	code.CodeLength += 3
	return
}

func (m *MakeClass) buildForRangeStatementForArray(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {

	//build array expression
	maxstack, _ = m.MakeExpression.build(class, code, s.StatmentForRangeAttr.Expression, context) // array on stack
	// if null skip
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6) // goto pop
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7) //goto for
	code.Codes[code.CodeLength+7] = cg.OP_pop
	code.CodeLength += 8
	s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))

	if s.StatmentForRangeAttr.Expression.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY {
		//get elements
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxstack {
			maxstack = 2
		}
		meta := ArrayMetas[s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ]
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.classname,
			Field:      "elements",
			Descriptor: meta.elementsFieldDescriptor,
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_JAVA_ARRAY, s.StatmentForRangeAttr.AutoVarForRangeArray.Elements)...)
		//get start
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.classname,
			Field:      "start",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.Start)...)
		//get end
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.classname,
			Field:      "end",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.End)...)
	} else { // java_array
		//get length
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxstack {
			maxstack = 2
		}
		code.Codes[code.CodeLength+1] = cg.OP_anewarray
		code.CodeLength += 2
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.End)...)
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_JAVA_ARRAY, s.StatmentForRangeAttr.AutoVarForRangeArray.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_JAVA_ARRAY, s.StatmentForRangeAttr.AutoVarForRangeArray.Elements)...)
	}

	var koffset uint16
	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	if s.StatmentForRangeAttr.ModelKV && s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		koffset = s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset
	} else {
		koffset = s.StatmentForRangeAttr.AutoVarForRangeArray.K
	}
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, koffset)...)

	loopbeginAt := code.CodeLength
	// load start
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.Start)...)
	// load k
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, koffset)...)
	// mk index
	code.Codes[code.CodeLength] = cg.OP_iadd
	code.Codes[code.CodeLength+1] = cg.OP_dup
	code.CodeLength += 2
	// load end
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.End)...)
	if 3 > maxstack {
		maxstack = 3
	}
	/*
		k + start >= end,break loop,pop index on stack
		check if need to break
	*/
	rangeend := (&cg.JumpBackPatch{}).FromCode(cg.OP_if_icmpge, code)
	//load elements
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRangeArray.Elements)...)
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	// load value
	switch s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_baload
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_saload
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
	}
	code.CodeLength++
	// before store to local v ,cast into real type
	if s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ == ast.VARIABLE_TYPE_STRING {
	} else if s.StatmentForRangeAttr.Expression.VariableType.ArrayType.IsPointer() {
		primitiveObjectConverter.castPointerTypeToRealType(class, code, s.StatmentForRangeAttr.Expression.VariableType.ArrayType)
	}

	//store to v tmp
	copyOP(code,
		storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ,
			s.StatmentForRangeAttr.AutoVarForRangeArray.V)...)

	//current stack is 0
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.StatmentForRangeAttr.IdentifierV.Var.BeenCaptured {
			panic("...")
		}
		copyOP(code,
			loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.AutoVarForRangeArray.V)...)
		copyOP(code,
			storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
		if s.StatmentForRangeAttr.ModelKV {
			if s.StatmentForRangeAttr.IdentifierK.Var.BeenCaptured {
				panic("...")
			}
			copyOP(code,
				loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.AutoVarForRangeArray.K)...)
			copyOP(code,
				storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset)...)
		}
	} else { // for k,v = range arr
		// store v
		//get ops,make ops ready
		stack, remainStack, ops, target, classname, name, descriptor := m.MakeExpression.getLeftValue(class,
			code, s.StatmentForRangeAttr.ExpressionV, context)
		if stack > maxstack {
			maxstack = stack
		}
		//load v
		copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ,
			s.StatmentForRangeAttr.AutoVarForRangeArray.V)...)
		if t := remainStack + s.StatmentForRangeAttr.Expression.VariableType.ArrayType.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		//convert to suitable type
		if target.IsInteger() && target.Typ != s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ {
			m.MakeExpression.numberTypeConverter(code, s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, target.Typ)
		}
		if t := remainStack + target.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		copyOPLeftValue(class, code, ops, classname, name, descriptor)
		if s.StatmentForRangeAttr.ModelKV { // set to k
			stack, remainStack, ops, target, classname, name, descriptor := m.MakeExpression.getLeftValue(class,
				code, s.StatmentForRangeAttr.ExpressionK, context)
			if stack > maxstack {
				maxstack = stack
			}
			if t := remainStack + 1; t > maxstack {
				maxstack = t
			}
			// load k
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, koffset)...)
			m.MakeExpression.numberTypeConverter(code, ast.VARIABLE_TYPE_INT, target.Typ)
			if t := target.JvmSlotSize() + remainStack; t > maxstack {
				maxstack = t
			}
			copyOPLeftValue(class, code, ops, classname, name, descriptor)
		}
	}

	// build block
	m.buildBlock(class, code, s.Block, context)
	//innc k
	s.ContinueOPOffset = code.CodeLength
	code.Codes[code.CodeLength] = cg.OP_iinc
	if koffset > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(koffset)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	//goto begin
	jumpto(cg.OP_goto, code, loopbeginAt)
	backPatchEs([]*cg.JumpBackPatch{rangeend}, code.CodeLength) // jump to here
	//pop index on stack
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	return
}

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	if s.StatmentForRangeAttr != nil {
		if s.StatmentForRangeAttr.Expression.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY ||
			s.StatmentForRangeAttr.Expression.VariableType.Typ == ast.VARIABLE_TYPE_JAVA_ARRAY {
			return m.buildForRangeStatementForArray(class, code, s, context)
		} else { // for map
			return m.buildForRangeStatementForMap(class, code, s, context)
		}
	}
	//init
	if s.Init != nil {
		stack, _ := m.MakeExpression.build(class, code, s.Init, context)
		if stack > maxstack {
			maxstack = stack
		}
	}
	loopBeginAt := code.CodeLength
	s.ContinueOPOffset = loopBeginAt
	//condition
	if s.Condition != nil {
		stack, es := m.MakeExpression.build(class, code, s.Condition, context)
		backPatchEs(es, code.CodeLength)
		if stack > maxstack {
			maxstack = stack
		}
		s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code))
	} else {

	}
	m.buildBlock(class, code, s.Block, context)
	if s.Post != nil {
		s.ContinueOPOffset = code.CodeLength
		stack, _ := m.MakeExpression.build(class, code, s.Post, context)
		if stack > maxstack {
			maxstack = stack
		}
	}
	jumpto(cg.OP_goto, code, loopBeginAt)
	return
}
