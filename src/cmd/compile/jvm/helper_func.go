package jvm

import (
	"fmt"
	"strings"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func writeExits(es []*cg.Exit, to int) {
	for _, e := range es {
		offset := int16(to - int(e.CurrentCodeLength))
		e.BranchBytes[0] = byte(offset >> 8)
		e.BranchBytes[1] = byte(offset)
	}
}

func gotoOffset(code *cg.AttributeCode, offset int) {
	if offset < 0 {
		panic("to is negative")
	}
	exit := (&cg.Exit{}).Init(cg.OP_goto, code)
	writeExits([]*cg.Exit{exit}, offset)
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
	var count uint16
	count = 1
	for _, v := range functionType.ParameterList {
		count += jvmSlotSize(v.Type)
	}
	if count > 255 {
		panic("over 255")
	}
	return byte(count)
}

func jvmSlotSize(typ *ast.Type) uint16 {
	if typ.RightValueValid() == false {
		panic("right value is not valid:" + typ.TypeString())
	}
	if typ.Type == ast.VariableTypeDouble ||
		typ.Type == ast.VariableTypeLong {
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

func newArrayBaseOnType(class *cg.ClassHighLevel, code *cg.AttributeCode, typ *ast.Type) {
	switch typ.Type {
	case ast.VariableTypeBool:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ArrayTypeBoolean
		code.CodeLength += 2
	case ast.VariableTypeByte:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ArrayTypeByte
		code.CodeLength += 2
	case ast.VariableTypeShort:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ArrayTypeShort
		code.CodeLength += 2
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ArrayTypeInt
		code.CodeLength += 2
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ArrayTypeLong
		code.CodeLength += 2
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ArrayTypeFloat
		code.CodeLength += 2
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_newarray
		code.Codes[code.CodeLength+1] = ArrayTypeDouble
		code.CodeLength += 2
	case ast.VariableTypeString:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(javaStringClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeMap:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(mapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeFunction:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(javaMethodHandleClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeObject:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(typ.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeArray:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		meta := ArrayMetas[typ.Array.Type]
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeJavaArray:
		code.Codes[code.CodeLength] = cg.OP_anewarray
		class.InsertClassConst(Descriptor.typeDescriptor(typ.Array), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
}
func storeArrayElementOp(typ ast.VariableTypeKind) (op byte) {
	switch typ {
	case ast.VariableTypeBool:
		op = cg.OP_bastore
	case ast.VariableTypeByte:
		op = cg.OP_bastore
	case ast.VariableTypeShort:
		op = cg.OP_sastore
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		op = cg.OP_iastore
	case ast.VariableTypeLong:
		op = cg.OP_lastore
	case ast.VariableTypeFloat:
		op = cg.OP_fastore
	case ast.VariableTypeDouble:
		op = cg.OP_dastore
	case ast.VariableTypeString:
		fallthrough
	case ast.VariableTypeMap:
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeArray:
		fallthrough
	case ast.VariableTypeJavaArray:
		op = cg.OP_aastore
	}
	return
}

/*
	return function call result stack size , this value should be proper handled
*/
func popCallResult(code *cg.AttributeCode, e *ast.Expression, ft *ast.FunctionType) (stackSize uint16) {
	stackSize = functionReturnJvmSize(ft)
	if ft.VoidReturn() == false && e.IsStatementExpression {
		if len(ft.ReturnList) == 1 {
			if jvmSlotSize(ft.ReturnList[0].Type) == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
				code.CodeLength++
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop2
				code.CodeLength++
			}
		} else { // > 1
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}

func functionReturnJvmSize(ft *ast.FunctionType) uint16 {
	if ft.VoidReturn() {
		return 0
	}
	if len(ft.ReturnList) > 1 {
		return 1
	} else {
		return jvmSlotSize(ft.ReturnList[0].Type)
	}
}

func setEnumArray(class *cg.ClassHighLevel, code *cg.AttributeCode, state *StackMapState, context *Context, e *ast.Enum) (maxStack uint16) {
	if e.DefaultValue == 0 {
		return
	}
	type autoVar = struct {
		k      uint16
		length uint16
	}
	var a autoVar
	a.k = code.MaxLocals
	a.length = code.MaxLocals + 1
	code.MaxLocals += 2
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	code.Codes[code.CodeLength] = cg.OP_arraylength
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, a.length)...)
	code.Codes[code.CodeLength] = cg.OP_iconst_0
	code.CodeLength++
	copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, a.k)...)
	state.appendLocals(class, &ast.Type{
		Type: ast.VariableTypeInt,
	})
	state.appendLocals(class, &ast.Type{
		Type: ast.VariableTypeInt,
	})
	context.MakeStackMap(code, state, code.CodeLength)
	offset := code.CodeLength
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, a.k)...)
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, a.length)...)
	exit := (&cg.Exit{}).Init(cg.OP_if_icmpge, code)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, a.k)...)
	loadInt32(class, code, e.DefaultValue)
	code.Codes[code.CodeLength] = cg.OP_iastore
	code.CodeLength++
	maxStack = 3
	code.Codes[code.CodeLength] = cg.OP_iinc
	code.Codes[code.CodeLength+1] = byte(a.k)
	code.Codes[code.CodeLength+2] = 1
	code.CodeLength += 3
	gotoOffset(code, offset)
	writeExits([]*cg.Exit{exit}, code.CodeLength)
	context.MakeStackMap(code, state, code.CodeLength)
	return

}
