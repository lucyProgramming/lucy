package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForRangeStatementForArray(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context, state *StackMapState) (maxstack uint16) {
	//build array expression
	maxstack, _ = m.MakeExpression.build(class, code, s.StatmentForRangeAttr.Expression, context, state) // array on stack

	// if null skip
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6) // goto pop
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7) //goto for
	code.Codes[code.CodeLength+7] = cg.OP_pop
	state.Stacks = append(state.Stacks,
		state.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.Expression.VariableType)...)
	code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(state, code.CodeLength+7))
	code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(state, code.CodeLength+11))
	state.popStack(1)
	code.CodeLength += 8
	s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	forState := (&StackMapState{}).FromLast(state)
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
		code.Codes[code.CodeLength+1] = cg.OP_arraylength
		code.CodeLength += 2
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.End)...)
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_JAVA_ARRAY, s.StatmentForRangeAttr.AutoVarForRangeArray.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_JAVA_ARRAY, s.StatmentForRangeAttr.AutoVarForRangeArray.Start)...)

	}

	{
		// eles
		if s.StatmentForRangeAttr.Expression.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY {
			meta := ArrayMetas[s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ]
			_, t, _ := Descriptor.ParseType([]byte(meta.elementsFieldDescriptor))
			forState.Locals = append(forState.Locals, forState.newStackMapVerificationTypeInfo(class, t)...)
		} else {
			forState.Locals = append(forState.Locals, forState.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.Expression.VariableType)...)
		}
		// start
		forState.Locals = append(forState.Locals,
			forState.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
		//end
		forState.Locals = append(forState.Locals,
			forState.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
		// k
		forState.Locals = append(forState.Locals,
			forState.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)

	}

	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.K)...)
	loopbeginAt := code.CodeLength
	code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(forState, loopbeginAt))
	// load start
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.Start)...)
	// load k
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.K)...)
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
			stack := closure.createCloureVar(class, code, s.StatmentForRangeAttr.IdentifierV.Var)
			if stack > maxstack {
				maxstack = stack
			}
			copyOP(code,
				loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
			copyOP(code,
				loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.AutoVarForRangeArray.V)...)
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
			copyOP(code,
				loadSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.AutoVarForRangeArray.V)...)
			copyOP(code,
				storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
			forState.Locals = append(forState.Locals,
				forState.newStackMapVerificationTypeInfo(class, s.StatmentForRangeAttr.IdentifierV.Var.Typ)...)
		}
		if s.StatmentForRangeAttr.ModelKV {
			if s.StatmentForRangeAttr.IdentifierK.Var.BeenCaptured {
				stack := closure.createCloureVar(class, code, s.StatmentForRangeAttr.IdentifierK.Var)
				if stack > maxstack {
					maxstack = stack
				}
				copyOP(code,
					loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.K)...)
				closure.storeLocalCloureVar(class, code, s.StatmentForRangeAttr.IdentifierK.Var)
				{
					t := &ast.VariableType{}
					t.Typ = ast.VARIABLE_TYPE_OBJECT
					t.Class = &ast.Class{}
					t.Class.Name = closure.getMeta(ast.VARIABLE_TYPE_INT).className
					forState.Locals = append(forState.Locals,
						forState.newStackMapVerificationTypeInfo(class, t)...)
				}
			} else {
				copyOP(code,
					loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.K)...)
				copyOP(code,
					storeSimpleVarOp(s.StatmentForRangeAttr.Expression.VariableType.ArrayType.Typ, s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset)...)
				forState.Locals = append(forState.Locals,
					forState.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
			}
		}
	} else { // for k,v = range arr
		// store v
		//get ops,make ops ready
		stack, remainStack, ops, target, classname, name, descriptor := m.MakeExpression.getLeftValue(class,
			code, s.StatmentForRangeAttr.ExpressionV, context, state)
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
				code, s.StatmentForRangeAttr.ExpressionK, context, state)
			if stack > maxstack {
				maxstack = stack
			}
			if t := remainStack + 1; t > maxstack {
				maxstack = t
			}
			// load k
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRangeArray.K)...)
			m.MakeExpression.numberTypeConverter(code, ast.VARIABLE_TYPE_INT, target.Typ)
			if t := target.JvmSlotSize() + remainStack; t > maxstack {
				maxstack = t
			}
			copyOPLeftValue(class, code, ops, classname, name, descriptor)
		}
	}

	// build block
	m.buildBlock(class, code, s.Block, context, forState)
	//innc k
	s.ContinueOPOffset = code.CodeLength
	code.Codes[code.CodeLength] = cg.OP_iinc
	if s.StatmentForRangeAttr.AutoVarForRangeArray.K > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(s.StatmentForRangeAttr.AutoVarForRangeArray.K)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	//goto begin
	jumpto(cg.OP_goto, code, loopbeginAt)
	backPatchEs([]*cg.JumpBackPatch{rangeend}, code.CodeLength) // jump to here
	//pop index on stack

	state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})...)
	code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, context.MakeStackMap(state, code.CodeLength))
	state.popStack(1)
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	return
}
