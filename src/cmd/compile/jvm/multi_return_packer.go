package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type MultiValueAutoVar struct {
	offset uint16
}

/*
	stack is 1 ,expect value on stack
*/
func storeMultiValueAutoVar(class *cg.ClassHighLevel, code *cg.AttributeCode, state *StackMapState) *MultiValueAutoVar {
	ret := &MultiValueAutoVar{}
	ret.offset = code.MaxLocals
	code.MaxLocals++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeObject,
		ret.offset)...)
	state.appendLocals(class, state.newObjectVariableType(javaRootObjectArray))
	return ret
}

func (packer *MultiValueAutoVar) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, typ *ast.Type, context *Context) (maxStack uint16) {
	maxStack = packer.unPackObject(class, code, k, context)
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
func (packer *MultiValueAutoVar) unPackObject(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, context *Context) (maxStack uint16) {
	if k > 127 {
		panic("over 127")
	}
	maxStack = 2
	//a.buildLoadArrayListAutoVar(code, context) // local array list on stack
	copyOPs(code,
		loadLocalVariableOps(ast.VariableTypeObject, packer.offset)...)
	loadInt32(class, code, int32(k))
	code.Codes[code.CodeLength] = cg.OP_aaload
	code.CodeLength++
	return
}
