package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) mkFunc(f *ast.Function) {
	if f.IsGlobal || f.IsClosureFunction() {
		context := &Context{f}
		method := m.buildFunction(m.mainclass, f, context, true)
		m.mainclass.Methods[method.Name] = []*cg.MethodHighLevel{method}
		method.AccessFlags = 0
		if method.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
		}
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.ClassHighLevel = m.mainclass
		return
	}
	context := &Context{f}
	class := m.mkClosureFunctionClass()
	m.buildFunction(class, f, context, false)
}
func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, f *ast.Function, context *Context, isstatic bool) *cg.MethodHighLevel {
	ret := &cg.MethodHighLevel{}
	ret.Name = f.Name
	f.Method = ret
	ret.Code.Codes = make([]byte, 65536)
	ret.Code.CodeLength = 0
	m.buildAtuoArrayListVar(class, &ret.Code, context)
	m.buildBlock(class, &ret.Code, f.Block, context)
	ret.Descriptor = f.Descriptor
	return ret
}
func (m *MakeClass) buildAtuoArrayListVar(class *cg.ClassHighLevel, code *cg.AttributeCode, context *Context) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClasses("java/util/ArrayList", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
		Class: "java/util/ArrayList",
		Name:  "<init>",
		Type:  "()V",
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
