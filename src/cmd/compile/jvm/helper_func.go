package jvm

import (
	"fmt"
	"strings"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func backPatchEs(es []*cg.JumpBackPatch, to int) {
	for _, e := range es {
		offset := int16(to - int(e.CurrentCodeLength))
		e.Bs[0] = byte(offset >> 8)
		e.Bs[1] = byte(offset)
	}
}

func jumpTo(op byte, code *cg.AttributeCode, to int) {
	b := (&cg.JumpBackPatch{}).FromCode(op, code)
	backPatchEs([]*cg.JumpBackPatch{b}, to)
}

func copyOP(code *cg.AttributeCode, op ...byte) {
	for k, v := range op {
		code.Codes[code.CodeLength+k] = v
	}
	code.CodeLength += len(op)
}

func copyOPLeftValue(class *cg.ClassHighLevel, code *cg.AttributeCode, ops []byte, classname,
	name, descriptor string) {
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
			class.InsertIntConst(value, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
}

func storeGlobalVar(class *cg.ClassHighLevel, mainClass *cg.ClassHighLevel, code *cg.AttributeCode,
	v *ast.VariableDefinition) {
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
		panic("right value is not valid:" + v.TypeString())
	}
	if v.Typ == ast.VARIABLE_TYPE_DOUBLE || ast.VARIABLE_TYPE_LONG == v.Typ {
		return 2
	} else {
		return 1
	}
}

func nameTemplateFunction(f *ast.Function) string {
	s := f.Name
	for _, v := range f.Typ.ParameterList {
		if v.Typ.IsPrimitive() {
			s += fmt.Sprintf("_%s", v.Typ.TypeString())
			continue
		}
		switch v.Typ.Typ {
		case ast.VARIABLE_TYPE_OBJECT:
			s += fmt.Sprintf("_%s", strings.Replace(v.Typ.Class.Name, "/", "_", -1))
		case ast.VARIABLE_TYPE_MAP:
			s += "_map"
		case ast.VARIABLE_TYPE_ARRAY:
			s += "_array"
		case ast.VARIABLE_TYPE_JAVA_ARRAY:
			s += "_java_array"
		case ast.VARIABLE_TYPE_ENUM:
			s += fmt.Sprintf("_%s", strings.Replace(v.Typ.Enum.Name, "/", "_", -1))
		}
	}
	return s
}
