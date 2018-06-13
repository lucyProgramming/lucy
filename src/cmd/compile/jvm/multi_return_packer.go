package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type MultiValuePacker struct {
}

/*
	stack is 1
*/
func (a *MultiValuePacker) buildLoadArrayListAutoVar(code *cg.AttributeCode, context *Context) (maxStack uint16) {
	maxStack = 1
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForMultiReturn.Offset)...)
	return
}

/*
	stack is 1
*/
func (a *MultiValuePacker) storeArrayListAutoVar(code *cg.AttributeCode, context *Context) {
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT,
		context.function.AutoVarForMultiReturn.Offset)...)

}

func (a *MultiValuePacker) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, typ *ast.VariableType, context *Context) (maxStack uint16) {
	maxStack = a.unPackObject(class, code, k, context)
	if typ.IsPointer() == false {
		typeConverter.getFromObject(class, code, typ)
		if t := jvmSize(typ); t > maxStack {
			maxStack = t
		}
	} else {
		typeConverter.castPointerTypeToRealType(class, code, typ)
	}
	return
}

/*
	object is all i need
*/
func (a *MultiValuePacker) unPackObject(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, context *Context) (maxStack uint16) {
	maxStack = 2
	a.buildLoadArrayListAutoVar(code, context) // local array list on stack
	if k > 127 {
		panic("over 127")
	}
	loadInt32(class, code, int32(k))
	code.Codes[code.CodeLength] = cg.OP_aaload
	code.CodeLength++
	return
}
