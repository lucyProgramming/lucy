package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type MultiValueAutoVar struct {
	localVarOffset uint16
}

/*
	stack is 1 ,expect value on stack
*/
func newMultiValueAutoVar(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	state *StackMapState) *MultiValueAutoVar {

	ret := &MultiValueAutoVar{}
	ret.localVarOffset = code.MaxLocals
	code.MaxLocals++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
		ret.localVarOffset)...)
	state.appendLocals(class, state.newObjectVariableType(javaRootObjectArray))
	return ret
}

func (this *MultiValueAutoVar) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode,
	valueIndex int, typ *ast.Type) (maxStack uint16) {
	maxStack = this.unPack2Object(class, code, valueIndex)
	if typ.IsPointer() == false {
		typeConverter.unPackPrimitives(class, code, typ)
		if t := jvmSlotSize(typ); t > maxStack {
			maxStack = t
		}
	} else {
		typeConverter.castPointer(class, code, typ)
	}
	return
}

/*
	object is all i need
*/
func (this *MultiValueAutoVar) unPack2Object(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	valueIndex int) (maxStack uint16) {
	if valueIndex > 127 {
		panic("over 127")
	}
	maxStack = 2
	//a.buildLoadArrayListAutoVar(code, context) // local array list on stack
	copyOPs(code,
		loadLocalVariableOps(ast.VariableTypeObject, this.localVarOffset)...)
	loadInt32(class, code, int32(valueIndex))
	code.Codes[code.CodeLength] = cg.OP_aaload
	code.CodeLength++
	return
}
