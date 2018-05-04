package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type ArrayListPacker struct {
}

/*
	stack is 1
*/
func (a *ArrayListPacker) buildLoadArrayListAutoVar(code *cg.AttributeCode, context *Context) (maxstack uint16) {
	maxstack = 1
	copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, context.function.AutoVarForMultiReturn.Offset)...)
	return
}

/*
	stack is 1
*/
func (a *ArrayListPacker) storeArrayListAutoVar(code *cg.AttributeCode, context *Context) {
	copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT,
		context.function.AutoVarForMultiReturn.Offset)...)

}

func (a *ArrayListPacker) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, typ *ast.VariableType, context *Context) (maxstack uint16) {
	maxstack = a.unPackObject(class, code, k, context)
	if typ.IsPointer() == false {
		typeConverter.getFromObject(class, code, typ)
		if t := jvmSize(typ); t > maxstack {
			maxstack = t
		}
	} else {
		typeConverter.castPointerTypeToRealType(class, code, typ)
	}
	return
}

/*
	object is all i need
*/
func (a *ArrayListPacker) unPackObject(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, context *Context) (maxstack uint16) {
	maxstack = 2
	a.buildLoadArrayListAutoVar(code, context) // local array list on stack
	if k > 127 {
		panic("over 127")
	}
	loadInt32(class, code, int32(k))
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     "get",
		Descriptor: "(I)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
