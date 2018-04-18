package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildColonAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	vs := e.Data.(*ast.ExpressionDeclareVariable)
	//first round
	index := len(vs.Vs) - 1
	currentStack := uint16(0)
	for index >= 0 {
		//this variable not been captured,also not declared here
		if vs.Vs[index].BeenCaptured {
			if vs.IfDeclareBefor[index] == false {
				stack := closure.createCloureVar(class, code, vs.Vs[index])
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
			}
			// load to stack
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, vs.Vs[index].LocalValOffset)...)
			state.Stacks = append(state.Stacks,
				state.newStackMapVerificationTypeInfo(class, state.newObjectVariableType(closure.getMeta(vs.Vs[index].Typ.Typ).className))...)
			currentStack += 1
		}
		index--
	}
	variables := vs.Vs
	for _, v := range vs.Values {
		if v.MayHaveMultiValue() && len(v.VariableTypes) > 1 {
			stack, _ := m.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			m.buildStoreArrayListAutoVar(code, context)
			for kk, tt := range v.VariableTypes {
				if variables[0].Name == ast.NO_NAME_IDENTIFIER {
					continue
				}
				stack = m.unPackArraylist(class, code, kk, tt, context)
				if t := stack + currentStack; t > maxstack {
					maxstack = t
				}
				if tt.IsNumber() && tt.Typ != variables[0].Typ.Typ {
					m.numberTypeConverter(code, tt.Typ, variables[0].Typ.Typ)
				}
				if variables[0].IsGlobal {
					storeGlobalVar(class, m.MakeClass.mainclass, code, variables[0])
				} else {
					m.MakeClass.storeLocalVar(class, code, variables[0])
					if variables[0].BeenCaptured {
						currentStack -= 1
					}
					m.MakeClass.appendLocalVar(class, code, variables[0], state)
				}
				variables = variables[1:]
			}
			continue
		}
		variableType := v.VariableType
		if v.MayHaveMultiValue() {
			variableType = v.VariableTypes[0]
		}
		stack, es := m.build(class, code, v, context, state)
		state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, v.VariableType)...)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength) //
			code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
				context.MakeStackMap(state, code.CodeLength))
		}
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		if variables[0].Name == ast.NO_NAME_IDENTIFIER {
			if variableType.JvmSlotSize() == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
				code.CodeLength++
				state.popStack(1)
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop2
				code.CodeLength++
				state.popStack(2)
			}
			variables = variables[1:]
			continue
		}
		if variableType.IsNumber() && variableType.Typ != variables[0].Typ.Typ {
			m.numberTypeConverter(code, variableType.Typ, variables[0].Typ.Typ)
		}
		if variables[0].IsGlobal {
			storeGlobalVar(class, m.MakeClass.mainclass, code, variables[0])
			if variableType.JvmSlotSize() == 1 {
				state.popStack(1)
			} else {
				state.popStack(2)
			}
		} else {
			m.MakeClass.storeLocalVar(class, code, variables[0])
			if variables[0].BeenCaptured {
				currentStack -= 1
			}
			m.MakeClass.appendLocalVar(class, code, variables[0], state)
			if variableType.JvmSlotSize() == 1 {
				if variables[0].BeenCaptured {
					state.popStack(2)
				} else {
					state.popStack(1)
				}
			} else {
				if variables[0].BeenCaptured {
					state.popStack(3)
				} else {
					state.popStack(2)
				}
			}
		}

		variables = variables[1:]
	}
	return
}
