package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildVar(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	vs := e.Data.(*ast.ExpressionDeclareVariable)
	for _, v := range vs.Variables {
		if v.BeenCaptured {
			v.LocalValOffset = code.MaxLocals
			code.MaxLocals++
		} else {
			v.LocalValOffset = code.MaxLocals
			code.MaxLocals += jvmSlotSize(v.Type)
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
		closure.createClosureVar(class, code, v.Type)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		{
			t := state.newObjectVariableType(closure.getMeta(v.Type.Type).className)
			state.pushStack(class, t)
			state.pushStack(class, t)
		}
		currentStack += 2
		index--
	}
	index = 0
	for _, v := range vs.InitValues {
		if v.HaveMultiValue() {
			stack := buildExpression.build(class, code, vs.InitValues[0], context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			autoVar := newMultiValueAutoVar(class, code, state)
			for kk, tt := range v.MultiValues {
				stack = autoVar.unPack(class, code, kk, tt)
				if t := stack + currentStack; t > maxStack {
					maxStack = t
				}
				if vs.Variables[index].IsGlobal {
					buildExpression.BuildPackage.storeGlobalVariable(class, code, vs.Variables[index])
					index++
					continue
				}
				buildExpression.BuildPackage.storeLocalVar(class, code, vs.Variables[index])
				if vs.Variables[index].BeenCaptured {
					copyOPs(code, storeLocalVariableOps(vs.Variables[index].Type.Type, vs.Variables[index].LocalValOffset)...)
					state.popStack(2)
					state.appendLocals(class, state.newObjectVariableType(closure.getMeta(vs.Variables[index].Type.Type).className))
				} else {
					state.appendLocals(class, vs.Variables[index].Type)
				}
				index++
			}
			continue
		}
		//
		stack := buildExpression.build(class, code, vs.InitValues[0], context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if vs.Variables[index].IsGlobal {
			buildExpression.BuildPackage.storeGlobalVariable(class, code, vs.Variables[index])
			index++
			continue
		}
		buildExpression.BuildPackage.storeLocalVar(class, code, vs.Variables[index])
		if vs.Variables[index].BeenCaptured {
			copyOPs(code, storeLocalVariableOps(vs.Variables[index].Type.Type, vs.Variables[index].LocalValOffset)...)
			state.popStack(2)
			state.appendLocals(class, state.newObjectVariableType(closure.getMeta(vs.Variables[index].Type.Type).className))
		} else {
			state.appendLocals(class, vs.Variables[index].Type)
		}
		index++
	}
	return
}
