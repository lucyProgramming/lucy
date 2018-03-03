package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForRangeStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	//build array expression
	maxstack, _ = m.MakeExpression.build(class, code, s.StatmentForRangeAttr.AarrayExpression, context) // array on stack

	// if null skip
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	code.Codes[code.CodeLength+1] = cg.OP_ifnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+2:code.CodeLength+4], 6) // goto pop
	code.Codes[code.CodeLength+4] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 7) //goto for
	code.Codes[code.CodeLength+7] = cg.OP_pop
	code.CodeLength += 8
	s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code))

	//get elements
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	if 2 > maxstack {
		maxstack = 2
	}
	meta := ArrayMetas[s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.Typ]
	code.Codes[code.CodeLength+1] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      meta.classname,
		Name:       "elements",
		Descriptor: meta.elementsFieldDescriptor,
	}, code.Codes[code.CodeLength+2:code.CodeLength+4])
	code.CodeLength += 4
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRange.Elements)...)
	//get start
	code.Codes[code.CodeLength] = cg.OP_dup
	code.Codes[code.CodeLength+1] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      meta.classname,
		Name:       "start",
		Descriptor: "I",
	}, code.Codes[code.CodeLength+2:code.CodeLength+4])
	code.CodeLength += 4
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRange.Start)...)
	//get end
	code.Codes[code.CodeLength] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      meta.classname,
		Name:       "end",
		Descriptor: "I",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRange.End)...)
	var koffset uint16
	// k set to 0
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	if s.StatmentForRangeAttr.ModelKV && s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		koffset = s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset
	} else {
		koffset = s.StatmentForRangeAttr.AutoVarForRange.K
	}
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, koffset)...)
	loopbeginAt := code.CodeLength
	// load start
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRange.Start)...)
	// load k
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, koffset)...)
	// mk index
	code.Codes[code.CodeLength] = cg.OP_iadd
	code.Codes[code.CodeLength+1] = cg.OP_dup
	code.CodeLength += 2
	//check if need to break
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, s.StatmentForRangeAttr.AutoVarForRange.End)...)
	if 3 > maxstack {
		maxstack = 3
	}
	/*
		k + start >= end,break loop,pop index on stack
	*/
	rangeend := (&cg.JumpBackPatch{}).FromCode(cg.OP_if_icmpge, code)
	//load elements
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRange.Elements)...)
	code.Codes[code.CodeLength] = cg.OP_swap
	code.CodeLength++
	// load value
	switch s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		code.Codes[code.CodeLength] = cg.OP_baload
		code.CodeLength++
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_saload
		code.CodeLength++
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_iaload
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_laload
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_faload
		code.CodeLength++
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_daload
		code.CodeLength++
	case ast.VARIABLE_TYPE_STRING:
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
	case ast.VARIABLE_TYPE_OBJECT:
		code.Codes[code.CodeLength] = cg.OP_aaload
		code.CodeLength++
	case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
		meta := ArrayMetas[s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.ArrayType.Typ]
		code.Codes[code.CodeLength] = cg.OP_aaload // cast into real type
		code.Codes[code.CodeLength+1] = cg.OP_checkcast
		class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+2:code.CodeLength+4])
		code.CodeLength += 4
	}
	// store to v
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		copyOP(code,
			storeSimpleVarOp(s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.Typ,
				s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
	} else {
		copyOP(code,
			storeSimpleVarOp(s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.Typ,
				s.StatmentForRangeAttr.AutoVarForRange.V)...)
	}
	//current stack is 0
	if s.Condition.Typ == ast.EXPRESSION_TYPE_ASSIGN {
		//get ops,make ops ready
		var vExpression *ast.Expression
		if s.StatmentForRangeAttr.ModelKV {
			vExpression = s.StatmentForRangeAttr.Lefts[1]
		} else {
			vExpression = s.StatmentForRangeAttr.Lefts[0]
		}
		stack, remainStack, ops, target, classname, name, descriptor := m.MakeExpression.getLeftValue(class,
			code, vExpression, context)
		if stack > maxstack {
			maxstack = stack
		}
		//load v
		copyOP(code, loadSimpleVarOp(s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.Typ,
			s.StatmentForRangeAttr.AutoVarForRange.V)...)
		if t := remainStack + s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		//convert to suitable type
		if target.IsInteger() && target.Typ != s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.Typ {
			m.MakeExpression.numberTypeConverter(code, s.StatmentForRangeAttr.AarrayExpression.VariableType.ArrayType.Typ, target.Typ)
		}
		if t := remainStack + target.JvmSlotSize(); t > maxstack {
			maxstack = t
		}
		copyOPLeftValue(class, code, ops, classname, name, descriptor)
		if s.StatmentForRangeAttr.ModelKV { // set to k
			stack, remainStack, ops, target, classname, name, descriptor := m.MakeExpression.getLeftValue(class,
				code, vExpression, context)
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
		return m.buildForRangeStatement(class, code, s, context)
	}
	//init
	if s.Init != nil {
		stack, _ := m.MakeExpression.build(class, code, s.Init, context)
		if stack > maxstack {
			maxstack = stack
		}
	}
	s.LoopBegin = code.CodeLength
	s.ContinueOPOffset = s.LoopBegin
	//condition
	if s.Condition != nil {
		stack, es := m.MakeExpression.build(class, code, s.Condition, context)
		backPatchEs(es, code.CodeLength)
		if stack > maxstack {
			maxstack = stack
		}
		code.Codes[code.CodeLength] = cg.OP_ifeq
		b := cg.JumpBackPatch{}
		b.CurrentCodeLength = code.CodeLength
		b.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
		s.BackPatchs = append(s.BackPatchs, &b)
		code.CodeLength += 3
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
	jumpto(cg.OP_goto, code, s.LoopBegin)
	return
}
