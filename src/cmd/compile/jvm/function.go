package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, method *cg.MethodHighLevel, f *ast.Function) {
	context := &Context{}
	context.function = f
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
