package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildSwitchStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementSwitch, context *Context, state *StackMapState) (maxstack uint16) {
	// if equal,leave 0 on stack
	compare := func(t *ast.VariableType) {
		state.popStack(2)
		switch t.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_ENUM:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_isub
			code.CodeLength++
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lcmp
			code.CodeLength++
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_fcmpg
			code.CodeLength++
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dcmpg
			code.CodeLength++
		case ast.VARIABLE_TYPE_STRING:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      java_string_class,
				Method:     "compareTo",
				Descriptor: "(Ljava/lang/String;)I",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_MAP:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY:
			context.MakeStackMap(code, state, code.CodeLength+7)
			state.pushStack(class, &ast.VariableType{
				Typ: ast.VARIABLE_TYPE_BOOL,
			})
			context.MakeStackMap(code, state, code.CodeLength+8)
			state.popStack(1)
			code.Codes[code.CodeLength] = cg.OP_if_acmpeq
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			code.CodeLength += 8
		}
	}
	maxstack, _ = m.MakeExpression.build(class, code, s.Condition, context, state)
	//value is on stack
	var exit *cg.Exit
	size := jvmSize(s.Condition.Value)
	currentStack := size
	state.pushStack(class, s.Condition.Value)
	for _, c := range s.StatmentSwitchCases {
		if exit != nil {
			backfillExit([]*cg.Exit{exit}, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
		mathes := []*cg.Exit{}
		for _, ee := range c.Matches {
			if ee.MayHaveMultiValue() && len(ee.Values) > 1 {
				stack, _ := m.MakeExpression.build(class, code, ee, context, state)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				multiValuePacker.storeArrayListAutoVar(code, context)
				for kkk, ttt := range ee.Values {
					currentStack = size
					if size == 1 {
						code.Codes[code.CodeLength] = cg.OP_dup
					} else {
						code.Codes[code.CodeLength] = cg.OP_dup2
					}
					code.CodeLength++
					state.pushStack(class, s.Condition.Value)
					currentStack += size
					if currentStack > maxstack {
						maxstack = currentStack
					}
					stack = multiValuePacker.unPack(class, code, kkk, ttt, context)
					if t := stack + currentStack; t > maxstack {
						maxstack = t
					}
					state.pushStack(class, s.Condition.Value)
					compare(s.Condition.Value)
					// consume result on stack
					mathes = append(mathes, (&cg.Exit{}).FromCode(cg.OP_ifeq, code))
				}
				continue
			}
			currentStack = size
			// mk stack ready
			if size == 1 {
				code.Codes[code.CodeLength] = cg.OP_dup
			} else {
				code.Codes[code.CodeLength] = cg.OP_dup2
			}
			code.CodeLength++
			currentStack += size
			state.pushStack(class, s.Condition.Value)
			stack, _ := m.MakeExpression.build(class, code, ee, context, state)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			state.pushStack(class, s.Condition.Value)
			compare(s.Condition.Value)

			mathes = append(mathes, (&cg.Exit{}).FromCode(cg.OP_ifeq, code)) // comsume result on stack
		}
		// should be goto next,here is no match
		exit = (&cg.Exit{}).FromCode(cg.OP_goto, code)
		// if match goto here
		backfillExit(mathes, code.CodeLength)
		//before block,pop off stack
		context.MakeStackMap(code, state, code.CodeLength)
		if size == 1 {
			code.Codes[code.CodeLength] = cg.OP_pop
		} else {
			code.Codes[code.CodeLength] = cg.OP_pop2
		}
		code.CodeLength++
		//block is here
		if c.Block != nil {
			ss := (&StackMapState{}).FromLast(state)
			m.buildBlock(class, code, c.Block, context, ss)
			state.addTop(ss)
		}
		if c.Block == nil || c.Block.DeadEnding == false {
			s.BackPatchs = append(s.BackPatchs,
				(&cg.Exit{}).FromCode(cg.OP_goto, code)) // matched,goto switch outside
		}
	}
	backfillExit([]*cg.Exit{exit}, code.CodeLength)
	context.MakeStackMap(code, state, code.CodeLength)
	if size == 1 {
		code.Codes[code.CodeLength] = cg.OP_pop
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop2
	}
	code.CodeLength++
	state.popStack(1)
	// build default
	if s.Default != nil {
		var ss *StackMapState
		if s.Default.HaveVariableDefinition() {
			ss = (&StackMapState{}).FromLast(state)
		} else {
			ss = state
		}
		m.buildBlock(class, code, s.Default, context, ss)
		state.addTop(ss)
	}
	return
}
