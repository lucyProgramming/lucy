package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeClass *MakeClass) loadLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		return closure.loadLocalCloureVar(class, code, v)
	}
	maxstack = jvmSize(v.Type)
	copyOP(code, loadLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
	return
}

func (makeClass *MakeClass) storeLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		closure.storeLocalClosureVar(class, code, v)
		return
	}
	maxstack = jvmSize(v.Type)
	copyOP(code, storeLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
	return
}
