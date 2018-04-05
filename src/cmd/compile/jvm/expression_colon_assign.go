package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildColonAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
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
			currentStack += 1
		}
		index--
	}
	variables := vs.Vs
	for _, v := range vs.Values {
		if v.MayHaveMultiValue() && len(v.VariableTypes) > 1 {
			stack, _ := m.build(class, code, v, context)
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
				}
				variables = variables[1:]
			}
			continue
		}
		variableType := v.VariableType
		if v.MayHaveMultiValue() {
			variableType = v.VariableTypes[0]
		}
		stack, es := m.build(class, code, v, context)
		backPatchEs(es, code.CodeLength)
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		if variables[0].Name == ast.NO_NAME_IDENTIFIER {
			if variableType.JvmSlotSize() == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
				code.CodeLength++
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop2
				code.CodeLength++
			}
		} else {
			if variableType.IsNumber() && variableType.Typ != variables[0].Typ.Typ {
				m.numberTypeConverter(code, variableType.Typ, variables[0].Typ.Typ)
			}
			if variables[0].IsGlobal {
				storeGlobalVar(class, m.MakeClass.mainclass, code, variables[0])
			} else {
				m.MakeClass.storeLocalVar(class, code, variables[0])
				if variables[0].BeenCaptured {
					currentStack -= 1
				}
			}
		}
		variables = variables[1:]
	}
	return
}
