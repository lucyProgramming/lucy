package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildSwitchStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementSwitch, context *Context, state *StackMapState) (maxstack uint16) {
	// if equal,leave 0 on stack
	switchCompare := func(t *ast.VariableType) {
		switch t.Typ {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
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
			code.Codes[code.CodeLength] = cg.OP_if_acmpeq
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			code.CodeLength += 8
		}
	}
	maxstack, _ = m.MakeExpression.build(class, code, s.Condition, context, nil)
	//value is on stack
	var exit *cg.JumpBackPatch
	size := s.Condition.VariableType.JvmSlotSize()
	currentStack := size
	for k, c := range s.StatmentSwitchCases {
		if exit != nil {
			backPatchEs([]*cg.JumpBackPatch{exit}, code.CodeLength)
		}
		gotoBodyExits := []*cg.JumpBackPatch{}
		needPop := false
		for kk, ee := range c.Matches {
			if ee.MayHaveMultiValue() && len(ee.VariableTypes) > 0 {
				stack, _ := m.MakeExpression.build(class, code, ee, context, nil)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				m.MakeExpression.buildStoreArrayListAutoVar(code, context)
				for kkk, ttt := range ee.VariableTypes {
					currentStack = size
					if k == len(s.StatmentSwitchCases)-1 && kk == len(c.Matches)-1 && kkk == len(ee.VariableTypes)-1 {
					} else {
						if size == 1 {
							code.Codes[code.CodeLength] = cg.OP_dup
						} else {
							code.Codes[code.CodeLength] = cg.OP_dup2
						}
						code.CodeLength++
						currentStack += size
						needPop = true
						if currentStack > maxstack {
							maxstack = currentStack
						}
					}
					stack = m.MakeExpression.unPackArraylist(class, code, kkk, ttt, context)
					if t := stack + currentStack; t > maxstack {
						maxstack = t
					}
					switchCompare(s.Condition.VariableType)
					gotoBodyExits = append(gotoBodyExits, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code)) // comsume result on stack
				}
				continue
			}
			currentStack = size
			// mk stack ready
			if k == len(s.StatmentSwitchCases)-1 && kk == len(c.Matches)-1 { // last one
			} else { // not the last one,dup the stack
				if size == 1 {
					code.Codes[code.CodeLength] = cg.OP_dup
				} else {
					code.Codes[code.CodeLength] = cg.OP_dup2
				}
				code.CodeLength++
				needPop = true
				currentStack += size
				if currentStack > maxstack {
					maxstack = currentStack
				}
			}
			stack, _ := m.MakeExpression.build(class, code, ee, context, nil)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			switchCompare(s.Condition.VariableType)
			gotoBodyExits = append(gotoBodyExits, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code)) // comsume result on stack
		}
		// should be goto next,here is no match
		if k != len(s.StatmentSwitchCases)-1 || s.Default != nil {
			exit = (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)
		}
		// if match goto here
		backPatchEs(gotoBodyExits, code.CodeLength)
		//before block,pop off stack
		if needPop {
			if size == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop2
			}
			code.CodeLength++
		}
		//block is here
		if c.Block != nil {
			m.buildBlock(class, code, c.Block, context, state)
		}
		if k != len(s.StatmentSwitchCases)-1 || s.Default != nil {
			s.BackPatchs = append(s.BackPatchs, (&cg.JumpBackPatch{}).FromCode(cg.OP_goto, code)) // matched,goto switch outside
		}
	}
	// build default
	if s.Default != nil {
		backPatchEs([]*cg.JumpBackPatch{exit}, code.CodeLength)
		m.buildBlock(class, code, s.Default, context, state)
	}
	return
}
