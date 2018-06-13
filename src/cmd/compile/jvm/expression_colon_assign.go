package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildColonAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	vs := e.Data.(*ast.ExpressionDeclareVariable)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if len(vs.Variables) == 1 {
		v := vs.Variables[0]
		currentStack := uint16(0)
		if v.Name != ast.NO_NAME_IDENTIFIER && v.BeenCaptured {
			obj := state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className)
			if vs.IfDeclareBefor[0] {
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
				currentStack = 1
				state.pushStack(class, obj)
			} else {
				closure.createCloureVar(class, code, v.Typ)
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				currentStack = 2
				state.pushStack(class, obj)
				state.pushStack(class, obj)
			}
		}
		stack, es := m.build(class, code, vs.Values[0], context, state)
		if len(es) > 0 {
			backfillExit(es, code.CodeLength)
			state.pushStack(class, vs.Values[0].Value)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Name == ast.NO_NAME_IDENTIFIER {
			if jvmSize(vs.Values[0].Value) == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop2
			}
			code.CodeLength++
			return
		}
		maxStack += currentStack
		if v.IsGlobal {
			storeGlobalVar(class, m.MakeClass.mainclass, code, vs.Variables[0])
		} else {
			if vs.IfDeclareBefor[0] {
				m.MakeClass.storeLocalVar(class, code, vs.Variables[0])
			} else {
				v.LocalValOffset = code.MaxLocals
				if v.BeenCaptured {
					code.MaxLocals++
				} else {
					code.MaxLocals += jvmSize(v.Typ)
				}
				m.MakeClass.storeLocalVar(class, code, v)
				if vs.Variables[0].BeenCaptured {
					copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
					state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className))
				} else {
					state.appendLocals(class, v.Typ)
				}
			}
		}
		return
	}
	if len(vs.Values) == 1 {
		maxStack, _ = m.build(class, code, vs.Values[0], context, state)
	} else {
		maxStack = m.buildExpressions(class, code, vs.Values, context, state)
	}
	multiValuePacker.storeArrayListAutoVar(code, context)
	//first round
	for k, v := range vs.Variables {
		if v.Name == ast.NO_NAME_IDENTIFIER {
			continue
		}
		if v.IsGlobal {
			stack := multiValuePacker.unPack(class, code, k, v.Typ, context)
			if stack > maxStack {
				maxStack = stack
			}
			storeGlobalVar(class, m.MakeClass.mainclass, code, v)
			continue
		}
		//this variable not been captured,also not declared here
		if vs.IfDeclareBefor[k] {
			if v.BeenCaptured {
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
				stack := multiValuePacker.unPack(class, code, k, v.Typ, context)
				if t := 1 + stack; t > maxStack {
					maxStack = t
				}
			} else {
				stack := multiValuePacker.unPack(class, code, k, v.Typ, context)
				if stack > maxStack {
					maxStack = stack
				}
			}
			m.MakeClass.storeLocalVar(class, code, v)
			continue
		}
		v.LocalValOffset = code.MaxLocals
		currentStack := uint16(0)
		if v.BeenCaptured {
			code.MaxLocals++
			stack := closure.createCloureVar(class, code, v.Typ)
			if stack > maxStack {
				maxStack = stack
			}
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
			if 2 > maxStack {
				maxStack = 2
			}
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
			currentStack = 1
			state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className))
		} else {
			code.MaxLocals += jvmSize(v.Typ)
			state.appendLocals(class, v.Typ)
		}
		if t := currentStack + multiValuePacker.unPack(class, code, k, v.Typ, context); t > maxStack {
			maxStack = t
		}
		m.MakeClass.storeLocalVar(class, code, v)
	}
	return
}

func (m *MakeExpression) buildVar(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	vs := e.Data.(*ast.ExpressionDeclareVariable)
	for _, v := range vs.Variables {
		if v.BeenCaptured {
			v.LocalValOffset = code.MaxLocals
			code.MaxLocals++
		} else {
			v.LocalValOffset = code.MaxLocals
			code.MaxLocals += jvmSize(v.Typ)
		}
	}
	index := len(vs.Variables) - 1
	currentStack := uint16(0)
	for index >= 0 {
		if vs.Variables[index].BeenCaptured == false {
			index--
			continue
		}
		v := vs.Variables[index]
		closure.createCloureVar(class, code, v.Typ)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		{
			t := state.newObjectVariableType(closure.getMeta(v.Typ.Typ).className)
			state.pushStack(class, t)
			state.pushStack(class, t)
		}
		currentStack += 2
		index--
	}
	index = 0
	for _, v := range vs.Values {
		if v.MayHaveMultiValue() && len(v.Values) > 1 {
			stack, _ := m.build(class, code, vs.Values[0], context, state)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			for kk, tt := range v.Values {
				stack = multiValuePacker.unPack(class, code, kk, tt, context)
				if t := stack + currentStack; t > maxstack {
					maxstack = t
				}
				if vs.Variables[index].IsGlobal {
					storeGlobalVar(class, m.MakeClass.mainclass, code, vs.Variables[index])
					index++
					continue
				}
				m.MakeClass.storeLocalVar(class, code, vs.Variables[index])
				if vs.Variables[index].BeenCaptured {
					copyOP(code, storeSimpleVarOp(vs.Variables[index].Typ.Typ, vs.Variables[index].LocalValOffset)...)
					state.popStack(2)
					state.appendLocals(class, state.newObjectVariableType(closure.getMeta(vs.Variables[index].Typ.Typ).className))
				} else {
					state.appendLocals(class, vs.Variables[index].Typ)
				}
				index++
			}
			continue
		}
		//
		stack, es := m.build(class, code, vs.Values[0], context, state)
		if len(es) > 0 {
			backfillExit(es, code.CodeLength)
			state.pushStack(class, v.Value)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if vs.Variables[index].IsGlobal {
			storeGlobalVar(class, m.MakeClass.mainclass, code, vs.Variables[index])
			index++
			continue
		}
		m.MakeClass.storeLocalVar(class, code, vs.Variables[index])
		if vs.Variables[index].BeenCaptured {
			copyOP(code, storeSimpleVarOp(vs.Variables[index].Typ.Typ, vs.Variables[index].LocalValOffset)...)
			state.popStack(2)
			state.appendLocals(class, state.newObjectVariableType(closure.getMeta(vs.Variables[index].Typ.Typ).className))
		} else {
			state.appendLocals(class, vs.Variables[index].Typ)
		}
		index++
	}

	return
}
