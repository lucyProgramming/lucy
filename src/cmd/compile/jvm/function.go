package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) mkFunc(f *ast.Function) {
	method := &cg.MethodHighLevel{}
	context := &Context{}
	context.function = f
	if f.IsGlobal || f.IsClosureFunction() == false {
		method.Class = m.mainclass
		method.Name = f.Name
		method.Descriptor = f.MkDescriptor()
		method.AccessFlags = 0
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		if f.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
		} else {
			method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		}
		m.buildFunction(m.mainclass, method, f, context)
		m.mainclass.AppendMethod(method)
		return
	}
	class := m.mkClosureFunctionClass()
	m.buildFunction(class, method, f, context)
}
func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, method *cg.MethodHighLevel, f *ast.Function, context *Context) {
	context.method = method
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
	}()
	method.Code.MaxLocals = f.Varoffset
	m.buildBlock(class, &method.Code, f.Block, context)
	return
}
