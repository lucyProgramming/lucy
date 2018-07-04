package jvm

import (
	"fmt"
	"strings"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func fillOffsetForExits(es []*cg.Exit, to int) {
	for _, e := range es {
		offset := int16(to - int(e.CurrentCodeLength))
		e.BranchBytes[0] = byte(offset >> 8)
		e.BranchBytes[1] = byte(offset)
	}
}

func jumpTo(op byte, code *cg.AttributeCode, to int) {
	b := (&cg.Exit{}).FromCode(op, code)
	fillOffsetForExits([]*cg.Exit{b}, to)
}

func copyOPs(code *cg.AttributeCode, op ...byte) {
	for k, v := range op {
		code.Codes[code.CodeLength+k] = v
	}
	code.CodeLength += len(op)
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

func interfaceMethodArgsCount(functionType *ast.FunctionType) byte {
	var b uint16
	b = 1
	for _, v := range functionType.ParameterList {
		b += jvmSlotSize(v.Type)
	}
	if b > 255 {
		panic("over 255")
	}
	return byte(b)
}

func jvmSlotSize(typ *ast.Type) uint16 {
	if typ.RightValueValid() == false {
		panic("right value is not valid:" + typ.TypeString())
	}
	if typ.Type == ast.VariableTypeDouble || ast.VariableTypeLong == typ.Type {
		return 2
	} else {
		return 1
	}
}

func nameTemplateFunction(function *ast.Function) string {
	s := function.Name
	for _, v := range function.Type.ParameterList {
		if v.Type.IsPrimitive() {
			s += fmt.Sprintf("$%s", v.Type.TypeString())
			continue
		}
		switch v.Type.Type {
		case ast.VariableTypeObject:
			s += fmt.Sprintf("$%s", strings.Replace(v.Type.Class.Name, "/", "$", -1))
		case ast.VariableTypeMap:
			s += "_map"
		case ast.VariableTypeArray:
			s += "_array"
		case ast.VariableTypeJavaArray:
			s += "_java_array"
		case ast.VariableTypeEnum:
			s += fmt.Sprintf("$%s", strings.Replace(v.Type.Enum.Name, "/", "$", -1))
		}
	}
	return s
}

func insertTypeAssertClass(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.Type) {
	if t.Type == ast.VariableTypeObject {
		class.InsertClassConst(t.Class.Name, code.Codes[code.CodeLength:code.CodeLength+2])
	} else if t.Type == ast.VariableTypeArray { // arrays
		meta := ArrayMetas[t.Array.Type]
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength:code.CodeLength+2])
	} else {
		class.InsertClassConst(Descriptor.typeDescriptor(t), code.Codes[code.CodeLength:code.CodeLength+2])
	}
	code.CodeLength += 2
}
