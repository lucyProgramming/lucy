package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildColonAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	vs := e.Data.(*ast.ExpressionDeclareVariable)
	//first round
	index := len(vs.Vs) - 1
	currentStack := uint16(0)
	for index >= 0 {
		if vs.Vs[index].Name == ast.NO_NAME_IDENTIFIER {
			index--
			continue
		}
		//this variable not been captured,also not declared here
		if vs.Vs[index].BeenCaptured {
			t := state.newObjectVariableType(closure.getMeta(vs.Vs[index].Typ.Typ).className)
			if vs.IfDeclareBefor[index] == false {
				vs.Vs[index].LocalValOffset = code.MaxLocals
				code.MaxLocals += 1
				stack := closure.createCloureVar(class, code, vs.Vs[index].Typ)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				state.Stacks = append(state.Stacks,
					state.newStackMapVerificationTypeInfo(class, t))
				state.Stacks = append(state.Stacks,
					state.newStackMapVerificationTypeInfo(class, t))
				currentStack += 2
			} else {
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, vs.Vs[index].LocalValOffset)...)
				state.Stacks = append(state.Stacks,
					state.newStackMapVerificationTypeInfo(class, t))
				currentStack += 1
			}
		} else {
			vs.Vs[index].LocalValOffset = code.MaxLocals
			code.MaxLocals += jvmSize(vs.Vs[index].Typ)
		}
		index--
	}
	variables := vs.Vs
	ifcreateBefore := vs.IfDeclareBefor
	slice := func() {
		variables = variables[1:]
		ifcreateBefore = ifcreateBefore[1:]
	}
	for _, v := range vs.Values {
		if v.MayHaveMultiValue() && len(v.Values) > 1 {
			stack, _ := m.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			m.buildStoreArrayListAutoVar(code, context)
			for kk, tt := range v.Values {
				if variables[0].Name == ast.NO_NAME_IDENTIFIER {
					slice()
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
						if ifcreateBefore[0] {
							state.popStack(1)
							currentStack -= 1
						} else {
							copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, variables[0].LocalValOffset)...)
							state.popStack(2)
							currentStack -= 2
						}
						state.appendLocals(class,
							state.newObjectVariableType(closure.getMeta(variables[0].Typ.Typ).className))
					} else {
						state.appendLocals(class, variables[0].Typ)
					}
				}
				slice()
			}
			continue
		}
		variableType := v.Value
		if v.MayHaveMultiValue() {
			variableType = v.Values[0]
		}
		stack, es := m.build(class, code, v, context, state)
		if len(es) > 0 {
			state.Stacks = append(state.Stacks,
				state.newStackMapVerificationTypeInfo(class, v.Value))
			backPatchEs(es, code.CodeLength) //
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1) // must be bool expression
		}
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		if variables[0].Name == ast.NO_NAME_IDENTIFIER {
			slice()
			if jvmSize(variableType) == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
			} else { // 2
				code.Codes[code.CodeLength] = cg.OP_pop2
			}
			continue
		}
		if variableType.IsNumber() && variableType.Typ != variables[0].Typ.Typ {
			m.numberTypeConverter(code, variableType.Typ, variables[0].Typ.Typ)
		}
		if variables[0].IsGlobal {
			storeGlobalVar(class, m.MakeClass.mainclass, code, variables[0])
		} else {
			m.MakeClass.storeLocalVar(class, code, variables[0])
			if variables[0].BeenCaptured {
				if ifcreateBefore[0] {
					state.popStack(1)
					currentStack -= 1
				} else {
					copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, variables[0].LocalValOffset)...)
					state.popStack(2)
					currentStack -= 2
				}
				state.appendLocals(class,
					state.newObjectVariableType(closure.getMeta(variables[0].Typ.Typ).className))
			} else {
				state.appendLocals(class, variables[0].Typ)
			}
		}
		slice()
	}
	return
}
