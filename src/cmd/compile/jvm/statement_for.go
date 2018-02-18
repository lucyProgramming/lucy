package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildForRangeStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	//build array expression
	maxstack, _ = m.MakeExpression.build(class, code, s.StatmentForRangeAttr.AarrayExpression, context) // array on stack
	meta := ArrayMetas[s.StatmentForRangeAttr.AarrayExpression.VariableType.CombinationType.Typ]
	code.Codes[code.CodeLength] = cg.OP_dup //dup top
	if 2 > maxstack {
		maxstack = 2
	}
	//get elements
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
	if s.Condition.Typ == ast.EXPRESSION_TYPE_COLON_ASSIGN {
		// k set to 0
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
		var koffset uint16
		if s.StatmentForRangeAttr.ModelKV {
			koffset = s.StatmentForRangeAttr.IdentifierK.Var.LocalValOffset
		} else {
			koffset = s.StatmentForRangeAttr.AutoVarForRange.K
		}
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, koffset)...)
		loopbegin := code.CodeLength
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
		s.BackPatchs = append(s.BackPatchs, (*&cg.JumpBackPatch{}).FromCode(cg.OP_if_icmpge, code)) // k + start >= end,break loop
		//load elements
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, s.StatmentForRangeAttr.AutoVarForRange.Elements)...)
		code.Codes[code.CodeLength] = cg.OP_swap
		code.CodeLength++
		// load value
		switch s.StatmentForRangeAttr.AarrayExpression.VariableType.CombinationType.Typ {
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
			meta := ArrayMetas[s.StatmentForRangeAttr.AarrayExpression.VariableType.CombinationType.CombinationType.Typ]
			code.Codes[code.CodeLength] = cg.OP_aaload // cast into real type
			code.Codes[code.CodeLength+1] = cg.OP_checkcast
			class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+2:code.CodeLength+4])
			code.CodeLength += 4
		}
		// store to v
		copyOP(code, storeSimpleVarOp(s.StatmentForRangeAttr.AarrayExpression.VariableType.CombinationType.Typ, s.StatmentForRangeAttr.IdentifierV.Var.LocalValOffset)...)
		// build block
		m.buildBlock(class, code, s.Block, context)
		//innc k
		code.Codes[code.CodeLength] = cg.OP_iinc
		if koffset > 255 {
			panic("over 255")
		}
		code.Codes[code.CodeLength+1] = byte(koffset)
		code.Codes[code.CodeLength+2] = 1
		code.CodeLength += 3
		//goto begin
		jumpto(cg.OP_goto, code, loopbegin)
		return
	}

	panic(111)

	return
}

func (m *MakeClass) buildForStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementFor, context *Context) (maxstack uint16) {
	if s.StatmentForRangeAttr != nil {
		return m.buildForRangeStatement(class, code, s, context)
	}
	//init
	if s.Init != nil {
		//	code.MKLineNumber(s.Init.Pos.StartLine)
		stack, _ := m.MakeExpression.build(class, code, s.Init, context)
		if stack > maxstack {
			maxstack = stack
		}
	}
	s.LoopBegin = code.CodeLength
	s.ContinueOPOffset = s.LoopBegin
	//condition
	if s.Condition != nil {
		//code.MKLineNumber(s.Condition.Pos.StartLine)
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
		//	code.MKLineNumber(s.Post.Pos.StartLine)
		s.ContinueOPOffset = code.CodeLength
		stack, _ := m.MakeExpression.build(class, code, s.Post, context)
		if stack > maxstack {
			maxstack = stack
		}
	}
	jumpto(cg.OP_goto, code, s.LoopBegin)
	return
}
