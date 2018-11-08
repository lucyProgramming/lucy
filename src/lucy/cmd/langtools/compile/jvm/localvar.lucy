package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildPackage) loadLocalVar(class *cg.ClassHighLevel,
	code *cg.AttributeCode, v *ast.Variable) (maxStack uint16) {
	if v.BeenCapturedAsLeftValue > 0 {
		return closure.loadLocalClosureVar(class, code, v)
	}
	maxStack = jvmSlotSize(v.Type)
	copyOPs(code, loadLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
	return
}

func (this *BuildPackage) storeLocalVar(class *cg.ClassHighLevel,
	code *cg.AttributeCode, v *ast.Variable) (maxStack uint16) {
	if v.BeenCapturedAsLeftValue > 0 {
		closure.storeLocalClosureVar(class, code, v)
		return
	}
	maxStack = jvmSlotSize(v.Type)
	copyOPs(code, storeLocalVariableOps(v.Type.Type, v.LocalValOffset)...)
	return
}
