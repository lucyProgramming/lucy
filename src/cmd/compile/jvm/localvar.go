package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) loadLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		return closure.loadLocalCloureVar(class, code, v)
	}
	maxstack = jvmSize(v.Typ)
	copyOP(code, loadSimpleVarOps(v.Typ.Typ, v.LocalValOffset)...)
	return
}

func (m *MakeClass) storeLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		closure.storeLocalClosureVar(class, code, v)
		return
	}
	maxstack = jvmSize(v.Typ)
	copyOP(code, storeSimpleVarOps(v.Typ.Typ, v.LocalValOffset)...)
	return
}
