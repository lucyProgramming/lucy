package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildVar(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	vs := e.Data.(*ast.ExpressionDeclareVariable)
	//first round
	index := len(vs.Vs) - 1
	currentStack := uint16(0)
	for index >= 0 {
		if vs.Vs[index].BeenCaptured {
			t := state.newObjectVariableType(closure.getMeta(vs.Vs[index].Typ.Typ).className)
			vs.Vs[index].LocalValOffset = state.appendLocals(class, code, t)
			stack := closure.createCloureVar(class, code, vs.Vs[index])
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			state.Stacks = append(state.Stacks,
				state.newStackMapVerificationTypeInfo(class, t))
			// load to stack
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, vs.Vs[index].LocalValOffset)...)
			currentStack += 1
		} else {
			vs.Vs[index].LocalValOffset = state.appendLocals(class, code, vs.Vs[index].Typ)
		}
		index--
	}
	//
	variables := vs.Vs
	for _, v := range vs.Values {
		if v.MayHaveMultiValue() && len(v.Values) > 1 {
			stack, _ := m.build(class, code, v, context, nil)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			m.buildStoreArrayListAutoVar(code, context)
			for kk, tt := range v.Values {
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
						state.popStack(1) // pop closure object
						currentStack -= 1
					}
				}
				variables = variables[1:]
			}
			continue
		}
		variableType := v.Value
		stack, es := m.build(class, code, v, context, state)
		if len(es) > 0 {
			state.Stacks = append(state.Stacks, state.newStackMapVerificationTypeInfo(class, v.Value))
			backPatchEs(es, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		if variableType.IsNumber() && variableType.Typ != variables[0].Typ.Typ {
			m.numberTypeConverter(code, variableType.Typ, variables[0].Typ.Typ)
		}
		if variables[0].IsGlobal {
			storeGlobalVar(class, m.MakeClass.mainclass, code, variables[0])
		} else {
			m.MakeClass.storeLocalVar(class, code, variables[0])
			if variables[0].BeenCaptured {
				state.popStack(1)
				currentStack -= 1
			}
		}
		variables = variables[1:]
	}

	return

}
