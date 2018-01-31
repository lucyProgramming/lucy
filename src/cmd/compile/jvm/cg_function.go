package jvm

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) mkFunc(f *ast.Function) {
	method := &cg.MethodHighLevel{}
	context := &Context{f, nil}
	if f.IsGlobal || f.IsClosureFunction() == false {
		m.buildFunction(m.mainclass, method, f, context)
		method.AccessFlags = 0
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		if f.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
		} else {
			method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		}
		method.Class = m.mainclass
		method.Name = f.Name
		method.Descriptor = f.MkDescriptor()
		m.mainclass.AppendMethod(method)
		return
	}

	class := m.mkClosureFunctionClass()
	m.buildFunction(class, method, f, context)
}
func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, method *cg.MethodHighLevel, f *ast.Function, context *Context) {
	f.ClassMethod = method
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
	}()
	method.Code.MaxLocals = f.Varoffset
	m.buildAtuoArrayListVar(class, &method.Code, context)
	m.buildBlock(class, &method.Code, f.Block, context)
	fmt.Println("$$$$$$$$$$$$$$$$", method.Code.MaxStack)
	panic(method.Code.MaxStack)
	method.Descriptor = f.Descriptor
	return
}
func (m *MakeClass) buildAtuoArrayListVar(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses("java/util/ArrayList", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/util/ArrayList",
		Name:       "<init>",
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	switch context.function.ArrayListVarForMultiReturn.Offset {
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
		if context.function.ArrayListVarForMultiReturn.Offset > 255 {
			panic("over 255")
		}
		code.Codes[code.CodeLength] = cg.OP_astore
		code.Codes[code.CodeLength+1] = byte(context.function.ArrayListVarForMultiReturn.Offset)
		code.CodeLength += 2
	}
	return
}
