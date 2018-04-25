package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type AutoVarForRangeArray struct {
	Elements      uint16
	Start, End, K uint16
	V             uint16
}

func (m *MakeClass) buildForRangeStatementForArray(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context, state *StackMapState) (maxstack uint16) {
	//build array expression
	maxstack, _ = m.MakeExpression.build(class, code, s.RangeAttr.Expression, context, state) // array on stack

	// if null skip
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6) // goto pop
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7) //goto for
	code.Codes[code.CodeLength+7] = cg.OP_pop
	state.Stacks = append(state.Stacks,
		state.newStackMapVerificationTypeInfo(class, s.RangeAttr.Expression.Value))
	context.MakeStackMap(code, state, code.CodeLength+7)
	context.MakeStackMap(code, state, code.CodeLength+11)
	state.popStack(1)
	code.CodeLength += 8
	s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))
	forState := (&StackMapState{}).FromLast(state)
	defer func() {
		state.addTop(forState)
	}()
	var autoVar AutoVarForRangeArray
	{
		// eles
		if s.RangeAttr.Expression.Value.Typ == ast.VARIABLE_TYPE_ARRAY {
			meta := ArrayMetas[s.RangeAttr.Expression.Value.ArrayType.Typ]
			_, t, _ := Descriptor.ParseType([]byte(meta.elementsFieldDescriptor))
			autoVar.Elements = code.MaxLocals
			code.MaxLocals++
			forState.appendLocals(class, t)
		} else {
			autoVar.Elements = code.MaxLocals
			code.MaxLocals++
			forState.appendLocals(class, s.RangeAttr.Expression.Value)
		}
		// start
		autoVar.Start = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
		//end
		autoVar.End = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
		// K
		autoVar.K = code.MaxLocals
		code.MaxLocals++
		forState.appendLocals(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
	}

	if s.RangeAttr.Expression.Value.Typ == ast.VARIABLE_TYPE_ARRAY {
		//get elements
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxstack {
			maxstack = 2
		}
		meta := ArrayMetas[s.RangeAttr.Expression.Value.ArrayType.Typ]
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.classname,
			Field:      "elements",
			Descriptor: meta.elementsFieldDescriptor,
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_JAVA_ARRAY, autoVar.Elements)...)
		//get start
		code.Codes[code.CodeLength] = cg.OP_dup
		code.Codes[code.CodeLength+1] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.classname,
			Field:      "start",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.Start)...)
		//get end
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.classname,
			Field:      "end",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.End)...)
	} else { // java_array
		//get length
		code.Codes[code.CodeLength] = cg.OP_dup //dup top
		if 2 > maxstack {
			maxstack = 2
		}
		code.Codes[code.CodeLength+1] = cg.OP_arraylength
		code.CodeLength += 2
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.End)...)
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_JAVA_ARRAY, autoVar.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.Start)...)

	}

	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.K)...)
	loopbeginAt := code.CodeLength
	context.MakeStackMap(code, forState, loopbeginAt)
	// load start
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.Start)...)
	// load k
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.K)...)
	// mk index
	code.Codes[code.CodeLength] = cg.OP_iadd
	code.Codes[code.CodeLength+1] = cg.OP_dup
	code.CodeLength += 2
	// load end
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.End)...)
	if 3 > maxstack {
		maxstack = 3
	}
	/*
		k + start >= end,break loop,pop index on stack
		check if need to break
	*/
	rangeend := (&cg.JumpBackPatch{}).FromCode(cg.OP_if_icmpge, code)
	//load elements
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, autoVar.Elements)...)
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	// load value
	switch s.RangeAttr.Expression.Value.ArrayType.Typ {
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
	}
	code.CodeLength++
	// before store to local v ,cast into real type
	if s.RangeAttr.Expression.Value.ArrayType.Typ == ast.VARIABLE_TYPE_STRING {
	} else if s.RangeAttr.Expression.Value.ArrayType.IsPointer() {
		primitiveObjectConverter.castPointerTypeToRealType(class, code, s.RangeAttr.Expression.Value.ArrayType)
	}
	// v
	autoVar.V = code.MaxLocals
	code.MaxLocals += jvmSize(s.RangeAttr.IdentifierV.Var.Typ)
	//store to v tmp
	copyOP(code,
		storeSimpleVarOp(s.RangeAttr.Expression.Value.ArrayType.Typ,
			autoVar.V)...)

	forState.appendLocals(class, s.RangeAttr.IdentifierV.Var.Typ)
	//current stack is 0
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		if s.RangeAttr.IdentifierV.Var.BeenCaptured {
			panic(11)
		} else {

			copyOP(code,
				loadSimpleVarOp(s.RangeAttr.Expression.Value.ArrayType.Typ, autoVar.V)...)
			s.RangeAttr.IdentifierV.Var.LocalValOffset = code.MaxLocals
			code.MaxLocals += jvmSize(s.RangeAttr.IdentifierV.Var.Typ)
			copyOP(code,
				storeSimpleVarOp(s.RangeAttr.Expression.Value.ArrayType.Typ, s.RangeAttr.IdentifierV.Var.LocalValOffset)...)
			forState.appendLocals(class, s.RangeAttr.IdentifierV.Var.Typ)
		}
		if s.RangeAttr.ModelKV {
			if s.RangeAttr.IdentifierK.Var.BeenCaptured {
				panic(11)
			} else {
				copyOP(code,
					loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.K)...)
				s.RangeAttr.IdentifierK.Var.LocalValOffset = code.MaxLocals
				code.MaxLocals++
				copyOP(code,
					storeSimpleVarOp(s.RangeAttr.Expression.Value.ArrayType.Typ,
						s.RangeAttr.IdentifierK.Var.LocalValOffset)...)
				forState.appendLocals(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
			}
		}
	} else { // for k,v = range arr
		// store v
		//get ops,make ops ready
		stack, remainStack, ops, target, classname, name, descriptor := m.MakeExpression.getLeftValue(class,
			code, s.RangeAttr.ExpressionV, context, forState)
		if stack > maxstack {
			maxstack = stack
		}
		//load v
		copyOP(code, loadSimpleVarOp(s.RangeAttr.Expression.Value.ArrayType.Typ,
			autoVar.V)...)
		if t := remainStack + jvmSize(s.RangeAttr.Expression.Value.ArrayType); t > maxstack {
			maxstack = t
		}
		//convert to suitable type
		if target.IsInteger() && target.Typ != s.RangeAttr.Expression.Value.ArrayType.Typ {
			m.MakeExpression.numberTypeConverter(code, s.RangeAttr.Expression.Value.ArrayType.Typ, target.Typ)
		}
		if t := remainStack + jvmSize(target); t > maxstack {
			maxstack = t
		}
		copyOPLeftValue(class, code, ops, classname, name, descriptor)
		if s.RangeAttr.ModelKV { // set to k
			stack, remainStack, ops, target, classname, name, descriptor := m.MakeExpression.getLeftValue(class,
				code, s.RangeAttr.ExpressionK, context, forState)
			if stack > maxstack {
				maxstack = stack
			}
			if t := remainStack + 1; t > maxstack {
				maxstack = t
			}
			// load k
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, autoVar.K)...)
			m.MakeExpression.numberTypeConverter(code, ast.VARIABLE_TYPE_INT, target.Typ)
			if t := jvmSize(target) + remainStack; t > maxstack {
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
	if autoVar.K > 255 {
		panic("over 255")
	}
	code.Codes[code.CodeLength+1] = byte(autoVar.K)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	//goto begin
	jumpto(cg.OP_goto, code, loopbeginAt)
	backPatchEs([]*cg.JumpBackPatch{rangeend}, code.CodeLength) // jump to here
	//pop index on stack

	state.Stacks = append(state.Stacks,
		state.newStackMapVerificationTypeInfo(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT}))
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	code.Codes[code.CodeLength] = cg.OP_pop
	code.CodeLength++
	return
}
