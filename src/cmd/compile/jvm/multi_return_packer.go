package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type MultiValuePacker struct {
}

/*
	stack is 1
*/
func (a *MultiValuePacker) storeMultiValueAutoVar(code *cg.AttributeCode, context *Context) {
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
		context.function.AutoVariableForMultiReturn.Offset)...)
}

func (a *MultiValuePacker) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, typ *ast.Type, context *Context) (maxStack uint16) {
	maxStack = a.unPackObject(class, code, k, context)
	if typ.IsPointer() == false {
		typeConverter.unPackPrimitives(class, code, typ)
		if t := jvmSlotSize(typ); t > maxStack {
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
	if k > 127 {
		panic("over 127")
	}
	maxStack = 2
	//a.buildLoadArrayListAutoVar(code, context) // local array list on stack
	copyOPs(code,
		loadLocalVariableOps(ast.VariableTypeObject, context.function.AutoVariableForMultiReturn.Offset)...)
	loadInt32(class, code, int32(k))
	code.Codes[code.CodeLength] = cg.OP_aaload
	code.CodeLength++
	return
}
