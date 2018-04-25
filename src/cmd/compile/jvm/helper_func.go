package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	make a default construction
*/
func mkClassDefaultContruction(class *cg.ClassHighLevel) {
	method := &cg.MethodHighLevel{}
	method.Name = special_method_init
	method.Descriptor = "()V"
	method.AccessFlags |= cg.ACC_METHOD_PUBLIC
	method.Code = &cg.AttributeCode{}
	length := 5
	method.Code.Codes = make([]byte, length)
	method.Code.Codes[0] = cg.OP_aload_0
	method.Code.Codes[1] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      class.SuperClass,
		Method:     special_method_init,
		Descriptor: "()V",
	}, method.Code.Codes[2:4])
	method.Code.Codes[4] = cg.OP_return
	method.Code.MaxStack = 1
	method.Code.MaxLocals = 1
	method.Code.CodeLength = length
	class.AppendMethod(method)
}

func backPatchEs(es []*cg.JumpBackPatch, t int) {
	for _, e := range es {
		offset := int16(t - int(e.CurrentCodeLength))
		e.Bs[0] = byte(offset >> 8)
		e.Bs[1] = byte(offset)
	}
}

func jumpto(op byte, code *cg.AttributeCode, to int) {
	b := (&cg.JumpBackPatch{}).FromCode(op, code)
	backPatchEs([]*cg.JumpBackPatch{b}, to)
}

func copyOP(code *cg.AttributeCode, op ...byte) {
	for k, v := range op {
		code.Codes[code.CodeLength+k] = v
	}
	code.CodeLength += len(op)
}

func copyOPLeftValue(class *cg.ClassHighLevel, code *cg.AttributeCode, ops []byte, classname, name, descriptor string) {
	if len(ops) == 0 {
		return
	}
	code.Codes[code.CodeLength] = ops[0]
	code.CodeLength++
	if classname != "" || name != "" || descriptor != "" {
		if classname == "" || name == "" || descriptor == "" {
			panic("....")
		}
		if ops[0] == cg.OP_putstatic || ops[0] == cg.OP_putfield {
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      classname,
				Field:      name,
				Descriptor: descriptor,
			}, code.Codes[code.CodeLength:code.CodeLength+2])
		} else if ops[0] == cg.OP_invokevirtual {
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      classname,
				Method:     name,
				Descriptor: descriptor,
			}, code.Codes[code.CodeLength:code.CodeLength+2])
		} else {
			panic("...")
		}
		code.CodeLength += 2
	}
	copyOP(code, ops[1:]...)
}
func postion(class *cg.ClassHighLevel, code *cg.AttributeCode, s string) {
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst(s, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_pop
	code.CodeLength += 4

}
func loadInt32(class *cg.ClassHighLevel, code *cg.AttributeCode, value int32) {
	switch value {
	case -1:
		code.Codes[code.CodeLength] = cg.OP_iconst_m1
		code.CodeLength++
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
		if -127 >= value && value <= 128 {
			code.Codes[code.CodeLength] = cg.OP_bipush
			code.Codes[code.CodeLength+1] = byte(value)
			code.CodeLength += 2
		} else if -32768 <= value && 32767 >= value {
			code.Codes[code.CodeLength] = cg.OP_sipush
			code.Codes[code.CodeLength+1] = byte(int16(value) >> 8)
			code.Codes[code.CodeLength+2] = byte(value)
			code.CodeLength += 3
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertIntConst(int32(value), code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
}

func checkStackTopIfNagetiveThrowIndexOutOfRangeException(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context, state *StackMapState) (increment uint16) {

	increment = 1
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	context.MakeStackMap(code, state, code.CodeLength+6)
	context.MakeStackMap(code, state, code.CodeLength+15)
	code.Codes[code.CodeLength] = cg.OP_iflt
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 6)
	code.Codes[code.CodeLength+3] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 12)
	code.Codes[code.CodeLength+6] = cg.OP_pop
	code.Codes[code.CodeLength+7] = cg.OP_new
	class.InsertClassConst(java_index_out_of_range_exception_class, code.Codes[code.CodeLength+8:code.CodeLength+10])
	code.Codes[code.CodeLength+10] = cg.OP_dup
	code.Codes[code.CodeLength+11] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_index_out_of_range_exception_class,
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+12:code.CodeLength+14])
	code.Codes[code.CodeLength+14] = cg.OP_athrow
	code.CodeLength += 15
	return
}

func storeGlobalVar(class *cg.ClassHighLevel, mainClass *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) {
	code.Codes[code.CodeLength] = cg.OP_putstatic
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      mainClass.Name,
		Field:      v.Name,
		Descriptor: Descriptor.typeDescriptor(v.Typ),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
}

func interfaceMethodArgsCount(ft *ast.FunctionType) byte {
	var b byte
	b = 1
	for _, v := range ft.ParameterList {
		b += byte(jvmSize(v.Typ))
	}
	return b
}

func jvmSize(v *ast.VariableType) uint16 {
	if v.RightValueValid() == false {
		panic("right value is not valid," + v.TypeString())
	}
	if v.Typ == ast.VARIABLE_TYPE_DOUBLE || ast.VARIABLE_TYPE_LONG == v.Typ {
		return 2
	} else {
		return 1
	}
}
