package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildExpression) buildVarAssign(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	vs := e.Data.(*ast.ExpressionVarAssign)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if len(vs.Lefts) == 1 {
		v := vs.Lefts[0].Data.(*ast.ExpressionIdentifier).Variable
		//fmt.Println(v.Name, v.Pos.ErrMsgPrefix())
		currentStack := uint16(0)
		if v.BeenCapturedAsLeftValue > 0 {
			closure.createClosureVar(class, code, v.Type)
			code.Codes[code.CodeLength] = cg.OP_dup
			code.CodeLength++
			currentStack = 2
			obj := state.newObjectVariableType(closure.getMeta(v.Type.Type).className)
			state.pushStack(class, obj)
			state.pushStack(class, obj)
		}
		stack := this.build(class, code, vs.InitValues[0], context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.IsGlobal {
			this.BuildPackage.storeGlobalVariable(class, code, v)
		} else {
			v.LocalValOffset = code.MaxLocals

			this.BuildPackage.storeLocalVar(class, code, v)
			if v.BeenCapturedAsLeftValue > 0 {
				code.MaxLocals++
				copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, v.LocalValOffset)...)
				state.appendLocals(class, state.newObjectVariableType(closure.getMeta(v.Type.Type).className))
			} else {
				code.MaxLocals += jvmSlotSize(v.Type)
				state.appendLocals(class, v.Type)
			}
		}
		return
	}
	if len(vs.InitValues) == 1 {
		maxStack = this.build(class, code, vs.InitValues[0], context, state)
	} else {
		maxStack = this.buildExpressions(class, code, vs.InitValues, context, state)
	}
	autoVar := newMultiValueAutoVar(class, code, state)
	for k, v := range vs.Lefts {
		if v.Type != ast.ExpressionTypeIdentifier {
			stack, remainStack, ops, _ := this.getLeftValue(class, code, v, context, state)
			if stack > maxStack {
				maxStack = stack
			}
			if t := remainStack + autoVar.unPack(class, code, k, v.Value); t > maxStack {
				maxStack = t
			}
			copyOPs(code, ops...)
			continue
		}
		//identifier
		identifier := v.Data.(*ast.ExpressionIdentifier)
		if identifier.Name == ast.UnderScore {
			continue
		}
		variable := identifier.Variable
		if variable.IsGlobal {
			stack := autoVar.unPack(class, code, k, variable.Type)
			if stack > maxStack {
				maxStack = stack
			}
			this.BuildPackage.storeGlobalVariable(class, code, variable)
			continue
		}
		//this variable not been captured,also not declared here
		if vs.IfDeclaredBefore[k] {
			if variable.BeenCapturedAsLeftValue > 0 {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, variable.LocalValOffset)...)
				stack := autoVar.unPack(class, code, k, variable.Type)
				if t := 1 + stack; t > maxStack {
					maxStack = t
				}
			} else {
				stack := autoVar.unPack(class, code, k, variable.Type)
				if stack > maxStack {
					maxStack = stack
				}
			}
			this.BuildPackage.storeLocalVar(class, code, variable)
		} else {
			variable.LocalValOffset = code.MaxLocals
			currentStack := uint16(0)
			if variable.BeenCapturedAsLeftValue > 0 {
				code.MaxLocals++
				stack := closure.createClosureVar(class, code, variable.Type)
				if stack > maxStack {
					maxStack = stack
				}
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				if 2 > maxStack {
					maxStack = 2
				}
				copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject, variable.LocalValOffset)...)
				currentStack = 1
				state.appendLocals(class, state.newObjectVariableType(closure.getMeta(variable.Type.Type).className))
			} else {
				code.MaxLocals += jvmSlotSize(variable.Type)
				state.appendLocals(class, variable.Type)
			}
			if t := currentStack + autoVar.unPack(class, code, k, variable.Type); t > maxStack {
				maxStack = t
			}
			this.BuildPackage.storeLocalVar(class, code, variable)
		}
	}
	return
}
