package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

type ArrayListPacker struct {
}

/*
	stack is 1
*/
func (m *ArrayListPacker) buildLoadArrayListAutoVar(code *cg.AttributeCode, context *Context) {
	switch context.function.AutoVarForMultiReturn.Offset {
	case 0:
		code.Codes[code.CodeLength] = cg.OP_aload_0
		code.CodeLength++
	case 1:
		code.Codes[code.CodeLength] = cg.OP_aload_1
		code.CodeLength++
	case 2:
		code.Codes[code.CodeLength] = cg.OP_aload_2
		code.CodeLength++
	case 3:
		code.Codes[code.CodeLength] = cg.OP_aload_3
		code.CodeLength++
	default:
		if context.function.AutoVarForMultiReturn.Offset > 255 {
			panic("local var offset over 255")
		}
		code.Codes[code.CodeLength] = cg.OP_aload
		code.Codes[code.CodeLength+1] = byte(context.function.AutoVarForMultiReturn.Offset)
		code.CodeLength += 2
	}
}

/*
	stack is 1
*/
func (a *ArrayListPacker) buildStoreArrayListAutoVar(code *cg.AttributeCode, context *Context) {
	switch context.function.AutoVarForMultiReturn.Offset {
	case 0:
		code.Codes[code.CodeLength] = cg.OP_astore_0
		code.CodeLength++
	case 1:
		code.Codes[code.CodeLength] = cg.OP_astore_1
		code.CodeLength++
	case 2:
		code.Codes[code.CodeLength] = cg.OP_astore_2
		code.CodeLength++
	case 3:
		code.Codes[code.CodeLength] = cg.OP_astore_3
		code.CodeLength++
	default:
		if context.function.AutoVarForMultiReturn.Offset > 255 {
			panic("local var offset over 255")
		}
		code.Codes[code.CodeLength] = cg.OP_astore
		code.Codes[code.CodeLength+1] = byte(context.function.AutoVarForMultiReturn.Offset)
		code.CodeLength += 2
	}
}

func (a *ArrayListPacker) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode, k int, typ *ast.VariableType, context *Context) (maxstack uint16) {
	maxstack = 2
	a.buildLoadArrayListAutoVar(code, context) // local array list on stack
	switch k {
	case 0:
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
	case 1:
		code.Codes[code.CodeLength] = cg.OP_iconst_1
		code.CodeLength++
	case 2:
		code.Codes[code.CodeLength] = cg.OP_iconst_2
		code.CodeLength++
	case 3:
		code.Codes[code.CodeLength] = cg.OP_iconst_3
		code.CodeLength++
	case 4:
		code.Codes[code.CodeLength] = cg.OP_iconst_4
		code.CodeLength++
	case 5:
		code.Codes[code.CodeLength] = cg.OP_iconst_5
		code.CodeLength++
	default:
		if k > 127 {
			panic("over 127")
		}
		code.Codes[code.CodeLength] = cg.OP_bipush
		code.Codes[code.CodeLength+1] = byte(k)
		code.CodeLength += 2
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     "get",
		Descriptor: "(I)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if typ.IsPointer() == false {
		primitiveObjectConverter.getFromObject(class, code, typ)
	} else {
		primitiveObjectConverter.castPointerTypeToRealType(class, code, typ)
	}
	return
}

func (a *ArrayListPacker) unPackObject(class *cg.ClassHighLevel, code *cg.AttributeCode,
	k int, context *Context) (maxstack uint16) {
	maxstack = 2
	a.buildLoadArrayListAutoVar(code, context) // local array list on stack
	switch k {
	case 0:
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
	case 1:
		code.Codes[code.CodeLength] = cg.OP_iconst_1
		code.CodeLength++
	case 2:
		code.Codes[code.CodeLength] = cg.OP_iconst_2
		code.CodeLength++
	case 3:
		code.Codes[code.CodeLength] = cg.OP_iconst_3
		code.CodeLength++
	case 4:
		code.Codes[code.CodeLength] = cg.OP_iconst_4
		code.CodeLength++
	case 5:
		code.Codes[code.CodeLength] = cg.OP_iconst_5
		code.CodeLength++
	default:
		if k > 127 {
			panic("over 127")
		}
		code.Codes[code.CodeLength] = cg.OP_bipush
		code.Codes[code.CodeLength+1] = byte(k)
		code.CodeLength += 2
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_arrylist_class,
		Method:     "get",
		Descriptor: "(I)Ljava/lang/Object;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	return
}
